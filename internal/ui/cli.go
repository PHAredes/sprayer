package ui

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"sprayer/internal/apply"
	"sprayer/internal/job"
	"sprayer/internal/llm"
	"sprayer/internal/profile"
	"sprayer/internal/scraper"
)

// CLI implements the command-line interface logic.
type CLI struct {
	store        *job.Store
	profileStore *profile.Store
	llmClient    *llm.Client
}

func NewCLI() (*CLI, error) {
	s, err := job.NewStore()
	if err != nil {
		return nil, err
	}
	pStore, err := profile.NewStore(s.DB)
	if err != nil {
		return nil, err
	}
	return &CLI{
		store:        s,
		profileStore: pStore,
		llmClient:    llm.NewClient(),
	}, nil
}

func (c *CLI) Run() {
	if len(os.Args) < 2 {
		c.printUsage()
		return
	}

	switch os.Args[1] {
	case "scrape":
		c.handleScrape()
	case "list":
		c.handleList()
	case "apply":
		c.handleApply()
	case "profile":
		c.handleProfile()
	case "setup":
		c.handleSetup()
	case "tui":
		c.handleTUI()
	default:
		c.printUsage()
	}
}

func (c *CLI) printUsage() {
	fmt.Println(`Sprayer - The Agentic Job Application Tool

Usage:
  sprayer <command> [flags]

Commands:
  scrape   Fetch jobs from all sources
  list     List and filter jobs (pipeable)
  apply    Apply to a specific job (generates draft)
  list     List and filter jobs (pipeable)
  apply    Apply to a specific job (generates draft)
   profile  Manage profiles
   setup    Configure SMTP and LLM settings
   tui      Launch interactive terminal UI`)
}

func (c *CLI) handleScrape() {
	fs := flag.NewFlagSet("scrape", flag.ExitOnError)
	fast := fs.Bool("fast", false, "Use only reliable API scrapers (recommended)")
	all := fs.Bool("all", false, "Use all API scrapers (may have issues)")
	force := fs.Bool("force", false, "Force scrape even if recently run")

	// Parse flags first
	if len(os.Args) > 2 {
		fs.Parse(os.Args[2:])
	}

	keywords := fs.Args()
	if len(keywords) == 0 {
		keywords = []string{"golang", "rust", "remote"}
	}

	fmt.Printf("Scraping for: %v (fast=%v)\n", keywords, *fast)

	// Check history
	cacheKey := fmt.Sprintf("%v-fast=%v", keywords, *fast)
	lastRun, _ := c.store.GetLastScrape(cacheKey)
	if !*force && time.Since(lastRun) < 15*time.Minute {
		fmt.Printf("Skipping scrape (run %v ago). Use --force to override.\n", time.Since(lastRun).Round(time.Second))
		return
	}

	var s job.Scraper
	if *all {
		s = scraper.APIOnly() // All API scrapers
	} else if *fast {
		s = scraper.APIOnlyReliable() // Only reliable API scrapers
	} else {
		s = scraper.All(keywords, "Remote") // Full scraping with browser
	}

	jobs, err := s()
	if err != nil {
		fmt.Printf("Scrape error: %v\n", err)
		return
	}

	// Flag and sanitize before saving
	pipeline := job.Pipe(job.FlagTraps(), job.SanitizeDescriptions())
	processed := pipeline(jobs)

	c.store.Save(processed)
	c.store.SetLastScrape(cacheKey)
	fmt.Printf("Saved %d jobs.\n", len(processed))
}

func (c *CLI) handleList() {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	keywords := fs.String("keywords", "", "Filter by keywords (comma-sep)")
	minScore := fs.Int("min-score", 0, "Filter by minimum score")
	fs.Parse(os.Args[2:])

	jobs, _ := c.store.All()

	filters := []job.Filter{
		job.Dedup(),
		job.FlagTraps(),
		job.SanitizeDescriptions(),
	}
	if *keywords != "" {
		filters = append(filters, job.ByKeywords(strings.Split(*keywords, ",")))
	}
	if *minScore > 0 {
		filters = append(filters, job.ByMinScore(*minScore))
	}

	pipeline := job.Pipe(filters...)
	filtered := pipeline(jobs)

	for _, j := range filtered {
		trapIndicator := ""
		if j.HasTraps {
			trapIndicator = " [!] TRAPS FOUND"
		}
		fmt.Printf("[%d]%s %s @ %s (%s)\n", j.Score, trapIndicator, j.Title, j.Company, j.ID)
	}
}

func (c *CLI) handleApply() {
	fs := flag.NewFlagSet("apply", flag.ExitOnError)
	jobID := fs.String("job", "", "Job ID to apply to")
	prompt := fs.String("prompt", "email_cold", "Message prompt template")
	send := fs.Bool("send", false, "Send email immediately via SMTP")
	fs.Parse(os.Args[2:])

	if *jobID == "" {
		fmt.Println("Error: --job is required")
		return
	}

	j, err := c.store.ByID(*jobID)
	if err != nil {
		fmt.Printf("Job not found: %v\n", err)
		return
	}

	profiles, _ := c.profileStore.All()
	// Use first profile for now - can be enhanced later
	var p profile.Profile
	if len(profiles) > 0 {
		p = profiles[0]
	} else {
		p = profile.NewDefaultProfile()
	}

	fmt.Printf("Generating application for %s using profile %s...\n", j.Company, p.Name)

	subject, body, err := apply.GenerateEmail(*j, p, c.llmClient, *prompt)
	if err != nil {
		fmt.Printf("Generation failed: %v\n", err)
		return
	}

	path, err := apply.Draft(*j, p, subject, body)
	if err != nil {
		fmt.Printf("Draft failed: %v\n", err)
		return
	}

	fmt.Printf("Draft created: %s\n", path)

	if *send {
		fmt.Printf("Sending email via SMTP...\n")
		// Assume CV is attached if path exists and ends with .pdf, but Draft only returns .eml path?
		// apply.Draft saves .eml. Attachment is usually handled inside Draft or external.
		// Wait, Draft function saves the .eml file.
		// SendDirect needs the attachment path (PDF) separately.
		// Let's assume Profile has CV path.
		cvPath := p.CVPath
		err := apply.SendDirect(j.Email, subject, body, cvPath)
		if err != nil {
			fmt.Printf("Failed to send: %v\n", err)
		} else {
			fmt.Printf("Email sent successfully to %s!\n", j.Email)
		}
	}
}

func (c *CLI) handleProfile() {
	// Stub for now
	profiles, _ := c.profileStore.All()
	for _, p := range profiles {
		fmt.Printf("- %s (%s)\n", p.Name, p.ID)
	}
}

func (c *CLI) handleTUI() {
	fmt.Println("Starting enhanced TUI...")

	// Initialize the new TUI
	if err := InitializeTUI(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		return
	}
}
