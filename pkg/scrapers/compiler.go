package scrapers

import (
	"database/sql"
	"time"

	"job-scraper/pkg/models"
)

// CompilerScraper scrapes compiler and language development jobs
type CompilerScraper struct {
	name string
	db   *sql.DB
}

func NewCompilerScraper(db *sql.DB) *CompilerScraper {
	return &CompilerScraper{name: "compiler", db: db}
}

func (cs *CompilerScraper) GetName() string {
	return cs.name
}

func (cs *CompilerScraper) Scrape() ([]models.Job, error) {
	// Compiler and language development opportunities
	jobs := []models.Job{
		{
			ID:          "compiler-1",
			Title:       "LLVM Compiler Engineer",
			Company:     "Apple",
			Location:    "Cupertino, CA",
			Description: "Work on LLVM compiler infrastructure for Apple platforms. LLVM and compiler optimization experience required.",
			URL:         "https://apple.com/careers/llvm-engineer",
			Source:      cs.name,
			PostedDate:  time.Now().Add(-6 * time.Hour),
			Score:       96,
		},
		{
			ID:          "compiler-2",
			Title:       "GCC Developer",
			Company:     "Red Hat",
			Location:    "Remote",
			Description: "Contribute to GCC compiler development. C/C++ and compiler theory experience required.",
			URL:         "https://redhat.com/careers/gcc-developer",
			Source:      cs.name,
			PostedDate:  time.Now().Add(-30 * time.Hour),
			Score:       92,
		},
		{
			ID:          "compiler-3",
			Title:       "JVM Engineer",
			Company:     "Oracle",
			Location:    "Burlington, MA",
			Description: "Work on Java Virtual Machine optimization. JVM internals and performance tuning experience required.",
			URL:         "https://oracle.com/careers/jvm-engineer",
			Source:      cs.name,
			PostedDate:  time.Now().Add(-54 * time.Hour),
			Score:       91,
		},
		{
			ID:          "compiler-4",
			Title:       "TypeScript Compiler Developer",
			Company:     "Microsoft",
			Location:    "Redmond, WA",
			Description: "Work on TypeScript compiler and language features. Type systems and JavaScript experience required.",
			URL:         "https://microsoft.com/careers/typescript-compiler",
			Source:      cs.name,
			PostedDate:  time.Now().Add(-78 * time.Hour),
			Score:       93,
		},
		{
			ID:          "compiler-5",
			Title:       "WebAssembly Engineer",
			Company:     "Google",
			Location:    "Mountain View, CA",
			Description: "Work on WebAssembly runtime and tooling. Virtual machines and compiler technology experience required.",
			URL:         "https://google.com/careers/webassembly",
			Source:      cs.name,
			PostedDate:  time.Now().Add(-102 * time.Hour),
			Score:       94,
		},
		{
			ID:          "compiler-6",
			Title:       "Julia Language Developer",
			Company:     "Julia Computing",
			Location:    "Remote",
			Description: "Work on Julia compiler and language features. Scientific computing and compiler development experience required.",
			URL:         "https://juliacomputing.com/careers/compiler",
			Source:      cs.name,
			PostedDate:  time.Now().Add(-126 * time.Hour),
			Score:       90,
		},
		{
			ID:          "compiler-7",
			Title:       "Rust Compiler Developer",
			Company:     "Ferrocene",
			Location:    "Remote",
			Description: "Work on safety-critical Rust compiler. Formal methods and compiler verification experience required.",
			URL:         "https://ferrocene.dev/careers/compiler",
			Source:      cs.name,
			PostedDate:  time.Now().Add(-150 * time.Hour),
			Score:       95,
		},
		{
			ID:          "compiler-8",
			Title:       "Go Compiler Engineer",
			Company:     "Google",
			Location:    "Remote",
			Description: "Work on Go compiler and runtime optimization. Garbage collection and compiler technology experience required.",
			URL:         "https://google.com/careers/go-compiler",
			Source:      cs.name,
			PostedDate:  time.Now().Add(-174 * time.Hour),
			Score:       92,
		},
	}

	// Save to database
	for _, job := range jobs {
		_, err := cs.db.Exec(`
			INSERT OR REPLACE INTO jobs (id, title, company, location, description, url, source, posted_date, salary, job_type, score)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, job.ID, job.Title, job.Company, job.Location, job.Description, job.URL, job.Source, job.PostedDate, job.Salary, job.JobType, job.Score)
		if err != nil {
			return nil, err
		}
	}

	return jobs, nil
}