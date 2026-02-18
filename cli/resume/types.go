package resume

// Resume represents the complete resume data structure
type Resume struct {
	Contact    ContactInfo         `yaml:"contact"`
	Summary    string              `yaml:"summary" json:"summary"`
	Education  EducationEntry      `yaml:"education"`
	Skills     Skills              `yaml:"skills"`
	Experience []ExperienceEntry   `yaml:"experience"`
	Projects   []ProjectEntry      `yaml:"projects"`
	Leadership []LeadershipEntry   `yaml:"leadership"`
}

// ContactInfo holds personal contact information
type ContactInfo struct {
	Name     string `yaml:"name" json:"name"`
	Location string `yaml:"location" json:"location"`
	Email    string `yaml:"email" json:"email"`
	Phone    string `yaml:"phone" json:"phone"`
	LinkedIn string `yaml:"linkedin" json:"linkedin"`
	GitHub   string `yaml:"github" json:"github"`
}

// EducationEntry represents education details
type EducationEntry struct {
	Institution string `yaml:"institution" json:"institution"`
	Location    string `yaml:"location" json:"location"`
	Degree      string `yaml:"degree" json:"degree"`
	Minor       string `yaml:"minor" json:"minor"`
	GPA         string `yaml:"gpa" json:"gpa"`
	Honors      string `yaml:"honors" json:"honors"`
}

// Skills contains all skill categories
type Skills struct {
	Languages  []SkillItem `yaml:"languages" json:"languages"`
	Frameworks []SkillItem `yaml:"frameworks" json:"frameworks"`
	Cloud      []SkillItem `yaml:"cloud" json:"cloud"`
}

// SkillItem represents a single skill with tags
type SkillItem struct {
	Name string   `yaml:"name" json:"name"`
	Tags []string `yaml:"tags" json:"tags"`
}

// ExperienceEntry represents a work experience
type ExperienceEntry struct {
	ID        string   `yaml:"id"`
	Title     string   `yaml:"title"`
	Company   string   `yaml:"company"`
	Location  string   `yaml:"location"`
	StartDate string   `yaml:"start_date"`
	EndDate   string   `yaml:"end_date"`
	Tags      []string `yaml:"tags"`
	Bullets   []Bullet `yaml:"bullets"`
}

// ProjectEntry represents a project
type ProjectEntry struct {
	ID           string   `yaml:"id"`
	Title        string   `yaml:"title"`
	Technologies string   `yaml:"technologies"`
	GitHub       string   `yaml:"github,omitempty"`
	Tags         []string `yaml:"tags"`
	Bullets      []Bullet `yaml:"bullets"`
}

// Bullet represents a single bullet point with ID and tags
type Bullet struct {
	ID   string   `yaml:"id"`
	Text string   `yaml:"text"`
	Tags []string `yaml:"tags"`
}

// LeadershipEntry represents a leadership activity
type LeadershipEntry struct {
	ID   string   `yaml:"id"`
	Text string   `yaml:"text"`
	Tags []string `yaml:"tags"`
}

// GetAllIDs returns all IDs from the resume for listing
func (r *Resume) GetAllIDs() []ItemWithID {
	var items []ItemWithID

	// Experience bullets
	for _, exp := range r.Experience {
		for _, bullet := range exp.Bullets {
			items = append(items, ItemWithID{
				ID:       bullet.ID,
				Text:     bullet.Text,
				Tags:     bullet.Tags,
				Section:  "Experience",
				Category: exp.Company,
			})
		}
	}

	// Project bullets
	for _, proj := range r.Projects {
		for _, bullet := range proj.Bullets {
			items = append(items, ItemWithID{
				ID:       bullet.ID,
				Text:     bullet.Text,
				Tags:     bullet.Tags,
				Section:  "Projects",
				Category: proj.Title,
			})
		}
	}

	// Leadership entries
	for _, lead := range r.Leadership {
		items = append(items, ItemWithID{
			ID:       lead.ID,
			Text:     lead.Text,
			Tags:     lead.Tags,
			Section:  "Leadership",
			Category: "",
		})
	}

	return items
}

// ItemWithID represents an item that can be selected by ID
type ItemWithID struct {
	ID       string
	Text     string
	Tags     []string
	Section  string
	Category string
}

// FilterByIDs returns only items matching the given IDs
func (r *Resume) FilterByIDs(ids []string) map[string]bool {
	idSet := make(map[string]bool)
	for _, id := range ids {
		idSet[id] = true
	}
	return idSet
}

// FilterByTags returns IDs matching any of the given tags
func (r *Resume) FilterByTags(tags []string) map[string]bool {
	tagSet := make(map[string]bool)
	for _, tag := range tags {
		tagSet[tag] = true
	}

	selectedIDs := make(map[string]bool)

	for _, exp := range r.Experience {
		for _, bullet := range exp.Bullets {
			for _, tag := range bullet.Tags {
				if tagSet[tag] {
					selectedIDs[bullet.ID] = true
					break
				}
			}
		}
	}

	for _, proj := range r.Projects {
		for _, bullet := range proj.Bullets {
			for _, tag := range bullet.Tags {
				if tagSet[tag] {
					selectedIDs[bullet.ID] = true
					break
				}
			}
		}
	}

	for _, lead := range r.Leadership {
		for _, tag := range lead.Tags {
			if tagSet[tag] {
				selectedIDs[lead.ID] = true
				break
			}
		}
	}

	return selectedIDs
}
