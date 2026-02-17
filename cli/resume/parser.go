package resume

import (
	"os"

	"gopkg.in/yaml.v3"
)

// LoadResume reads and parses the resume YAML file
func LoadResume(path string) (*Resume, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var resume Resume
	if err := yaml.Unmarshal(data, &resume); err != nil {
		return nil, err
	}

	return &resume, nil
}
