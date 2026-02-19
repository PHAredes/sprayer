package apply

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ScratchEmail struct {
	ID           string    `json:"id"`
	EmailAddress string    `json:"email_address"`
	Provider     string    `json:"provider"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	ForwardTo    string    `json:"forward_to,omitempty"`
	JobID        string    `json:"job_id,omitempty"`
	Active       bool      `json:"active"`

	providerData map[string]interface{}
}

type Email struct {
	ID         string    `json:"id"`
	From       string    `json:"from"`
	To         string    `json:"to"`
	Subject    string    `json:"subject"`
	Body       string    `json:"body"`
	Text       string    `json:"text,omitempty"`
	HTML       string    `json:"html,omitempty"`
	ReceivedAt time.Time `json:"received_at"`
	Read       bool      `json:"read"`
}

type Provider interface {
	Name() string
	Create(forwardTo string) (*ScratchEmail, error)
	CheckInbox(id string) ([]Email, error)
	Deactivate(id string) error
}

type ScratchEmailManager struct {
	db       *sql.DB
	provider Provider
}

func NewScratchEmailManager(db *sql.DB) (*ScratchEmailManager, error) {
	if err := migrateScratchEmails(db); err != nil {
		return nil, fmt.Errorf("migrate scratch emails: %w", err)
	}

	provider := detectProvider()

	return &ScratchEmailManager{
		db:       db,
		provider: provider,
	}, nil
}

func migrateScratchEmails(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS scratch_emails (
			id              TEXT PRIMARY KEY,
			email_address   TEXT UNIQUE NOT NULL,
			provider        TEXT NOT NULL,
			created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at      DATETIME,
			forward_to      TEXT,
			job_id          TEXT,
			active          BOOLEAN DEFAULT 1,
			provider_data   TEXT
		)`)
	return err
}

func detectProvider() Provider {
	providerName := os.Getenv("SPRAYER_SCRATCH_EMAIL_PROVIDER")
	switch providerName {
	case "simplelogin":
		if apiKey := os.Getenv("SIMPLELOGIN_API_KEY"); apiKey != "" {
			return NewSimpleLoginProvider(apiKey)
		}
	case "mailtm":
		return NewMailTMProvider()
	}

	if apiKey := os.Getenv("SIMPLELOGIN_API_KEY"); apiKey != "" {
		return NewSimpleLoginProvider(apiKey)
	}

	return NewMailTMProvider()
}

func getDefaultForwardTo() string {
	return os.Getenv("SPRAYER_SCRATCH_EMAIL_FORWARD_TO")
}

func (m *ScratchEmailManager) CreateScratchEmail(forwardTo string) (*ScratchEmail, error) {
	if m.provider == nil {
		return nil, fmt.Errorf("no scratch email provider configured")
	}

	scratch, err := m.provider.Create(forwardTo)
	if err != nil {
		return nil, fmt.Errorf("provider create: %w", err)
	}

	scratch.ForwardTo = forwardTo
	scratch.Active = true

	if err := m.save(scratch); err != nil {
		return nil, fmt.Errorf("save scratch email: %w", err)
	}

	return scratch, nil
}

func (m *ScratchEmailManager) CreateScratchEmailForJob(forwardTo, jobID string) (*ScratchEmail, error) {
	scratch, err := m.CreateScratchEmail(forwardTo)
	if err != nil {
		return nil, err
	}

	scratch.JobID = jobID
	if err := m.save(scratch); err != nil {
		return nil, fmt.Errorf("update job association: %w", err)
	}

	return scratch, nil
}

func (m *ScratchEmailManager) CheckScratchEmail(id string) ([]Email, error) {
	scratch, err := m.ByID(id)
	if err != nil {
		return nil, fmt.Errorf("scratch email not found: %w", err)
	}

	if !scratch.Active {
		return nil, fmt.Errorf("scratch email is deactivated")
	}

	return m.provider.CheckInbox(scratch.ID)
}

func (m *ScratchEmailManager) DeactivateScratchEmail(id string) error {
	scratch, err := m.ByID(id)
	if err != nil {
		return fmt.Errorf("scratch email not found: %w", err)
	}

	if err := m.provider.Deactivate(scratch.ID); err != nil {
		return fmt.Errorf("provider deactivate: %w", err)
	}

	scratch.Active = false
	return m.save(scratch)
}

func (m *ScratchEmailManager) GetScratchEmail(id string) (*ScratchEmail, error) {
	return m.ByID(id)
}

func (m *ScratchEmailManager) ListScratchEmails() ([]ScratchEmail, error) {
	return m.All()
}

func (m *ScratchEmailManager) ByID(id string) (*ScratchEmail, error) {
	row := m.db.QueryRow(`
		SELECT id, email_address, provider, created_at, expires_at, 
		       forward_to, job_id, active, provider_data
		FROM scratch_emails WHERE id = ?`, id)

	return scanScratchEmail(row)
}

func (m *ScratchEmailManager) ByJobID(jobID string) (*ScratchEmail, error) {
	row := m.db.QueryRow(`
		SELECT id, email_address, provider, created_at, expires_at, 
		       forward_to, job_id, active, provider_data
		FROM scratch_emails WHERE job_id = ?`, jobID)

	return scanScratchEmail(row)
}

func (m *ScratchEmailManager) All() ([]ScratchEmail, error) {
	rows, err := m.db.Query(`
		SELECT id, email_address, provider, created_at, expires_at, 
		       forward_to, job_id, active, provider_data
		FROM scratch_emails ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []ScratchEmail
	for rows.Next() {
		var se ScratchEmail
		var providerData []byte
		err := rows.Scan(
			&se.ID, &se.EmailAddress, &se.Provider, &se.CreatedAt, &se.ExpiresAt,
			&se.ForwardTo, &se.JobID, &se.Active, &providerData,
		)
		if err != nil {
			return nil, err
		}
		if len(providerData) > 0 {
			json.Unmarshal(providerData, &se.providerData)
		}
		emails = append(emails, se)
	}
	return emails, nil
}

func (m *ScratchEmailManager) ActiveEmails() ([]ScratchEmail, error) {
	rows, err := m.db.Query(`
		SELECT id, email_address, provider, created_at, expires_at, 
		       forward_to, job_id, active, provider_data
		FROM scratch_emails WHERE active = 1 ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []ScratchEmail
	for rows.Next() {
		se, err := scanScratchEmailFromRow(rows)
		if err != nil {
			return nil, err
		}
		emails = append(emails, *se)
	}
	return emails, nil
}

func (m *ScratchEmailManager) save(scratch *ScratchEmail) error {
	providerData, _ := json.Marshal(scratch.providerData)

	_, err := m.db.Exec(`
		INSERT OR REPLACE INTO scratch_emails 
		(id, email_address, provider, created_at, expires_at, forward_to, 
		 job_id, active, provider_data)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		scratch.ID, scratch.EmailAddress, scratch.Provider, scratch.CreatedAt,
		scratch.ExpiresAt, scratch.ForwardTo, scratch.JobID,
		scratch.Active, providerData)

	return err
}

func scanScratchEmail(row *sql.Row) (*ScratchEmail, error) {
	var se ScratchEmail
	var providerData []byte
	err := row.Scan(
		&se.ID, &se.EmailAddress, &se.Provider, &se.CreatedAt, &se.ExpiresAt,
		&se.ForwardTo, &se.JobID, &se.Active, &providerData,
	)
	if err != nil {
		return nil, err
	}
	if len(providerData) > 0 {
		json.Unmarshal(providerData, &se.providerData)
	}
	return &se, nil
}

func scanScratchEmailFromRow(rows *sql.Rows) (*ScratchEmail, error) {
	var se ScratchEmail
	var providerData []byte
	err := rows.Scan(
		&se.ID, &se.EmailAddress, &se.Provider, &se.CreatedAt, &se.ExpiresAt,
		&se.ForwardTo, &se.JobID, &se.Active, &providerData,
	)
	if err != nil {
		return nil, err
	}
	if len(providerData) > 0 {
		json.Unmarshal(providerData, &se.providerData)
	}
	return &se, nil
}

type MailTMProvider struct {
	client    *http.Client
	baseURL   string
	accountID string
	password  string
	token     string
}

func NewMailTMProvider() *MailTMProvider {
	return &MailTMProvider{
		client:  &http.Client{Timeout: 30 * time.Second},
		baseURL: "https://api.mail.tm",
	}
}

func (p *MailTMProvider) Name() string {
	return "mail.tm"
}

func (p *MailTMProvider) Create(forwardTo string) (*ScratchEmail, error) {
	domainsResp, err := p.getDomains()
	if err != nil {
		return nil, fmt.Errorf("get domains: %w", err)
	}
	if len(domainsResp) == 0 {
		return nil, fmt.Errorf("no domains available")
	}
	domain := domainsResp[0].Domain

	timestamp := time.Now().UnixNano()
	address := fmt.Sprintf("sprayer_%d@%s", timestamp, domain)
	password := generateRandomPassword()

	account, err := p.createAccount(address, password)
	if err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}

	token, err := p.login(address, password)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	p.accountID = account.ID
	p.password = password
	p.token = token

	return &ScratchEmail{
		ID:           account.ID,
		EmailAddress: address,
		Provider:     p.Name(),
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		ForwardTo:    forwardTo,
		Active:       true,
		providerData: map[string]interface{}{
			"password": password,
			"token":    token,
		},
	}, nil
}

func (p *MailTMProvider) CheckInbox(id string) ([]Email, error) {
	if p.token == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	req, _ := http.NewRequest("GET", p.baseURL+"/messages", nil)
	req.Header.Set("Authorization", "Bearer "+p.token)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch messages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var messages []struct {
		ID        string                   `json:"id"`
		From      struct{ Address string } `json:"from"`
		Subject   string                   `json:"subject"`
		CreatedAt time.Time                `json:"createdAt"`
		Seen      bool                     `json:"seen"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, fmt.Errorf("decode messages: %w", err)
	}

	var emails []Email
	for _, msg := range messages {
		email := Email{
			ID:         msg.ID,
			From:       msg.From.Address,
			To:         "",
			Subject:    msg.Subject,
			ReceivedAt: msg.CreatedAt,
			Read:       msg.Seen,
		}

		body, err := p.getMessageBody(msg.ID)
		if err == nil {
			email.Body = body
		}

		emails = append(emails, email)
	}

	return emails, nil
}

func (p *MailTMProvider) getMessageBody(id string) (string, error) {
	req, _ := http.NewRequest("GET", p.baseURL+"/messages/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+p.token)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var msg struct {
		Text string `json:"text"`
		HTML string `json:"html"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return "", err
	}

	if msg.Text != "" {
		return msg.Text, nil
	}
	return msg.HTML, nil
}

func (p *MailTMProvider) Deactivate(id string) error {
	return nil
}

func (p *MailTMProvider) getDomains() ([]struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
}, error) {
	resp, err := p.client.Get(p.baseURL + "/domains")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Members []struct {
			ID     string `json:"id"`
			Domain string `json:"domain"`
		} `json:"hydra:member"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Members, nil
}

func (p *MailTMProvider) createAccount(address, password string) (*struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}, error) {
	payload := map[string]string{
		"address":  address,
		"password": password,
	}
	body, _ := json.Marshal(payload)

	resp, err := p.client.Post(p.baseURL+"/accounts", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create account failed: %s", string(respBody))
	}

	var account struct {
		ID      string `json:"id"`
		Address string `json:"address"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (p *MailTMProvider) login(address, password string) (string, error) {
	payload := map[string]string{
		"address":  address,
		"password": password,
	}
	body, _ := json.Marshal(payload)

	resp, err := p.client.Post(p.baseURL+"/token", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("login failed")
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Token, nil
}

type SimpleLoginProvider struct {
	apiKey  string
	client  *http.Client
	baseURL string
}

func NewSimpleLoginProvider(apiKey string) *SimpleLoginProvider {
	return &SimpleLoginProvider{
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 30 * time.Second},
		baseURL: "https://app.simplelogin.io/api",
	}
}

func (p *SimpleLoginProvider) Name() string {
	return "simplelogin"
}

func (p *SimpleLoginProvider) Create(forwardTo string) (*ScratchEmail, error) {
	payload := map[string]interface{}{
		"alias_prefix": "sprayer",
		"mailbox_id":   nil,
		"note":         "Sprayer job application",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", p.baseURL+"/v2/alias/custom/new", bytes.NewReader(body))
	req.Header.Set("Authentication", p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create alias: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create alias failed: %s", string(respBody))
	}

	var result struct {
		ID      int    `json:"id"`
		Email   string `json:"email"`
		Enabled bool   `json:"enabled"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &ScratchEmail{
		ID:           fmt.Sprintf("%d", result.ID),
		EmailAddress: result.Email,
		Provider:     p.Name(),
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Time{},
		ForwardTo:    forwardTo,
		Active:       true,
		providerData: map[string]interface{}{
			"alias_id": result.ID,
		},
	}, nil
}

func (p *SimpleLoginProvider) CheckInbox(id string) ([]Email, error) {
	req, _ := http.NewRequest("GET", p.baseURL+"/v2/aliases/"+id+"/activities", nil)
	req.Header.Set("Authentication", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get activities: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Activities []struct {
			Action    string    `json:"action"`
			From      string    `json:"from"`
			To        string    `json:"to"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"activities"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	var emails []Email
	for _, act := range result.Activities {
		if act.Action == "forward" {
			emails = append(emails, Email{
				ID:         fmt.Sprintf("%d", act.Timestamp.UnixNano()),
				From:       act.From,
				To:         act.To,
				Subject:    "Forwarded email",
				ReceivedAt: act.Timestamp,
				Read:       false,
			})
		}
	}

	return emails, nil
}

func (p *SimpleLoginProvider) Deactivate(id string) error {
	req, _ := http.NewRequest("POST", p.baseURL+"/v2/aliases/"+id+"/toggle", nil)
	req.Header.Set("Authentication", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("deactivate alias: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("deactivate failed: %d", resp.StatusCode)
	}

	return nil
}

func generateRandomPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 16)
	rand.Read(b)
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

func (s *ScratchEmail) IsExpired() bool {
	if s.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(s.ExpiresAt)
}

func (s *ScratchEmail) TimeUntilExpiry() time.Duration {
	if s.ExpiresAt.IsZero() {
		return 0
	}
	return time.Until(s.ExpiresAt)
}

func (s *ScratchEmail) StatusText() string {
	if !s.Active {
		return "Deactivated"
	}
	if s.IsExpired() {
		return "Expired"
	}
	if s.ExpiresAt.IsZero() {
		return "Active (no expiry)"
	}
	remaining := s.TimeUntilExpiry()
	if remaining < time.Hour {
		return fmt.Sprintf("Expires in %.0f minutes", remaining.Minutes())
	}
	return fmt.Sprintf("Expires in %.1f hours", remaining.Hours())
}
