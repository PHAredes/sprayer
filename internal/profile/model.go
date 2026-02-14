package profile

// Profile represents a person-specific application profile.
// Links to a CV variant, cover letter template, and keywords for matching.
type Profile struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Keywords     []string `json:"keywords"`
	CVPath       string   `json:"cv_path"`
	CoverPath    string   `json:"cover_path"`
	ContactEmail string   `json:"contact_email"`
	PreferRemote bool     `json:"prefer_remote"`
	Locations    []string `json:"locations"`
}
