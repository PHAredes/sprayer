package profile

import (
	"database/sql"
	"encoding/json"
	"strings"
)

// Store handles profile persistence.
type Store struct {
	db *sql.DB
}

// NewStore wraps a database connection for profile storage.
func NewStore(db *sql.DB) (*Store, error) {
	if err := migrate(db); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS profiles (
			id            TEXT PRIMARY KEY,
			name          TEXT,
			keywords      TEXT,
			cv_path       TEXT,
			cover_path    TEXT,
			contact_email TEXT,
			prefer_remote BOOLEAN DEFAULT 0,
			locations     TEXT
		)`)
	return err
}

// Save upserts a profile.
func (s *Store) Save(p Profile) error {
	kw, _ := json.Marshal(p.Keywords)
	locs, _ := json.Marshal(p.Locations)
	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO profiles
		(id, name, keywords, cv_path, cover_path, contact_email, prefer_remote, locations)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.Name, string(kw), p.CVPath, p.CoverPath,
		p.ContactEmail, p.PreferRemote, string(locs))
	return err
}

// All returns all profiles.
func (s *Store) All() ([]Profile, error) {
	rows, err := s.db.Query(`
		SELECT id, name, keywords, cv_path, cover_path, contact_email, prefer_remote, locations
		FROM profiles ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []Profile
	for rows.Next() {
		var p Profile
		var kwJSON, locsJSON string
		err := rows.Scan(&p.ID, &p.Name, &kwJSON, &p.CVPath, &p.CoverPath,
			&p.ContactEmail, &p.PreferRemote, &locsJSON)
		if err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(kwJSON), &p.Keywords)
		json.Unmarshal([]byte(locsJSON), &p.Locations)
		profiles = append(profiles, p)
	}
	return profiles, nil
}

// ByID returns a single profile.
func (s *Store) ByID(id string) (*Profile, error) {
	row := s.db.QueryRow(`
		SELECT id, name, keywords, cv_path, cover_path, contact_email, prefer_remote, locations
		FROM profiles WHERE id = ?`, strings.ToLower(id))

	var p Profile
	var kwJSON, locsJSON string
	err := row.Scan(&p.ID, &p.Name, &kwJSON, &p.CVPath, &p.CoverPath,
		&p.ContactEmail, &p.PreferRemote, &locsJSON)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(kwJSON), &p.Keywords)
	json.Unmarshal([]byte(locsJSON), &p.Locations)
	return &p, nil
}

// Delete removes a profile.
func (s *Store) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM profiles WHERE id = ?", id)
	return err
}
