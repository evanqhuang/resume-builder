package resume

import (
	"os"
	"testing"
)

func TestLoadResume(t *testing.T) {
	// Create a temporary YAML file
	content := `contact:
  name: Test User
  email: test@example.com
  phone: 123-456-7890
  linkedin: linkedin.com/in/test
  github: github.com/test

education:
  institution: Test University
  location: Test City
  degree: B.S. Computer Science
  gpa: "3.5/4.0"

skills:
  languages:
    - name: Go
      tags: [go, backend]
  frameworks:
    - name: Docker
      tags: [docker, containers]
  cloud:
    - name: AWS
      tags: [aws, cloud]

experience:
  - id: test-exp
    title: Software Engineer
    company: Test Company
    location: Test City
    start_date: Jan 2020
    end_date: Present
    tags: [test]
    bullets:
      - id: test-bullet
        text: Test accomplishment
        tags: [go, testing]

projects: []
leadership: []
`

	tmpFile, err := os.CreateTemp("", "resume-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test loading
	resume, err := LoadResume(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load resume: %v", err)
	}

	// Verify data
	if resume.Contact.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", resume.Contact.Name)
	}

	if len(resume.Experience) != 1 {
		t.Errorf("Expected 1 experience entry, got %d", len(resume.Experience))
	}

	if resume.Experience[0].ID != "test-exp" {
		t.Errorf("Expected experience ID 'test-exp', got '%s'", resume.Experience[0].ID)
	}

	if len(resume.Experience[0].Bullets) != 1 {
		t.Errorf("Expected 1 bullet, got %d", len(resume.Experience[0].Bullets))
	}

	if resume.Experience[0].Bullets[0].ID != "test-bullet" {
		t.Errorf("Expected bullet ID 'test-bullet', got '%s'", resume.Experience[0].Bullets[0].ID)
	}
}

func TestLoadResumeFileNotFound(t *testing.T) {
	_, err := LoadResume("/nonexistent/file.yaml")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}
