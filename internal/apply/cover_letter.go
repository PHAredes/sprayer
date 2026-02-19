package apply

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"sprayer/internal/job"
	"sprayer/internal/llm"
	"sprayer/internal/profile"
)

type CoverLetterManager struct {
	client     *llm.Client
	cache      map[string]*CachedCoverLetter
	cacheDir   string
	mu         sync.RWMutex
	persistent bool
}

type CachedCoverLetter struct {
	Content      string    `json:"content"`
	Edited       bool      `json:"edited"`
	Generated    time.Time `json:"generated"`
	JobID        string    `json:"job_id"`
	JobTitle     string    `json:"job_title"`
	Company      string    `json:"company"`
	OriginalHash string    `json:"original_hash"`
}

type CoverLetterCacheData struct {
	CoverLetters map[string]*CachedCoverLetter `json:"cover_letters"`
	SavedAt      time.Time                     `json:"saved_at"`
}

func NewCoverLetterManager(client *llm.Client) *CoverLetterManager {
	return &CoverLetterManager{
		client:     client,
		cache:      make(map[string]*CachedCoverLetter),
		persistent: false,
	}
}

func NewPersistentCoverLetterManager(client *llm.Client, cacheDir string) (*CoverLetterManager, error) {
	m := &CoverLetterManager{
		client:     client,
		cache:      make(map[string]*CachedCoverLetter),
		cacheDir:   cacheDir,
		persistent: true,
	}

	if err := m.loadFromDisk(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load cover letter cache: %w", err)
	}

	return m, nil
}

func (m *CoverLetterManager) loadFromDisk() error {
	if m.cacheDir == "" {
		return nil
	}

	cachePath := filepath.Join(m.cacheDir, "cover_letters.json")
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return err
	}

	var cacheData CoverLetterCacheData
	if err := json.Unmarshal(data, &cacheData); err != nil {
		return fmt.Errorf("unmarshal cache: %w", err)
	}

	m.mu.Lock()
	m.cache = cacheData.CoverLetters
	m.mu.Unlock()

	return nil
}

func (m *CoverLetterManager) saveToDisk() error {
	if !m.persistent || m.cacheDir == "" {
		return nil
	}

	if err := os.MkdirAll(m.cacheDir, 0755); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}

	m.mu.RLock()
	cacheData := CoverLetterCacheData{
		CoverLetters: m.cache,
		SavedAt:      time.Now(),
	}
	m.mu.RUnlock()

	data, err := json.MarshalIndent(cacheData, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal cache: %w", err)
	}

	cachePath := filepath.Join(m.cacheDir, "cover_letters.json")
	return os.WriteFile(cachePath, data, 0644)
}

func (m *CoverLetterManager) GetCoverLetter(j *job.Job, p *profile.Profile) (string, error) {
	if p.CoverPath != "" {
		content, err := os.ReadFile(p.CoverPath)
		if err == nil {
			return string(content), nil
		}
	}

	return m.Generate(j, p)
}

func (m *CoverLetterManager) Generate(j *job.Job, p *profile.Profile) (string, error) {
	if m.client == nil || !m.client.Available() {
		return "", fmt.Errorf("LLM client not available")
	}

	m.mu.RLock()
	if cached, ok := m.cache[j.ID]; ok {
		if time.Since(cached.Generated) < 24*time.Hour && !cached.Edited {
			m.mu.RUnlock()
			return cached.Content, nil
		}
	}
	m.mu.RUnlock()

	content, err := GenerateCoverLetter(j, p, m.client)
	if err != nil {
		return "", err
	}

	m.mu.Lock()
	m.cache[j.ID] = &CachedCoverLetter{
		Content:   content,
		Edited:    false,
		Generated: time.Now(),
		JobID:     j.ID,
		JobTitle:  j.Title,
		Company:   j.Company,
	}
	m.mu.Unlock()

	m.saveToDisk()

	return content, nil
}

func (m *CoverLetterManager) Edit(jobID, newContent string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cached, ok := m.cache[jobID]; ok {
		cached.Content = newContent
		cached.Edited = true
		m.saveToDisk()
		return nil
	}

	m.cache[jobID] = &CachedCoverLetter{
		Content:   newContent,
		Edited:    true,
		Generated: time.Now(),
		JobID:     jobID,
	}

	m.saveToDisk()
	return nil
}

func (m *CoverLetterManager) GetCached(jobID string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if cached, ok := m.cache[jobID]; ok {
		return cached.Content, true
	}
	return "", false
}

func (m *CoverLetterManager) GetCachedMeta(jobID string) (*CachedCoverLetter, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if cached, ok := m.cache[jobID]; ok {
		return cached, true
	}
	return nil, false
}

func (m *CoverLetterManager) ClearCache() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache = make(map[string]*CachedCoverLetter)
	return m.saveToDisk()
}

func (m *CoverLetterManager) ListCached() []CachedCoverLetter {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []CachedCoverLetter
	for _, cached := range m.cache {
		list = append(list, *cached)
	}
	return list
}

func (m *CoverLetterManager) Available() bool {
	return m.client != nil && m.client.Available()
}

func (m *CoverLetterManager) IsCached(jobID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.cache[jobID]
	return ok
}

func (m *CoverLetterManager) RemoveCached(jobID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.cache, jobID)
	return m.saveToDisk()
}

func SaveCoverLetter(content, jobID, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("create output dir: %w", err)
	}

	filename := fmt.Sprintf("cover_letter_%s_%d.txt", sanitize(jobID), time.Now().Unix())
	filePath := filepath.Join(outputDir, filename)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("write cover letter: %w", err)
	}

	return filePath, nil
}

func LoadCoverLetter(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read cover letter: %w", err)
	}
	return string(content), nil
}

func CombineWithEmail(emailBody, coverLetter string) string {
	emailBody = strings.TrimSpace(emailBody)
	coverLetter = strings.TrimSpace(coverLetter)

	if emailBody == "" {
		return coverLetter
	}
	if coverLetter == "" {
		return emailBody
	}

	return emailBody + "\n\n---\n\n" + coverLetter
}
