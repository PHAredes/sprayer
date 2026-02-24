package profile

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

// ImportFormat represents supported import formats
type ImportFormat string

const (
	FormatJSON ImportFormat = "json"
	FormatYAML ImportFormat = "yaml"
	FormatCSV  ImportFormat = "csv"
)

// ProfileImporter handles profile import from various formats
type ProfileImporter struct {
	parser *CVParser
}

// NewProfileImporter creates a new profile importer
func NewProfileImporter() *ProfileImporter {
	return &ProfileImporter{
		parser: NewCVParser(),
	}
}

// ImportProfile imports a profile from a file
func (pi *ProfileImporter) ImportProfile(filepath string, format ImportFormat) (Profile, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Profile{}, fmt.Errorf("failed to read file: %w", err)
	}

	switch format {
	case FormatJSON:
		return pi.importFromJSON(content)
	case FormatYAML:
		return pi.importFromYAML(content)
	case FormatCSV:
		return pi.importFromCSV(content)
	default:
		// Auto-detect format from file extension
		importPath := filepath
		ext := strings.ToLower(strings.TrimPrefix(getFileExt(importPath), "."))
		switch ext {
		case "json":
			return pi.importFromJSON(content)
		case "yaml", "yml":
			return pi.importFromYAML(content)
		case "csv":
			return pi.importFromCSV(content)
		default:
			return pi.importFromJSON(content) // Default to JSON
		}
	}
}

// ImportProfileFromCV imports a profile from a CV file
func (pi *ProfileImporter) ImportProfileFromCV(cvPath string, name string) (Profile, error) {
	cvData, err := pi.parser.ParseCVFromFile(cvPath)
	if err != nil {
		return Profile{}, fmt.Errorf("failed to parse CV: %w", err)
	}

	profile := GenerateProfileFromCV(cvData, name)
	profile.CVData = cvData
	profile.CVMinScore = 20 // Default minimum CV score

	return profile, nil
}

// ExportProfile exports a profile to various formats
func (pi *ProfileImporter) ExportProfile(profile Profile, filepath string, format ImportFormat) error {
	var content []byte
	var err error

	switch format {
	case FormatJSON:
		content, err = json.MarshalIndent(profile, "", "  ")
	case FormatYAML:
		content, err = yaml.Marshal(profile)
	default:
		content, err = json.MarshalIndent(profile, "", "  ")
	}

	if err != nil {
		return fmt.Errorf("failed to marshal profile: %w", err)
	}

	return ioutil.WriteFile(filepath, content, 0644)
}

// ValidateProfile validates imported profile data
func (pi *ProfileImporter) ValidateProfile(profile Profile) error {
	if profile.Name == "" {
		return fmt.Errorf("profile name is required")
	}

	if len(profile.Keywords) == 0 && len(profile.PreferredTech) == 0 {
		return fmt.Errorf("profile must have keywords or preferred technologies")
	}

	if profile.MinScore < 0 || profile.MinScore > 100 {
		return fmt.Errorf("min score must be between 0 and 100")
	}

	if profile.MaxScore < profile.MinScore || profile.MaxScore > 100 {
		return fmt.Errorf("max score must be between min score and 100")
	}

	if profile.CVMinScore < 0 {
		return fmt.Errorf("CV min score cannot be negative")
	}

	return nil
}

// importFromJSON imports from JSON format
func (pi *ProfileImporter) importFromJSON(content []byte) (Profile, error) {
	var profile Profile
	err := json.Unmarshal(content, &profile)
	if err != nil {
		return Profile{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if err := pi.ValidateProfile(profile); err != nil {
		return Profile{}, fmt.Errorf("invalid profile: %w", err)
	}

	return profile, nil
}

// importFromYAML imports from YAML format
func (pi *ProfileImporter) importFromYAML(content []byte) (Profile, error) {
	var profile Profile
	err := yaml.Unmarshal(content, &profile)
	if err != nil {
		return Profile{}, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if err := pi.ValidateProfile(profile); err != nil {
		return Profile{}, fmt.Errorf("invalid profile: %w", err)
	}

	return profile, nil
}

// importFromCSV imports from CSV format (simplified)
func (pi *ProfileImporter) importFromCSV(content []byte) (Profile, error) {
	lines := strings.Split(string(content), "\n")
	if len(lines) < 2 {
		return Profile{}, fmt.Errorf("CSV must have header and at least one data row")
	}

	// Simple CSV parsing - in production, use proper CSV library
	header := strings.Split(lines[0], ",")
	data := strings.Split(lines[1], ",")

	if len(header) != len(data) {
		return Profile{}, fmt.Errorf("CSV header and data row length mismatch")
	}

	profile := NewDefaultProfile()

	// Map CSV columns to profile fields (simplified)
	for i, col := range header {
		if i >= len(data) {
			break
		}
		value := strings.TrimSpace(data[i])

		switch strings.ToLower(col) {
		case "name":
			profile.Name = value
		case "keywords":
			profile.Keywords = strings.Split(value, ";")
		case "locations":
			profile.Locations = strings.Split(value, ";")
		case "prefer_remote":
			profile.PreferRemote = value == "true"
		case "min_score":
			fmt.Sscanf(value, "%d", &profile.MinScore)
		case "max_score":
			fmt.Sscanf(value, "%d", &profile.MaxScore)
		}
	}

	if err := pi.ValidateProfile(profile); err != nil {
		return Profile{}, fmt.Errorf("invalid profile: %w", err)
	}

	return profile, nil
}

// SampleProfileJSON returns a sample profile in JSON format for reference
func SampleProfileJSON() string {
	sample := NewDefaultProfile()
	sample.Name = "Software Engineer"
	sample.Keywords = []string{"golang", "rust", "react", "kubernetes"}
	sample.Locations = []string{"remote", "san francisco", "new york"}
	sample.PreferRemote = true
	sample.MinScore = 60
	sample.MaxScore = 95

	jsonBytes, _ := json.MarshalIndent(sample, "", "  ")
	return string(jsonBytes)
}

// SampleProfileYAML returns a sample profile in YAML format for reference
func SampleProfileYAML() string {
	sample := NewDefaultProfile()
	sample.Name = "Data Scientist"
	sample.Keywords = []string{"python", "machine learning", "tensorflow", "sql"}
	sample.Locations = []string{"remote", "boston", "seattle"}
	sample.PreferRemote = true
	sample.MinScore = 70
	sample.MaxScore = 100

	yamlBytes, _ := yaml.Marshal(sample)
	return string(yamlBytes)
}

// getFileExt extracts file extension from path
func getFileExt(path string) string {
	parts := strings.Split(path, ".")
	if len(parts) > 1 {
		return "." + parts[len(parts)-1]
	}
	return ""
}
