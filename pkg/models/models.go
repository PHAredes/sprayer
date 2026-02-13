package models

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Job represents a job posting
type Job struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Company     string    `json:"company"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Source      string    `json:"source"`
	PostedDate  time.Time `json:"posted_date"`
	Salary      string    `json:"salary,omitempty"`
	JobType     string    `json:"job_type,omitempty"`
	Score       int       `json:"score"`
	// Tracking fields
	EmailSent    bool      `json:"email_sent"`
	EmailOpened  bool      `json:"email_opened"`
	LinkClicked  bool      `json:"link_clicked"`
	Applied      bool      `json:"applied"`
	AppliedDate  time.Time `json:"applied_date,omitempty"`
}

// Profile represents user configuration
type Profile struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Keywords    []string `json:"keywords"`
	Locations   []string `json:"locations"`
	Companies   []string `json:"companies"`
	MinScore    int      `json:"min_score"`
	CreatedAt   time.Time `json:"created_at"`
	LastUsed    time.Time `json:"last_used"`
}

// CVData represents user CV information
type CVData struct {
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	Phone          string   `json:"phone"`
	Location       string   `json:"location"`
	Summary        string   `json:"summary"`
	Experience     []string `json:"experience"`
	Skills         []string `json:"skills"`
	Education      []string `json:"education"`
	Projects       []string `json:"projects"`
	Certifications []string `json:"certifications"`
}

// Config represents application configuration
type Config struct {
	DataDirectory string   `json:"data_directory"`
	Database      string   `json:"database"`
	DefaultProfile string  `json:"default_profile"`
	ExportFormat  string   `json:"export_format"`
	Scrapers      []string `json:"scrapers"`
}

// Database setup
func InitDB() (*sql.DB, error) {
	dataDir := filepath.Join(os.Getenv("HOME"), ".jobscraper")
	os.MkdirAll(dataDir, 0755)
	
	dbPath := filepath.Join(dataDir, "jobscraper.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create tables
	tables := []string{
		`CREATE TABLE IF NOT EXISTS jobs (
			id TEXT PRIMARY KEY,
			title TEXT,
			company TEXT,
			location TEXT,
			description TEXT,
			url TEXT,
			source TEXT,
			posted_date DATETIME,
			salary TEXT,
			job_type TEXT,
			score INTEGER,
			email_sent BOOLEAN DEFAULT 0,
			email_opened BOOLEAN DEFAULT 0,
			link_clicked BOOLEAN DEFAULT 0,
			applied BOOLEAN DEFAULT 0,
			applied_date DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS profiles (
			id TEXT PRIMARY KEY,
			name TEXT,
			keywords TEXT,
			locations TEXT,
			companies TEXT,
			min_score INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_used DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS cv_data (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT,
			phone TEXT,
			location TEXT,
			summary TEXT,
			experience TEXT,
			skills TEXT,
			education TEXT,
			projects TEXT,
			certifications TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS tracking_events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			job_id TEXT,
			recipient TEXT,
			email_hash TEXT,
			ip_address TEXT,
			user_agent TEXT,
			event_type TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (job_id) REFERENCES jobs(id)
		)`,
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

// Load default configuration
func LoadDefaultConfig() *Config {
	dataDir := filepath.Join(os.Getenv("HOME"), ".jobscraper")
	return &Config{
		DataDirectory: dataDir,
		Database:      "jobscraper.db",
		DefaultProfile: "default",
		ExportFormat:  "json",
		Scrapers:      []string{"mock_scraper"},
	}
}

// Save configuration to file
func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Load configuration from file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}