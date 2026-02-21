package scraper

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"sprayer/internal/job"
	"sprayer/internal/parse"
)

// HN scrapes the monthly "Who is Hiring?" thread via the HN Algolia API.
func HN() job.Scraper {
	return func() ([]job.Job, error) {
		// Find the latest "Who is Hiring?" story
		storyURL := "https://hn.algolia.com/api/v1/search?query=%22Ask%20HN%3A%20Who%20is%20hiring%22&tags=story&hitsPerPage=1"
		storyResp, err := httpGet(storyURL)
		if err != nil {
			return nil, fmt.Errorf("HN story search: %w", err)
		}

		var storyResult struct {
			Hits []struct {
				ObjectID string `json:"objectID"`
				Title    string `json:"title"`
			} `json:"hits"`
		}
		if err := json.Unmarshal(storyResp, &storyResult); err != nil {
			return nil, fmt.Errorf("HN story parse: %w", err)
		}
		if len(storyResult.Hits) == 0 {
			return nil, fmt.Errorf("no 'Who is Hiring?' thread found")
		}

		storyID := storyResult.Hits[0].ObjectID

		// Fetch top-level comments (job postings)
		var allJobs []job.Job
		page := 0
		for {
			commentsURL := fmt.Sprintf(
				"https://hn.algolia.com/api/v1/search?tags=comment,story_%s&hitsPerPage=100&page=%d",
				storyID, page,
			)
			commentsResp, err := httpGet(commentsURL)
			if err != nil {
				break
			}

			var commentsResult struct {
				Hits []struct {
					ObjectID  string `json:"objectID"`
					CommentText string `json:"comment_text"`
					CreatedAt string `json:"created_at"`
				} `json:"hits"`
				NbPages int `json:"nbPages"`
			}
			if err := json.Unmarshal(commentsResp, &commentsResult); err != nil {
				break
			}

			for _, hit := range commentsResult.Hits {
				text := hit.CommentText
				if len(text) < 50 {
					continue // Skip very short comments (replies, not job posts)
				}

				j := parseHNComment(hit.ObjectID, text, hit.CreatedAt)
				allJobs = append(allJobs, j)
			}

			page++
			if page >= commentsResult.NbPages || page >= 5 {
				break
			}
			time.Sleep(200 * time.Millisecond) // Rate limit
		}

		return allJobs, nil
	}
}

func parseHNComment(id, text, createdAt string) job.Job {
	// HN job posts typically start with "Company | Role | Location"
	lines := strings.SplitN(text, "\n", 2)
	header := lines[0]

	// Try to parse pipe-separated header
	parts := strings.Split(header, "|")
	var company, title, location string
	if len(parts) >= 3 {
		company = strings.TrimSpace(parts[0])
		title = strings.TrimSpace(parts[1])
		location = strings.TrimSpace(parts[2])
	} else if len(parts) == 2 {
		company = strings.TrimSpace(parts[0])
		title = strings.TrimSpace(parts[1])
	} else {
		// Fallback: first line is the title
		title = strings.TrimSpace(header)
	}

	// Strip HTML tags from description
	desc := stripHTML(text)

	// Extract email and salary via parsing
	email := parse.ExtractFirstEmail(desc)
	salary := parse.ExtractSalary(desc)

	posted, _ := time.Parse(time.RFC3339, createdAt)

	return job.Job{
		ID:          fmt.Sprintf("hn-%s", id),
		Title:       title,
		Company:     company,
		Location:    location,
		Description: desc,
		URL:         fmt.Sprintf("https://news.ycombinator.com/item?id=%s", id),
		Source:      "hackernews",
		PostedDate:  posted,
		Salary:      salary,
		Email:       email,
		Score:       scoreHNJob(title, desc),
	}
}

func scoreHNJob(title, desc string) int {
	score := 50
	lower := strings.ToLower(title + " " + desc)

	// Boost for desirable keywords
	boosts := map[string]int{
		"remote": 10, "golang": 10, "go": 5, "rust": 10,
		"compiler": 15, "haskell": 10, "functional": 8,
		"senior": 5, "staff": 8, "principal": 10,
	}
	for kw, boost := range boosts {
		if strings.Contains(lower, kw) {
			score += boost
		}
	}

	if score > 100 {
		score = 100
	}
	return score
}

func stripHTML(s string) string {
	// Simple HTML tag stripper
	var out strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			out.WriteRune(r)
		}
	}
	return strings.TrimSpace(out.String())
}

func httpGet(url string) ([]byte, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	var buf [1 << 20]byte // 1MB max
	n := 0
	for {
		read, err := resp.Body.Read(buf[n:])
		n += read
		if err != nil {
			break
		}
	}
	return buf[:n], nil
}

// idFromContent generates a deterministic ID from content.
func idFromContent(source, content string) string {
	h := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%s-%x", source, h[:8])
}
