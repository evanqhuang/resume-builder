package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/evanqhuang/resume-cli/resume"
)

//go:embed templates/modern.tex
var modernTemplate string

// TemplateData holds data for LaTeX template
type TemplateData struct {
	Contact       resume.ContactInfo
	Summary       string
	Education     resume.EducationEntry
	Skills        resume.Skills
	Experience    []ExperienceData
	Projects      []ProjectData
	Leadership    []string
	IncludeSkills bool
}

// ExperienceData holds experience data for template
type ExperienceData struct {
	Title     string
	Company   string
	Location  string
	StartDate string
	EndDate   string
	Bullets   []string
}

// ProjectData holds project data for template
type ProjectData struct {
	Title        string
	Technologies string
	GitHub       string
	Bullets      []string
}

// GenerateLatex generates LaTeX source from resume data
func GenerateLatex(r *resume.Resume, selectedIDs map[string]bool) (string, error) {
	data := prepareTemplateData(r, selectedIDs)

	tmpl, err := template.New("resume").Funcs(template.FuncMap{
		"escape": escapeLaTeX,
	}).Parse(modernTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func prepareTemplateData(r *resume.Resume, selectedIDs map[string]bool) TemplateData {
	includeAll := len(selectedIDs) == 0

	data := TemplateData{
		Contact:       r.Contact,
		Summary:       r.Summary,
		Education:     r.Education,
		Skills:        r.Skills,
		IncludeSkills: true, // Always include skills section
	}

	// Process experience entries
	for _, exp := range r.Experience {
		var bullets []string
		for _, bullet := range exp.Bullets {
			if includeAll || selectedIDs[bullet.ID] {
				bullets = append(bullets, bullet.Text)
			}
		}
		// Only include experience entries that have at least one bullet
		if len(bullets) > 0 {
			data.Experience = append(data.Experience, ExperienceData{
				Title:     exp.Title,
				Company:   exp.Company,
				Location:  exp.Location,
				StartDate: exp.StartDate,
				EndDate:   exp.EndDate,
				Bullets:   bullets,
			})
		}
	}

	// Process project entries
	for _, proj := range r.Projects {
		var bullets []string
		for _, bullet := range proj.Bullets {
			if includeAll || selectedIDs[bullet.ID] {
				bullets = append(bullets, bullet.Text)
			}
		}
		// Only include projects that have at least one bullet
		if len(bullets) > 0 {
			data.Projects = append(data.Projects, ProjectData{
				Title:        proj.Title,
				Technologies: proj.Technologies,
				GitHub:       proj.GitHub,
				Bullets:      bullets,
			})
		}
	}

	// Process leadership entries
	for _, lead := range r.Leadership {
		if includeAll || selectedIDs[lead.ID] {
			data.Leadership = append(data.Leadership, lead.Text)
		}
	}

	return data
}

// escapeLaTeX escapes special LaTeX characters
func escapeLaTeX(s string) string {
	// Must escape backslash first to avoid double-escaping
	result := strings.ReplaceAll(s, `\`, `\textbackslash{}`)

	replacements := []struct {
		old, new string
	}{
		{"&", `\&`},
		{"%", `\%`},
		{"$", `\$`},
		{"#", `\#`},
		{"_", `\_`},
		{"{", `\{`},
		{"}", `\}`},
		{"~", `\textasciitilde{}`},
		{"^", `\textasciicircum{}`},
		{"<", `\textless{}`},
		{">", `\textgreater{}`},
		{"→", `$\rightarrow$`},
	}

	for _, r := range replacements {
		result = strings.ReplaceAll(result, r.old, r.new)
	}
	return result
}

// Common xelatex installation paths to check
var xelatexPaths = []string{
	"/Library/TeX/texbin/xelatex",                          // MacTeX on macOS
	"/usr/local/texlive/2025/bin/universal-darwin/xelatex", // TeX Live 2025 macOS
	"/usr/local/texlive/2024/bin/universal-darwin/xelatex", // TeX Live 2024 macOS
	"/usr/local/texlive/2023/bin/universal-darwin/xelatex", // TeX Live 2023 macOS
	"/usr/local/texlive/2025/bin/x86_64-linux/xelatex",    // TeX Live 2025 Linux
	"/usr/local/texlive/2024/bin/x86_64-linux/xelatex",    // TeX Live 2024 Linux
	"/usr/local/texlive/2023/bin/x86_64-linux/xelatex",    // TeX Live 2023 Linux
	"/usr/bin/xelatex",                                     // System installation
}

// FindXelatex locates the xelatex executable
func FindXelatex() (string, error) {
	// First try PATH
	if path, err := exec.LookPath("xelatex"); err == nil {
		return path, nil
	}

	// Check common installation paths
	for _, path := range xelatexPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("xelatex not found. Install LaTeX (e.g., 'brew install --cask mactex' on macOS)")
}

// GeneratePDF generates a PDF from resume data and returns the bytes
func GeneratePDF(r *resume.Resume, selectedIDs map[string]bool) ([]byte, error) {
	// Generate LaTeX content
	latexContent, err := GenerateLatex(r, selectedIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to generate LaTeX: %w", err)
	}

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "resume-pdf-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write LaTeX to temp file
	texFile := filepath.Join(tmpDir, "resume.tex")
	if err := os.WriteFile(texFile, []byte(latexContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write LaTeX file: %w", err)
	}

	// Find xelatex
	xelatexPath, err := FindXelatex()
	if err != nil {
		return nil, err
	}

	// Compile to PDF (run twice for proper formatting)
	for i := 0; i < 2; i++ {
		cmd := exec.Command(xelatexPath, "-interaction=nonstopmode", "-output-directory="+tmpDir, texFile)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("xelatex failed: %w\nOutput: %s", err, string(output))
		}
	}

	// Read the generated PDF
	pdfFile := filepath.Join(tmpDir, "resume.pdf")
	pdfBytes, err := os.ReadFile(pdfFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF: %w", err)
	}

	return pdfBytes, nil
}
