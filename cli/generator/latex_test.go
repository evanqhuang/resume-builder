package generator

import (
	"strings"
	"testing"

	"github.com/evanqhuang/resume-cli/resume"
)

func TestEscapeLaTeX(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello & World", `Hello \& World`},
		{"Cost: $100", `Cost: \$100`},
		{"50% complete", `50\% complete`},
		{"#hashtag", `\#hashtag`},
		{"file_name", `file\_name`},
		{"{braces}", `\{braces\}`},
		{"~tilde", `\textasciitilde{}tilde`},
		{"x^2", `x\textasciicircum{}2`},
		{"normal text", "normal text"},
	}

	for _, tt := range tests {
		result := escapeLaTeX(tt.input)
		if result != tt.expected {
			t.Errorf("escapeLaTeX(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestGenerateLatex(t *testing.T) {
	r := &resume.Resume{
		Contact: resume.ContactInfo{
			Name:     "Test User",
			Email:    "test@example.com",
			Phone:    "123-456-7890",
			LinkedIn: "linkedin.com/in/test",
			GitHub:   "github.com/test",
		},
		Education: resume.EducationEntry{
			Institution: "Test University",
			Location:    "Test City",
			Degree:      "B.S. Computer Science",
			GPA:         "3.5/4.0",
		},
		Skills: resume.Skills{
			Languages: []resume.SkillItem{
				{Name: "Go", Tags: []string{"go"}},
			},
			Frameworks: []resume.SkillItem{
				{Name: "Docker", Tags: []string{"docker"}},
			},
			Cloud: []resume.SkillItem{
				{Name: "AWS", Tags: []string{"aws"}},
			},
		},
		Experience: []resume.ExperienceEntry{
			{
				ID:        "test-exp",
				Title:     "Software Engineer",
				Company:   "Test Company",
				Location:  "Test City",
				StartDate: "Jan 2020",
				EndDate:   "Present",
				Bullets: []resume.Bullet{
					{ID: "bullet-1", Text: "Did something cool with Go & Docker"},
				},
			},
		},
	}

	latex, err := GenerateLatex(r, nil)
	if err != nil {
		t.Fatalf("GenerateLatex failed: %v", err)
	}

	// Check that document class is present
	if !strings.Contains(latex, `\documentclass`) {
		t.Error("LaTeX output missing document class")
	}

	// Check that contact info is present
	if !strings.Contains(latex, "Test User") {
		t.Error("LaTeX output missing contact name")
	}

	// Check that special characters are escaped
	if !strings.Contains(latex, `\&`) {
		t.Error("LaTeX output not escaping & character")
	}

	// Check that experience is included
	if !strings.Contains(latex, "Test Company") {
		t.Error("LaTeX output missing experience")
	}
}

func TestGenerateLatexWithFiltering(t *testing.T) {
	r := &resume.Resume{
		Contact: resume.ContactInfo{
			Name: "Test User",
		},
		Education: resume.EducationEntry{
			Institution: "Test University",
		},
		Skills: resume.Skills{
			Languages: []resume.SkillItem{{Name: "Go"}},
		},
		Experience: []resume.ExperienceEntry{
			{
				ID:       "exp-1",
				Title:    "Engineer",
				Company:  "Company A",
				Bullets: []resume.Bullet{
					{ID: "bullet-1", Text: "Bullet 1"},
					{ID: "bullet-2", Text: "Bullet 2"},
				},
			},
		},
	}

	// Filter to include only bullet-1
	selectedIDs := map[string]bool{"bullet-1": true}
	latex, err := GenerateLatex(r, selectedIDs)
	if err != nil {
		t.Fatalf("GenerateLatex failed: %v", err)
	}

	// Should include bullet-1
	if !strings.Contains(latex, "Bullet 1") {
		t.Error("LaTeX output missing filtered bullet")
	}

	// Should NOT include bullet-2
	if strings.Contains(latex, "Bullet 2") {
		t.Error("LaTeX output should not include unfiltered bullet")
	}
}
