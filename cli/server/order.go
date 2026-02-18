package server

import (
	"os"

	"github.com/evanqhuang/resume-cli/resume"
	"gopkg.in/yaml.v3"
)

// SectionOrder represents custom ordering for resume sections
type SectionOrder struct {
	Experience []string `yaml:"experience" json:"experience"`
	Projects   []string `yaml:"projects" json:"projects"`
	Leadership []string `yaml:"leadership" json:"leadership"`
}

// LoadOrder reads order.yaml or returns default order from resume
func LoadOrder(orderPath string, r *resume.Resume) (*SectionOrder, error) {
	data, err := os.ReadFile(orderPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default order from resume
			return GetDefaultOrder(r), nil
		}
		return nil, err
	}

	var order SectionOrder
	if err := yaml.Unmarshal(data, &order); err != nil {
		return nil, err
	}
	return &order, nil
}

// SaveOrder writes order to order.yaml
func SaveOrder(orderPath string, order *SectionOrder) error {
	data, err := yaml.Marshal(order)
	if err != nil {
		return err
	}
	return os.WriteFile(orderPath, data, 0644)
}

// GetDefaultOrder extracts IDs in their original YAML order
func GetDefaultOrder(r *resume.Resume) *SectionOrder {
	order := &SectionOrder{
		Experience: make([]string, len(r.Experience)),
		Projects:   make([]string, len(r.Projects)),
		Leadership: make([]string, len(r.Leadership)),
	}

	for i, exp := range r.Experience {
		order.Experience[i] = exp.ID
	}
	for i, proj := range r.Projects {
		order.Projects[i] = proj.ID
	}
	for i, lead := range r.Leadership {
		order.Leadership[i] = lead.ID
	}

	return order
}

// PartialSectionOrder allows updating a single section's order.
type PartialSectionOrder struct {
	Experience *[]string `yaml:"experience,omitempty" json:"experience,omitempty"`
	Projects   *[]string `yaml:"projects,omitempty" json:"projects,omitempty"`
	Leadership *[]string `yaml:"leadership,omitempty" json:"leadership,omitempty"`
}

// MergeOrder applies a partial order update to an existing SectionOrder.
func MergeOrder(existing *SectionOrder, partial *PartialSectionOrder) {
	if partial.Experience != nil {
		existing.Experience = *partial.Experience
	}
	if partial.Projects != nil {
		existing.Projects = *partial.Projects
	}
	if partial.Leadership != nil {
		existing.Leadership = *partial.Leadership
	}
}
