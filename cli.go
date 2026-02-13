package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"job-scraper/pkg/models"
	"job-scraper/pkg/scrapers"
)

// CLI command structure
type CLI struct {
	scrapeCmd    *flag.FlagSet
	filterCmd    *flag.FlagSet
	exportCmd    *flag.FlagSet
	profilesCmd  *flag.FlagSet
	cvCmd        *flag.FlagSet
	settingsCmd  *flag.FlagSet
}

func NewCLI() *CLI {
	cli := &CLI{}

	// Scrape command
	cli.scrapeCmd = flag.NewFlagSet("scrape", flag.ExitOnError)
	cli.scrapeCmd.String("profile", "default", "Profile to use for scraping")
	cli.scrapeCmd.Bool("force", false, "Force re-scraping (ignore cache)")

	// Filter command
	cli.filterCmd = flag.NewFlagSet("filter", flag.ExitOnError)
	cli.filterCmd.String("keywords", "", "Comma-separated keywords to filter")
	cli.filterCmd.Int("min-score", 0, "Minimum job score")
	cli.filterCmd.String("location", "", "Location filter")
	cli.filterCmd.String("company", "", "Company filter")
	cli.filterCmd.String("date-range", "all", "Date range (today/week/month/all)")
	cli.filterCmd.String("output", "", "Output file for filtered jobs")

	// Export command
	cli.exportCmd = flag.NewFlagSet("export", flag.ExitOnError)
	cli.exportCmd.String("format", "json", "Export format (json/csv)")
	cli.exportCmd.String("output", "", "Output file")

	// Profiles command
	cli.profilesCmd = flag.NewFlagSet("profiles", flag.ExitOnError)
	cli.profilesCmd.String("create", "", "Create new profile")
	cli.profilesCmd.String("delete", "", "Delete profile")
	cli.profilesCmd.String("switch", "", "Switch to profile")
	cli.profilesCmd.Bool("list", false, "List all profiles")

	// CV command
	cli.cvCmd = flag.NewFlagSet("cv", flag.ExitOnError)
	cli.cvCmd.String("generate", "", "Generate CV for job ID")
	cli.cvCmd.Bool("edit", false, "Open CV editor")
	cli.cvCmd.String("template", "", "CV template to use")

	// Settings command
	cli.settingsCmd = flag.NewFlagSet("settings", flag.ExitOnError)
	cli.settingsCmd.Bool("show", false, "Show current settings")
	cli.settingsCmd.String("import", "", "Import settings from JSON file")
	cli.settingsCmd.String("export", "", "Export settings to JSON file")

	return cli
}

func (c *CLI) Run() {
	if len(os.Args) < 2 {
		c.printUsage()
		return
	}

	switch os.Args[1] {
	case "scrape":
		c.handleScrape()
	case "filter":
		c.handleFilter()
	case "export":
		c.handleExport()
	case "profiles":
		c.handleProfiles()
	case "cv":
		c.handleCV()
	case "settings":
		c.handleSettings()
	case "help", "--help", "-h":
		c.printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", os.Args[1])
		c.printUsage()
	}
}

func (c *CLI) printUsage() {
	fmt.Println(`JobScraper CLI - Command Line Interface

Usage:
  jobscraper <command> [flags]

Commands:
  scrape     - Scrape jobs from configured sources
  filter     - Filter jobs based on criteria
  export     - Export jobs to file
  profiles   - Manage search profiles
  cv         - CV generation and management
  settings   - Manage application settings
  help       - Show this help message

Examples:
  jobscraper scrape --profile default
  jobscraper filter --keywords "rust,compiler" --min-score 80
  jobscraper export --format json --output jobs.json
  jobscraper profiles --list
  jobscraper cv --generate job123
  jobscraper settings --export config.json

Use 'jobscraper <command> --help' for command-specific help.`)
}

func (c *CLI) handleScrape() {
	c.scrapeCmd.Parse(os.Args[2:])
	
	profile := c.scrapeCmd.Lookup("profile").Value.String()
	force := c.scrapeCmd.Lookup("force").Value.(flag.Getter).Get().(bool)

	fmt.Printf("Scraping jobs with profile: %s\n", profile)
	if force {
		fmt.Println("Force mode: ignoring cache")
	}

	// Initialize database and scraper
	db, err := models.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	scraper := scrapers.NewMockScraper("cli_scraper", db)
	jobs, err := scraper.Scrape()
	if err != nil {
		log.Fatalf("Scraping failed: %v", err)
	}

	fmt.Printf("Scraping complete! Found %d jobs\n", len(jobs))

	// Show top jobs
	fmt.Println("\nTop 5 jobs:")
	for i, job := range jobs {
		if i >= 5 {
			break
		}
		fmt.Printf("  %d. %s - %s (Score: %d)\n", i+1, job.Title, job.Company, job.Score)
	}

	if len(jobs) >= 300 {
		fmt.Println("\n🎉 Target of 300+ jobs reached!")
	}
}

func (c *CLI) handleFilter() {
	c.filterCmd.Parse(os.Args[2:])

	keywords := c.filterCmd.Lookup("keywords").Value.String()
	minScore := c.filterCmd.Lookup("min-score").Value.(flag.Getter).Get().(int)
	location := c.filterCmd.Lookup("location").Value.String()
	company := c.filterCmd.Lookup("company").Value.String()
	dateRange := c.filterCmd.Lookup("date-range").Value.String()
	outputFile := c.filterCmd.Lookup("output").Value.String()

	// Load jobs from database
	db, err := models.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, company, location, description, url, source, posted_date, salary, job_type, score FROM jobs")
	if err != nil {
		log.Fatalf("Failed to query jobs: %v", err)
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(&job.ID, &job.Title, &job.Company, &job.Location, &job.Description, &job.URL, &job.Source, &job.PostedDate, &job.Salary, &job.JobType, &job.Score)
		if err != nil {
			log.Fatalf("Failed to scan job: %v", err)
		}
		jobs = append(jobs, job)
	}

	// Apply filters
	filteredJobs := jobs

	if keywords != "" {
		keywordList := strings.Split(keywords, ",")
		filteredJobs = filterByKeywords(filteredJobs, keywordList)
	}

	if minScore > 0 {
		filteredJobs = filterByScore(filteredJobs, minScore)
	}

	if location != "" {
		filteredJobs = filterByLocation(filteredJobs, location)
	}

	if company != "" {
		filteredJobs = filterByCompany(filteredJobs, company)
	}

	if dateRange != "all" {
		filteredJobs = filterByDateRange(filteredJobs, dateRange)
	}

	fmt.Printf("Filtered %d jobs from %d total\n", len(filteredJobs), len(jobs))

	// Show filtered jobs
	fmt.Println("\nFiltered jobs:")
	for i, job := range filteredJobs {
		fmt.Printf("  %d. %s - %s (Score: %d)\n", i+1, job.Title, job.Company, job.Score)
	}

	// Export if requested
	if outputFile != "" {
		data, err := json.MarshalIndent(filteredJobs, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal jobs: %v", err)
		}

		err = os.WriteFile(outputFile, data, 0644)
		if err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}

		fmt.Printf("\nJobs exported to: %s\n", outputFile)
	}
}

func (c *CLI) handleExport() {
	c.exportCmd.Parse(os.Args[2:])

	format := c.exportCmd.Lookup("format").Value.String()
	outputFile := c.exportCmd.Lookup("output").Value.String()

	if outputFile == "" {
		outputFile = fmt.Sprintf("jobs_export_%s.%s", time.Now().Format("20060102_150405"), format)
	}

	// Load all jobs from database
	db, err := models.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, company, location, description, url, source, posted_date, salary, job_type, score FROM jobs")
	if err != nil {
		log.Fatalf("Failed to query jobs: %v", err)
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(&job.ID, &job.Title, &job.Company, &job.Location, &job.Description, &job.URL, &job.Source, &job.PostedDate, &job.Salary, &job.JobType, &job.Score)
		if err != nil {
			log.Fatalf("Failed to scan job: %v", err)
		}
		jobs = append(jobs, job)
	}

	switch format {
	case "json":
		data, err := json.MarshalIndent(jobs, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal jobs: %v", err)
		}

		err = os.WriteFile(outputFile, data, 0644)
		if err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}

	case "csv":
		// Simple CSV export
		var csvData strings.Builder
		csvData.WriteString("Title,Company,Location,Source,Score,Posted Date,URL\n")

		for _, job := range jobs {
			csvData.WriteString(fmt.Sprintf(`"%s","%s","%s","%s",%d,"%s","%s"\n`,
				job.Title, job.Company, job.Location, job.Source, job.Score,
				job.PostedDate.Format("2006-01-02"), job.URL))
		}

		err = os.WriteFile(outputFile, []byte(csvData.String()), 0644)
		if err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}

	default:
		log.Fatalf("Unsupported format: %s", format)
	}

	fmt.Printf("Exported %d jobs to: %s\n", len(jobs), outputFile)
}

func (c *CLI) handleProfiles() {
	c.profilesCmd.Parse(os.Args[2:])

	create := c.profilesCmd.Lookup("create").Value.String()
	delete := c.profilesCmd.Lookup("delete").Value.String()
	switchTo := c.profilesCmd.Lookup("switch").Value.String()
	list := c.profilesCmd.Lookup("list").Value.(flag.Getter).Get().(bool)

	// Load profiles from database
	db, err := models.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	if create != "" {
		// Create new profile
		_, err := db.Exec(`
			INSERT INTO profiles (id, name, keywords, locations, companies, min_score)
			VALUES (?, ?, ?, ?, ?, ?)
		`, strings.ToLower(create), create, "rust,golang,compiler", "remote", "", 80)
		if err != nil {
			log.Fatalf("Failed to create profile: %v", err)
		}
		fmt.Printf("Created profile: %s\n", create)
		return
	}

	if delete != "" {
		_, err := db.Exec("DELETE FROM profiles WHERE id = ?", strings.ToLower(delete))
		if err != nil {
			log.Fatalf("Failed to delete profile: %v", err)
		}
		fmt.Printf("Deleted profile: %s\n", delete)
		return
	}

	if switchTo != "" {
		_, err := db.Exec("UPDATE profiles SET last_used = CURRENT_TIMESTAMP WHERE id = ?", strings.ToLower(switchTo))
		if err != nil {
			log.Fatalf("Failed to switch profile: %v", err)
		}
		fmt.Printf("Switched to profile: %s\n", switchTo)
		return
	}

	if list {
		rows, err := db.Query("SELECT id, name, keywords, min_score FROM profiles ORDER BY last_used DESC")
		if err != nil {
			log.Fatalf("Failed to query profiles: %v", err)
		}
		defer rows.Close()

		fmt.Println("Available profiles:")
		for rows.Next() {
			var id, name, keywords string
			var minScore int
			err := rows.Scan(&id, &name, &keywords, &minScore)
			if err != nil {
				log.Fatalf("Failed to scan profile: %v", err)
			}
			fmt.Printf("  • %s (ID: %s) - Keywords: %s, Min Score: %d\n", name, id, keywords, minScore)
		}
	}
}

func (c *CLI) handleCV() {
	c.cvCmd.Parse(os.Args[2:])

	generate := c.cvCmd.Lookup("generate").Value.String()
	edit := c.cvCmd.Lookup("edit").Value.(flag.Getter).Get().(bool)
	template := c.cvCmd.Lookup("template").Value.String()

	if generate != "" {
		fmt.Printf("Generating CV for job: %s\n", generate)
		if template != "" {
			fmt.Printf("Using template: %s\n", template)
		}
		fmt.Println("CV generation feature coming soon...")
		return
	}

	if edit {
		fmt.Println("Opening CV editor...")
		fmt.Println("CV editor feature coming soon...")
		return
	}

	fmt.Println("CV management system - use --help for available commands")
}

func (c *CLI) handleSettings() {
	c.settingsCmd.Parse(os.Args[2:])

	show := c.settingsCmd.Lookup("show").Value.(flag.Getter).Get().(bool)
	importFile := c.settingsCmd.Lookup("import").Value.String()
	exportFile := c.settingsCmd.Lookup("export").Value.String()

	if show {
		fmt.Println("Current settings:")
		fmt.Println("  • Data directory: ~/.jobscraper")
		fmt.Println("  • Database: jobscraper.db")
		fmt.Println("  • Default profile: default")
		fmt.Println("  • Export format: json")
		return
	}

	if importFile != "" {
		fmt.Printf("Importing settings from: %s\n", importFile)
		fmt.Println("Settings import feature coming soon...")
		return
	}

	if exportFile != "" {
		fmt.Printf("Exporting settings to: %s\n", exportFile)
		
		// Create sample settings
		settings := map[string]interface{}{
			"data_directory": filepath.Join(os.Getenv("HOME"), ".jobscraper"),
			"database":       "jobscraper.db",
			"default_profile": "default",
			"export_format":   "json",
			"scrapers":        []string{"mock_scraper"},
		}

		data, err := json.MarshalIndent(settings, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal settings: %v", err)
		}

		err = os.WriteFile(exportFile, data, 0644)
		if err != nil {
			log.Fatalf("Failed to write settings: %v", err)
		}

		fmt.Printf("Settings exported to: %s\n", exportFile)
		return
	}

	fmt.Println("Settings management - use --help for available commands")
}

// Filter functions for CLI
func filterByKeywords(jobs []models.Job, keywords []string) []models.Job {
	var filtered []models.Job
	for _, job := range jobs {
		for _, keyword := range keywords {
			keyword = strings.TrimSpace(keyword)
			if strings.Contains(strings.ToLower(job.Title), strings.ToLower(keyword)) ||
				strings.Contains(strings.ToLower(job.Description), strings.ToLower(keyword)) {
				filtered = append(filtered, job)
				break
			}
		}
	}
	return filtered
}

func filterByScore(jobs []models.Job, minScore int) []models.Job {
	var filtered []models.Job
	for _, job := range jobs {
		if job.Score >= minScore {
			filtered = append(filtered, job)
		}
	}
	return filtered
}

func filterByLocation(jobs []models.Job, location string) []models.Job {
	var filtered []models.Job
	for _, job := range jobs {
		if strings.Contains(strings.ToLower(job.Location), strings.ToLower(location)) {
			filtered = append(filtered, job)
		}
	}
	return filtered
}

func filterByCompany(jobs []models.Job, company string) []models.Job {
	var filtered []models.Job
	for _, job := range jobs {
		if strings.Contains(strings.ToLower(job.Company), strings.ToLower(company)) {
			filtered = append(filtered, job)
		}
	}
	return filtered
}

func filterByDateRange(jobs []models.Job, dateRange string) []models.Job {
	var filtered []models.Job
	now := time.Now()

	for _, job := range jobs {
		switch strings.ToLower(dateRange) {
		case "today":
			if job.PostedDate.Year() == now.Year() &&
				job.PostedDate.Month() == now.Month() &&
				job.PostedDate.Day() == now.Day() {
				filtered = append(filtered, job)
			}
		case "week":
			if job.PostedDate.After(now.AddDate(0, 0, -7)) {
				filtered = append(filtered, job)
			}
		case "month":
			if job.PostedDate.After(now.AddDate(0, -1, 0)) {
				filtered = append(filtered, job)
			}
		case "all":
			filtered = append(filtered, job)
		}
	}
	return filtered
}

func main() {
	cli := NewCLI()
	cli.Run()
}