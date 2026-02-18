package server

import (
	"github.com/evanqhuang/resume-cli/resume"
)

// TransformedResume is the response format expected by the frontend
type TransformedResume struct {
	Contact    resume.ContactInfo       `json:"contact"`
	Summary    string                   `json:"summary"`
	Education  resume.EducationEntry    `json:"education"`
	Skills     []SkillCategory          `json:"skills"`
	Experience []TransformedExperience  `json:"experience"`
	Projects   []TransformedProject     `json:"projects"`
	Leadership []TransformedLeadership  `json:"leadership"`
}

// SkillCategory represents a skill category with selected items
type SkillCategory struct {
	Category string              `json:"category"`
	Items    []TransformedSkill  `json:"items"`
}

// TransformedSkill is a skill item with a selected flag
type TransformedSkill struct {
	Name     string   `json:"name"`
	Tags     []string `json:"tags"`
	Selected bool     `json:"selected"`
}

// TransformedExperience is an experience entry with selected flags
type TransformedExperience struct {
	ID        string              `json:"id"`
	Title     string              `json:"title"`
	Company   string              `json:"company"`
	Location  string              `json:"location"`
	StartDate string              `json:"start_date"`
	EndDate   string              `json:"end_date"`
	Tags      []string            `json:"tags"`
	Bullets   []TransformedBullet `json:"bullets"`
	Selected  bool                `json:"selected"`
}

// TransformedProject is a project entry with selected flags
type TransformedProject struct {
	ID           string              `json:"id"`
	Title        string              `json:"title"`
	Technologies string              `json:"technologies"`
	GitHub       string              `json:"github,omitempty"`
	Tags         []string            `json:"tags"`
	Bullets      []TransformedBullet `json:"bullets"`
	Selected     bool                `json:"selected"`
}

// TransformedBullet is a bullet with a selected flag
type TransformedBullet struct {
	ID       string   `json:"id"`
	Text     string   `json:"text"`
	Tags     []string `json:"tags"`
	Selected bool     `json:"selected"`
}

// TransformedLeadership is a leadership entry with a selected flag
type TransformedLeadership struct {
	ID       string   `json:"id"`
	Text     string   `json:"text"`
	Tags     []string `json:"tags"`
	Selected bool     `json:"selected"`
}

// TransformResume converts a Resume to the frontend-expected format
func TransformResume(r *resume.Resume) *TransformedResume {
	return &TransformedResume{
		Contact:    r.Contact,
		Summary:    r.Summary,
		Education:  r.Education,
		Skills:     transformSkills(r.Skills),
		Experience: transformExperience(r.Experience),
		Projects:   transformProjects(r.Projects),
		Leadership: transformLeadership(r.Leadership),
	}
}

func transformSkills(skills resume.Skills) []SkillCategory {
	return []SkillCategory{
		{
			Category: "languages",
			Items:    transformSkillItems(skills.Languages),
		},
		{
			Category: "frameworks",
			Items:    transformSkillItems(skills.Frameworks),
		},
		{
			Category: "cloud",
			Items:    transformSkillItems(skills.Cloud),
		},
	}
}

func transformSkillItems(items []resume.SkillItem) []TransformedSkill {
	result := make([]TransformedSkill, len(items))
	for i, item := range items {
		tags := item.Tags
		if tags == nil {
			tags = []string{}
		}
		result[i] = TransformedSkill{
			Name:     item.Name,
			Tags:     tags,
			Selected: true,
		}
	}
	return result
}

func transformExperience(entries []resume.ExperienceEntry) []TransformedExperience {
	result := make([]TransformedExperience, len(entries))
	for i, entry := range entries {
		tags := entry.Tags
		if tags == nil {
			tags = []string{}
		}
		result[i] = TransformedExperience{
			ID:        entry.ID,
			Title:     entry.Title,
			Company:   entry.Company,
			Location:  entry.Location,
			StartDate: entry.StartDate,
			EndDate:   entry.EndDate,
			Tags:      tags,
			Bullets:   transformBullets(entry.Bullets),
			Selected:  true,
		}
	}
	return result
}

func transformProjects(entries []resume.ProjectEntry) []TransformedProject {
	result := make([]TransformedProject, len(entries))
	for i, entry := range entries {
		tags := entry.Tags
		if tags == nil {
			tags = []string{}
		}
		result[i] = TransformedProject{
			ID:           entry.ID,
			Title:        entry.Title,
			Technologies: entry.Technologies,
			GitHub:       entry.GitHub,
			Tags:         tags,
			Bullets:      transformBullets(entry.Bullets),
			Selected:     true,
		}
	}
	return result
}

func transformBullets(bullets []resume.Bullet) []TransformedBullet {
	result := make([]TransformedBullet, len(bullets))
	for i, bullet := range bullets {
		tags := bullet.Tags
		if tags == nil {
			tags = []string{}
		}
		result[i] = TransformedBullet{
			ID:       bullet.ID,
			Text:     bullet.Text,
			Tags:     tags,
			Selected: true,
		}
	}
	return result
}

func transformLeadership(entries []resume.LeadershipEntry) []TransformedLeadership {
	result := make([]TransformedLeadership, len(entries))
	for i, entry := range entries {
		tags := entry.Tags
		if tags == nil {
			tags = []string{}
		}
		result[i] = TransformedLeadership{
			ID:       entry.ID,
			Text:     entry.Text,
			Tags:     tags,
			Selected: true,
		}
	}
	return result
}
