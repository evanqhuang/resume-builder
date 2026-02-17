package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/evanqhuang/resume-cli/config"
	"github.com/evanqhuang/resume-cli/generator"
	"github.com/evanqhuang/resume-cli/matching"
	"github.com/evanqhuang/resume-cli/resume"
	"github.com/evanqhuang/resume-cli/server"
	"github.com/spf13/cobra"
)

const (
	version = "0.1.0"

	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

var (
	resumePath  string
	outputFile  string
	jobDescFile string
	jobDescText string
	itemIDs     []string
	itemTags    []string
	serverPort  int
)

func main() {
	// Load .env file if present
	if err := config.LoadEnv(); err != nil {
		fmt.Fprintf(os.Stderr, "%sWarning: failed to load .env file: %v%s\n", colorYellow, err, colorReset)
	}

	rootCmd := &cobra.Command{
		Use:     "resume-cli",
		Short:   "A CLI tool for building tailored resumes",
		Long:    "Generate LaTeX resumes and match resume items to job descriptions using AI",
		Version: version,
	}

	// Default resume path
	defaultResumePath := "../resume.yaml"
	if abs, err := filepath.Abs(defaultResumePath); err == nil {
		defaultResumePath = abs
	}

	rootCmd.PersistentFlags().StringVarP(&resumePath, "resume", "r", defaultResumePath, "Path to resume.yaml file")

	// Add subcommands
	rootCmd.AddCommand(matchCmd())
	rootCmd.AddCommand(generateCmd())
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(serveCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%sError: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}
}

func serveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start HTTP server for resume API",
		Long:  "Start an HTTP server to serve the resume API for the frontend",
		RunE:  runServe,
	}

	cmd.Flags().IntVarP(&serverPort, "port", "p", 8080, "Port to run the server on")

	return cmd
}

func matchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "match",
		Short: "Analyze job description and score resume items",
		Long:  "Use AI to score resume items against a job description",
		RunE:  runMatch,
	}

	cmd.Flags().StringVarP(&jobDescFile, "file", "f", "", "Path to file containing job description")
	cmd.Flags().StringVarP(&jobDescText, "job", "j", "", "Job description text (inline)")

	return cmd
}

func generateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate resume PDF",
		Long:  "Generate a LaTeX resume and compile it to PDF",
		RunE:  runGenerate,
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "resume.pdf", "Output PDF file path")
	cmd.Flags().StringSliceVar(&itemIDs, "ids", []string{}, "Comma-separated list of item IDs to include")
	cmd.Flags().StringSliceVar(&itemTags, "tags", []string{}, "Comma-separated list of tags to filter items")

	return cmd
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all resume items with IDs and tags",
		Long:  "Display all resume items grouped by section with their IDs and tags",
		RunE:  runList,
	}
}

func runMatch(cmd *cobra.Command, args []string) error {
	// Get job description
	var jobDesc string
	if jobDescFile != "" {
		data, err := os.ReadFile(jobDescFile)
		if err != nil {
			return fmt.Errorf("failed to read job description file: %w", err)
		}
		jobDesc = string(data)
		if len(jobDesc) == 0 {
			return fmt.Errorf("job description file is empty")
		}
	} else if jobDescText != "" {
		jobDesc = jobDescText
	} else {
		return fmt.Errorf("either --file or --job must be specified")
	}

	// Validate API key is set
	if os.Getenv("OPENROUTER_API_KEY") == "" {
		return fmt.Errorf("OPENROUTER_API_KEY environment variable not set. Set it or add to .env file")
	}

	// Load resume
	fmt.Printf("%sLoading resume from: %s%s\n", colorCyan, resumePath, colorReset)
	r, err := resume.LoadResume(resumePath)
	if err != nil {
		return fmt.Errorf("failed to load resume: %w", err)
	}

	// Analyze with OpenRouter
	fmt.Printf("%sAnalyzing with AI...%s\n", colorYellow, colorReset)
	result, err := matching.AnalyzeJob(r, jobDesc)
	if err != nil {
		return fmt.Errorf("failed to analyze job: %w", err)
	}

	// Get all items and sort by score
	items := r.GetAllIDs()
	type scoredItem struct {
		item  resume.ItemWithID
		score float64
	}

	var scored []scoredItem
	for _, item := range items {
		if score, ok := result.Scores[item.ID]; ok {
			scored = append(scored, scoredItem{item: item, score: score})
		}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Display results
	fmt.Printf("\n%s=== Matching Results ===%s\n\n", colorGreen, colorReset)
	for _, s := range scored {
		scoreColor := getScoreColor(s.score)
		fmt.Printf("%s[%.0f]%s %s%s%s\n", scoreColor, s.score, colorReset, colorBlue, s.item.ID, colorReset)

		// Truncate text for display
		text := s.item.Text
		if len(text) > 100 {
			text = text[:97] + "..."
		}
		fmt.Printf("  %s\n", text)

		if len(s.item.Tags) > 0 {
			fmt.Printf("  %sTags:%s %s\n", colorPurple, colorReset, strings.Join(s.item.Tags, ", "))
		}
		fmt.Println()
	}

	return nil
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Validate resume file exists
	if _, err := os.Stat(resumePath); os.IsNotExist(err) {
		return fmt.Errorf("resume file not found: %s", resumePath)
	}

	// Load resume
	fmt.Printf("%sLoading resume from: %s%s\n", colorCyan, resumePath, colorReset)
	r, err := resume.LoadResume(resumePath)
	if err != nil {
		return fmt.Errorf("failed to load resume: %w", err)
	}

	// Determine which items to include
	var selectedIDs map[string]bool
	if len(itemIDs) > 0 {
		selectedIDs = r.FilterByIDs(itemIDs)
		fmt.Printf("%sFiltering by IDs: %s%s\n", colorYellow, strings.Join(itemIDs, ", "), colorReset)
		if len(selectedIDs) == 0 {
			return fmt.Errorf("no items found matching the specified IDs")
		}
	} else if len(itemTags) > 0 {
		selectedIDs = r.FilterByTags(itemTags)
		fmt.Printf("%sFiltering by tags: %s%s\n", colorYellow, strings.Join(itemTags, ", "), colorReset)
		if len(selectedIDs) == 0 {
			return fmt.Errorf("no items found matching the specified tags")
		}
	} else {
		selectedIDs = make(map[string]bool) // Empty map means include all
		fmt.Printf("%sIncluding all items%s\n", colorYellow, colorReset)
	}

	// Generate LaTeX
	fmt.Printf("%sGenerating LaTeX...%s\n", colorCyan, colorReset)
	latexContent, err := generator.GenerateLatex(r, selectedIDs)
	if err != nil {
		return fmt.Errorf("failed to generate LaTeX: %w", err)
	}

	// Write LaTeX to temporary file
	texFile := strings.TrimSuffix(outputFile, ".pdf") + ".tex"
	if err := os.WriteFile(texFile, []byte(latexContent), 0644); err != nil {
		return fmt.Errorf("failed to write LaTeX file: %w", err)
	}
	fmt.Printf("%sWrote LaTeX to: %s%s\n", colorGreen, texFile, colorReset)

	// Find xelatex
	xelatexPath, err := generator.FindXelatex()
	if err != nil {
		return fmt.Errorf("xelatex not found in PATH. Please install a TeX distribution (MacTeX, TeX Live, or MiKTeX)")
	}

	// Compile to PDF
	fmt.Printf("%sCompiling PDF with xelatex...%s\n", colorCyan, colorReset)
	if err := compilePDF(texFile, xelatexPath); err != nil {
		return fmt.Errorf("failed to compile PDF: %w", err)
	}

	// Clean up intermediate files
	cleanupFiles(texFile)

	fmt.Printf("%s✓ Successfully generated: %s%s\n", colorGreen, outputFile, colorReset)
	return nil
}

func runList(cmd *cobra.Command, args []string) error {
	// Validate resume file exists
	if _, err := os.Stat(resumePath); os.IsNotExist(err) {
		return fmt.Errorf("resume file not found: %s", resumePath)
	}

	// Load resume
	r, err := resume.LoadResume(resumePath)
	if err != nil {
		return fmt.Errorf("failed to load resume: %w", err)
	}

	items := r.GetAllIDs()

	// Group by section
	sections := make(map[string][]resume.ItemWithID)
	for _, item := range items {
		sections[item.Section] = append(sections[item.Section], item)
	}

	// Display in order: Experience, Projects, Leadership
	sectionOrder := []string{"Experience", "Projects", "Leadership"}

	for _, section := range sectionOrder {
		items, ok := sections[section]
		if !ok || len(items) == 0 {
			continue
		}

		fmt.Printf("\n%s=== %s ===%s\n\n", colorGreen, section, colorReset)

		currentCategory := ""
		for _, item := range items {
			if item.Category != "" && item.Category != currentCategory {
				currentCategory = item.Category
				fmt.Printf("%s%s%s\n", colorYellow, currentCategory, colorReset)
			}

			fmt.Printf("  %s%s%s\n", colorBlue, item.ID, colorReset)

			text := item.Text
			if len(text) > 100 {
				text = text[:97] + "..."
			}
			fmt.Printf("    %s\n", text)

			if len(item.Tags) > 0 {
				fmt.Printf("    %sTags:%s %s\n", colorPurple, colorReset, strings.Join(item.Tags, ", "))
			}
			fmt.Println()
		}
	}

	return nil
}

func runServe(cmd *cobra.Command, args []string) error {
	// Validate resume file exists
	if _, err := os.Stat(resumePath); os.IsNotExist(err) {
		return fmt.Errorf("resume file not found: %s", resumePath)
	}

	fmt.Printf("%sStarting server with resume: %s%s\n", colorCyan, resumePath, colorReset)
	fmt.Printf("%sServer will be available at: %shttp://localhost:%d%s\n", colorGreen, colorWhite, serverPort, colorReset)

	return server.Start(resumePath, serverPort)
}

func compilePDF(texFile, xelatexPath string) error {
	// Get the directory containing the tex file for output
	outputDir := filepath.Dir(texFile)
	if outputDir == "" || outputDir == "." {
		outputDir, _ = os.Getwd()
	}
	// Convert to absolute path if needed
	absTexFile, err := filepath.Abs(texFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	absOutputDir := filepath.Dir(absTexFile)

	// Run xelatex twice for proper formatting
	for i := 0; i < 2; i++ {
		cmd := exec.Command(xelatexPath, "-interaction=nonstopmode", "-output-directory="+absOutputDir, absTexFile)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("xelatex failed: %w\nOutput: %s", err, string(output))
		}
	}
	return nil
}

func cleanupFiles(texFile string) {
	base := strings.TrimSuffix(texFile, ".tex")
	extensions := []string{".aux", ".log", ".out"}

	for _, ext := range extensions {
		os.Remove(base + ext)
	}
}

func getScoreColor(score float64) string {
	switch {
	case score >= 90:
		return colorGreen
	case score >= 70:
		return colorYellow
	case score >= 50:
		return colorCyan
	default:
		return colorRed
	}
}
