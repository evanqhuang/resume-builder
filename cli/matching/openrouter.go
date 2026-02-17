package matching

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/evanqhuang/resume-cli/resume"
)

const openRouterURL = "https://openrouter.ai/api/v1/chat/completions"

// MatchResult contains scored items from the resume (for CLI usage)
type MatchResult struct {
	Scores map[string]float64 // ID -> score (0-100)
}

// JobAnalysisResult contains the full analysis response (for API usage)
type JobAnalysisResult struct {
	Keywords       []string           `json:"keywords"`
	Scores         map[string]float64 `json:"scores"`
	SuggestedItems []string           `json:"suggested_items"`
}

// ScoredItem represents an item with its relevance score
type ScoredItem struct {
	ID    string
	Text  string
	Tags  []string
	Score float64
}

// OpenRouterRequest represents the API request structure
type OpenRouterRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterResponse represents the API response
type OpenRouterResponse struct {
	Choices []Choice `json:"choices"`
}

// Choice represents a response choice
type Choice struct {
	Message Message `json:"message"`
}

// AnalyzeJob calls OpenRouter API to score resume items against job description
func AnalyzeJob(r *resume.Resume, jobDescription string) (*MatchResult, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY environment variable not set")
	}

	model := os.Getenv("OPENROUTER_MODEL")
	if model == "" {
		model = "anthropic/claude-sonnet-4"
	}

	// Build prompt with all resume items
	prompt := buildMatchingPrompt(r, jobDescription)

	reqBody := OpenRouterRequest{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", openRouterURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp OpenRouterResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned")
	}

	// Parse the JSON scores from the response
	scores, err := parseScores(apiResp.Choices[0].Message.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse scores: %w", err)
	}

	return &MatchResult{Scores: scores}, nil
}

func buildMatchingPrompt(r *resume.Resume, jobDescription string) string {
	var prompt bytes.Buffer

	prompt.WriteString("You are analyzing a resume against a job description. ")
	prompt.WriteString("Score each resume item's relevance to the job on a scale of 0-100, where:\n")
	prompt.WriteString("- 90-100: Highly relevant, directly addresses key requirements\n")
	prompt.WriteString("- 70-89: Relevant, demonstrates related skills\n")
	prompt.WriteString("- 50-69: Somewhat relevant, transferable skills\n")
	prompt.WriteString("- 30-49: Tangentially related\n")
	prompt.WriteString("- 0-29: Not relevant\n\n")

	prompt.WriteString("Job Description:\n")
	prompt.WriteString(jobDescription)
	prompt.WriteString("\n\nResume Items:\n\n")

	// Experience bullets
	for _, exp := range r.Experience {
		for _, bullet := range exp.Bullets {
			fmt.Fprintf(&prompt, "ID: %s\nText: %s\nTags: %v\n\n", bullet.ID, bullet.Text, bullet.Tags)
		}
	}

	// Project bullets
	for _, proj := range r.Projects {
		for _, bullet := range proj.Bullets {
			fmt.Fprintf(&prompt, "ID: %s\nText: %s\nTags: %v\n\n", bullet.ID, bullet.Text, bullet.Tags)
		}
	}

	// Leadership entries
	for _, lead := range r.Leadership {
		fmt.Fprintf(&prompt, "ID: %s\nText: %s\nTags: %v\n\n", lead.ID, lead.Text, lead.Tags)
	}

	prompt.WriteString("\nReturn your response as a JSON object with this exact format:\n")
	prompt.WriteString("{\n")
	prompt.WriteString("  \"scores\": {\n")
	prompt.WriteString("    \"item-id-1\": 95,\n")
	prompt.WriteString("    \"item-id-2\": 82,\n")
	prompt.WriteString("    ...\n")
	prompt.WriteString("  }\n")
	prompt.WriteString("}\n")
	prompt.WriteString("\nOnly include the JSON object in your response, no other text.")

	return prompt.String()
}

func parseScores(content string) (map[string]float64, error) {
	// Try to extract JSON from the content
	// Sometimes the model wraps it in markdown code blocks
	start := bytes.Index([]byte(content), []byte("{"))
	end := bytes.LastIndex([]byte(content), []byte("}"))

	if start == -1 || end == -1 || start >= end {
		return nil, fmt.Errorf("no JSON object found in response")
	}

	jsonContent := content[start : end+1]

	var result struct {
		Scores map[string]float64 `json:"scores"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return result.Scores, nil
}

// AnalyzeJobForAPI calls OpenRouter API for the web API (returns keywords, scores, suggested_items)
func AnalyzeJobForAPI(r *resume.Resume, jobTitle, company, jobDescription string) (*JobAnalysisResult, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY environment variable not set")
	}

	model := os.Getenv("OPENROUTER_MODEL")
	if model == "" {
		model = "anthropic/claude-sonnet-4"
	}

	// Collect all item IDs
	allItemIDs := collectAllItemIDs(r)

	// Build prompt
	prompt := buildAPIAnalysisPrompt(r, jobTitle, company, jobDescription, allItemIDs)

	reqBody := OpenRouterRequest{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", openRouterURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp OpenRouterResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned")
	}

	// Parse the full analysis response
	return parseAnalysisResponse(apiResp.Choices[0].Message.Content, allItemIDs)
}

func collectAllItemIDs(r *resume.Resume) []string {
	var ids []string

	// Skills
	for _, skill := range r.Skills.Languages {
		ids = append(ids, "skill-"+strings.ToLower(strings.ReplaceAll(skill.Name, " ", "-")))
	}
	for _, skill := range r.Skills.Frameworks {
		ids = append(ids, "skill-"+strings.ToLower(strings.ReplaceAll(skill.Name, " ", "-")))
	}
	for _, skill := range r.Skills.Cloud {
		ids = append(ids, "skill-"+strings.ToLower(strings.ReplaceAll(skill.Name, " ", "-")))
	}

	// Experience
	for _, exp := range r.Experience {
		ids = append(ids, exp.ID)
		for _, bullet := range exp.Bullets {
			ids = append(ids, bullet.ID)
		}
	}

	// Projects
	for _, proj := range r.Projects {
		ids = append(ids, proj.ID)
		for _, bullet := range proj.Bullets {
			ids = append(ids, bullet.ID)
		}
	}

	// Leadership
	for _, lead := range r.Leadership {
		ids = append(ids, lead.ID)
	}

	return ids
}

func buildAPIAnalysisPrompt(r *resume.Resume, jobTitle, company, jobDescription string, allItemIDs []string) string {
	var prompt bytes.Buffer

	prompt.WriteString("You are a resume optimization expert. Analyze this job description and score each resume item for relevance.\n\n")
	fmt.Fprintf(&prompt, "Job Title: %s\n", jobTitle)
	fmt.Fprintf(&prompt, "Company: %s\n\n", company)
	prompt.WriteString("Job Description:\n")
	prompt.WriteString(jobDescription)
	prompt.WriteString("\n\nResume Summary:\n")

	// Skills
	prompt.WriteString("SKILLS:\n")
	for _, skill := range r.Skills.Languages {
		skillID := "skill-" + strings.ToLower(strings.ReplaceAll(skill.Name, " ", "-"))
		tags := strings.Join(skill.Tags, ", ")
		if len(skill.Tags) > 5 {
			tags = strings.Join(skill.Tags[:5], ", ")
		}
		fmt.Fprintf(&prompt, "  %s: %s (%s)\n", skillID, skill.Name, tags)
	}
	for _, skill := range r.Skills.Frameworks {
		skillID := "skill-" + strings.ToLower(strings.ReplaceAll(skill.Name, " ", "-"))
		tags := strings.Join(skill.Tags, ", ")
		if len(skill.Tags) > 5 {
			tags = strings.Join(skill.Tags[:5], ", ")
		}
		fmt.Fprintf(&prompt, "  %s: %s (%s)\n", skillID, skill.Name, tags)
	}
	for _, skill := range r.Skills.Cloud {
		skillID := "skill-" + strings.ToLower(strings.ReplaceAll(skill.Name, " ", "-"))
		tags := strings.Join(skill.Tags, ", ")
		if len(skill.Tags) > 5 {
			tags = strings.Join(skill.Tags[:5], ", ")
		}
		fmt.Fprintf(&prompt, "  %s: %s (%s)\n", skillID, skill.Name, tags)
	}

	// Experience
	prompt.WriteString("\nEXPERIENCE:\n")
	for _, exp := range r.Experience {
		fmt.Fprintf(&prompt, "  %s: %s at %s\n", exp.ID, exp.Title, exp.Company)
		for i, bullet := range exp.Bullets {
			if i >= 3 {
				break
			}
			text := bullet.Text
			if len(text) > 100 {
				text = text[:100] + "..."
			}
			fmt.Fprintf(&prompt, "    %s: %s\n", bullet.ID, text)
		}
	}

	// Projects
	prompt.WriteString("\nPROJECTS:\n")
	for _, proj := range r.Projects {
		fmt.Fprintf(&prompt, "  %s: %s\n", proj.ID, proj.Title)
		for i, bullet := range proj.Bullets {
			if i >= 2 {
				break
			}
			text := bullet.Text
			if len(text) > 100 {
				text = text[:100] + "..."
			}
			fmt.Fprintf(&prompt, "    %s: %s\n", bullet.ID, text)
		}
	}

	prompt.WriteString("\nAvailable Item IDs:\n")
	prompt.WriteString(strings.Join(allItemIDs, ", "))

	prompt.WriteString("\n\nProvide a JSON response with:\n")
	prompt.WriteString("1. \"keywords\": array of 10-15 key technical skills/terms from the job description\n")
	prompt.WriteString("2. \"scores\": object mapping each item ID to a relevance score (0-100)\n")
	prompt.WriteString("3. \"suggested_items\": array of item IDs you recommend including (score >= 60)\n\n")
	prompt.WriteString("Focus on:\n")
	prompt.WriteString("- Technical skills match\n")
	prompt.WriteString("- Domain/industry relevance\n")
	prompt.WriteString("- Impact and achievements that align with job requirements\n")
	prompt.WriteString("- Keywords and terminology overlap\n\n")
	prompt.WriteString("Return ONLY valid JSON, no other text.\n\n")
	prompt.WriteString("Example format:\n")
	prompt.WriteString("{\n")
	prompt.WriteString("  \"keywords\": [\"python\", \"distributed systems\", \"aws\"],\n")
	prompt.WriteString("  \"scores\": {\n")
	prompt.WriteString("    \"cap1-event-driven-transaction-processing\": 95,\n")
	prompt.WriteString("    \"skill-python\": 85\n")
	prompt.WriteString("  },\n")
	prompt.WriteString("  \"suggested_items\": [\"cap1-event-driven-transaction-processing\", \"skill-python\"]\n")
	prompt.WriteString("}\n")

	return prompt.String()
}

func parseAnalysisResponse(content string, allItemIDs []string) (*JobAnalysisResult, error) {
	// Strip markdown code blocks if present
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```json") {
		content = content[7:]
	}
	if strings.HasPrefix(content, "```") {
		content = content[3:]
	}
	if strings.HasSuffix(content, "```") {
		content = content[:len(content)-3]
	}
	content = strings.TrimSpace(content)

	// Try to extract JSON
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")

	if start == -1 || end == -1 || start >= end {
		return nil, fmt.Errorf("no JSON object found in response")
	}

	jsonContent := content[start : end+1]

	var data struct {
		Keywords       []string           `json:"keywords"`
		Scores         map[string]float64 `json:"scores"`
		SuggestedItems []string           `json:"suggested_items"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Build set of valid IDs for filtering
	validIDs := make(map[string]bool)
	for _, id := range allItemIDs {
		validIDs[id] = true
	}

	// Filter scores to only valid IDs
	validScores := make(map[string]float64)
	for k, v := range data.Scores {
		if validIDs[k] {
			validScores[k] = v
		}
	}

	// Filter suggested items to only valid IDs
	var validSuggested []string
	for _, item := range data.SuggestedItems {
		if validIDs[item] {
			validSuggested = append(validSuggested, item)
		}
	}

	return &JobAnalysisResult{
		Keywords:       data.Keywords,
		Scores:         validScores,
		SuggestedItems: validSuggested,
	}, nil
}
