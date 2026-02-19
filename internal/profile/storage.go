package profile

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
			id                TEXT PRIMARY KEY,
			name              TEXT,
			keywords          TEXT,
			cv_path           TEXT,
			cover_path        TEXT,
			contact_email     TEXT,
			prefer_remote     BOOLEAN DEFAULT 0,
			locations         TEXT,
			min_score         INTEGER DEFAULT 0,
			max_score         INTEGER DEFAULT 100,
			exclude_traps     BOOLEAN DEFAULT 1,
			must_have_email   BOOLEAN DEFAULT 0,
			preferred_tech    TEXT,
			avoid_tech        TEXT,
			preferred_companies TEXT,
			avoid_companies   TEXT
		)`)
	return err
}

// Save upserts a profile.
func (s *Store) Save(p Profile) error {
	kw, _ := json.Marshal(p.Keywords)
	locs, _ := json.Marshal(p.Locations)
	prefTech, _ := json.Marshal(p.PreferredTech)
	avoidTech, _ := json.Marshal(p.AvoidTech)
	prefCompanies, _ := json.Marshal(p.PreferredCompanies)
	avoidCompanies, _ := json.Marshal(p.AvoidCompanies)
	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO profiles
		(id, name, keywords, cv_path, cover_path, contact_email, prefer_remote, locations,
		 min_score, max_score, exclude_traps, must_have_email,
		 preferred_tech, avoid_tech, preferred_companies, avoid_companies)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.Name, string(kw), p.CVPath, p.CoverPath,
		p.ContactEmail, p.PreferRemote, string(locs),
		p.MinScore, p.MaxScore, p.ExcludeTraps, p.MustHaveEmail,
		string(prefTech), string(avoidTech), string(prefCompanies), string(avoidCompanies))
	return err
}

// All returns all profiles.
func (s *Store) All() ([]Profile, error) {
	rows, err := s.db.Query(`
		SELECT id, name, keywords, cv_path, cover_path, contact_email, prefer_remote, locations,
		       min_score, max_score, exclude_traps, must_have_email,
		       preferred_tech, avoid_tech, preferred_companies, avoid_companies
		FROM profiles ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []Profile
	for rows.Next() {
		var p Profile
		var kwJSON, locsJSON string
		var prefTechJSON, avoidTechJSON string
		var prefCompaniesJSON, avoidCompaniesJSON string
		err := rows.Scan(&p.ID, &p.Name, &kwJSON, &p.CVPath, &p.CoverPath,
			&p.ContactEmail, &p.PreferRemote, &locsJSON,
			&p.MinScore, &p.MaxScore, &p.ExcludeTraps, &p.MustHaveEmail,
			&prefTechJSON, &avoidTechJSON, &prefCompaniesJSON, &avoidCompaniesJSON)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(kwJSON), &p.Keywords); err != nil {
			return nil, fmt.Errorf("unmarshal keywords: %w", err)
		}
		if err := json.Unmarshal([]byte(locsJSON), &p.Locations); err != nil {
			return nil, fmt.Errorf("unmarshal locations: %w", err)
		}
		if err := json.Unmarshal([]byte(prefTechJSON), &p.PreferredTech); err != nil {
			p.PreferredTech = nil
		}
		if err := json.Unmarshal([]byte(avoidTechJSON), &p.AvoidTech); err != nil {
			p.AvoidTech = nil
		}
		if err := json.Unmarshal([]byte(prefCompaniesJSON), &p.PreferredCompanies); err != nil {
			p.PreferredCompanies = nil
		}
		if err := json.Unmarshal([]byte(avoidCompaniesJSON), &p.AvoidCompanies); err != nil {
			p.AvoidCompanies = nil
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

// ByID returns a single profile.
func (s *Store) ByID(id string) (*Profile, error) {
	row := s.db.QueryRow(`
		SELECT id, name, keywords, cv_path, cover_path, contact_email, prefer_remote, locations,
		       min_score, max_score, exclude_traps, must_have_email,
		       preferred_tech, avoid_tech, preferred_companies, avoid_companies
		FROM profiles WHERE id = ?`, strings.ToLower(id))

	var p Profile
	var kwJSON, locsJSON string
	var prefTechJSON, avoidTechJSON string
	var prefCompaniesJSON, avoidCompaniesJSON string
	err := row.Scan(&p.ID, &p.Name, &kwJSON, &p.CVPath, &p.CoverPath,
		&p.ContactEmail, &p.PreferRemote, &locsJSON,
		&p.MinScore, &p.MaxScore, &p.ExcludeTraps, &p.MustHaveEmail,
		&prefTechJSON, &avoidTechJSON, &prefCompaniesJSON, &avoidCompaniesJSON)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(kwJSON), &p.Keywords); err != nil {
		return nil, fmt.Errorf("unmarshal keywords: %w", err)
	}
	if err := json.Unmarshal([]byte(locsJSON), &p.Locations); err != nil {
		return nil, fmt.Errorf("unmarshal locations: %w", err)
	}
	if err := json.Unmarshal([]byte(prefTechJSON), &p.PreferredTech); err != nil {
		p.PreferredTech = nil
	}
	if err := json.Unmarshal([]byte(avoidTechJSON), &p.AvoidTech); err != nil {
		p.AvoidTech = nil
	}
	if err := json.Unmarshal([]byte(prefCompaniesJSON), &p.PreferredCompanies); err != nil {
		p.PreferredCompanies = nil
	}
	if err := json.Unmarshal([]byte(avoidCompaniesJSON), &p.AvoidCompanies); err != nil {
		p.AvoidCompanies = nil
	}
	return &p, nil
}

// Delete removes a profile.
func (s *Store) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM profiles WHERE id = ?", id)
	return err
}
