package job

import (
	"database/sql"
	"os"
	"path/filepath"

	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Store handles job persistence.
type Store struct {
	DB *sql.DB
}

// NewStore opens (or creates) the SQLite database.
func NewStore() (*Store, error) {
	dir := filepath.Join(os.Getenv("HOME"), ".sprayer")
	os.MkdirAll(dir, 0755)

	db, err := sql.Open("sqlite3", filepath.Join(dir, "sprayer.db"))
	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return &Store{DB: db}, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS jobs (
			id          TEXT PRIMARY KEY,
			title       TEXT,
			company     TEXT,
			location    TEXT,
			description TEXT,
			url         TEXT,
			source      TEXT,
			posted_date DATETIME,
			salary      TEXT,
			job_type    TEXT,
			email       TEXT,
			score       INTEGER,
			applied     BOOLEAN DEFAULT 0,
			applied_date DATETIME,
			created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS history (
			key TEXT PRIMARY KEY,
			last_run DATETIME
		)`)
	return err
}

// Save upserts jobs into the database.
func (s *Store) Save(jobs []Job) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT OR REPLACE INTO jobs
		(id, title, company, location, description, url, source, posted_date, salary, job_type, email, score, applied, applied_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, j := range jobs {
		_, err := stmt.Exec(j.ID, j.Title, j.Company, j.Location, j.Description,
			j.URL, j.Source, j.PostedDate, j.Salary, j.JobType, j.Email,
			j.Score, j.Applied, j.AppliedDate)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// All returns every job in the database.
func (s *Store) All() ([]Job, error) {
	rows, err := s.DB.Query(`
		SELECT id, title, company, location, description, url, source,
		       posted_date, salary, job_type, email, score, applied, applied_date
		FROM jobs ORDER BY score DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanJobs(rows)
}

// ByID returns a single job.
func (s *Store) ByID(id string) (*Job, error) {
	row := s.DB.QueryRow(`
		SELECT id, title, company, location, description, url, source,
		       posted_date, salary, job_type, email, score, applied, applied_date
		FROM jobs WHERE id = ?`, id)

	var j Job
	err := row.Scan(&j.ID, &j.Title, &j.Company, &j.Location, &j.Description,
		&j.URL, &j.Source, &j.PostedDate, &j.Salary, &j.JobType, &j.Email,
		&j.Score, &j.Applied, &j.AppliedDate)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func scanJobs(rows *sql.Rows) ([]Job, error) {
	var jobs []Job
	for rows.Next() {
		var j Job
		err := rows.Scan(&j.ID, &j.Title, &j.Company, &j.Location, &j.Description,
			&j.URL, &j.Source, &j.PostedDate, &j.Salary, &j.JobType, &j.Email,
			&j.Score, &j.Applied, &j.AppliedDate)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

// GetLastScrape returns the last time a scrape was run for the given key.
func (s *Store) GetLastScrape(key string) (time.Time, error) {
	var t time.Time
	err := s.DB.QueryRow("SELECT last_run FROM history WHERE key = ?", key).Scan(&t)
	if err == sql.ErrNoRows {
		return time.Time{}, nil
	}
	return t, err
}

// SetLastScrape updates the last scrape time for the given key.
func (s *Store) SetLastScrape(key string) error {
	_, err := s.DB.Exec("INSERT OR REPLACE INTO history (key, last_run) VALUES (?, ?)", key, time.Now())
	return err
}

// Close closes the database.
func (s *Store) Close() error {
	return s.DB.Close()
}
