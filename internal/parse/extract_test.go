package parse_test

import (
	"testing"
	"testing/quick"

	"sprayer/internal/parse"
)

func TestExtractEmails_Property(t *testing.T) {
	f := func(user, domain string) bool {
		// Filter out strings that would make invalid emails or break regex expectations
		if user == "" || domain == "" || len(user) > 20 || len(domain) > 20 {
			return true
		}
		// Basic ascii check for property test stability
		for _, r := range user + domain {
			if r < 'a' || r > 'z' {
				return true
			}
		}
		// Construct a valid email
		email := user + "@" + domain + ".com"
		// Inject it into random garbage text
		text := "Some garbage " + email + " more garbage"
		
		emails := parse.ExtractEmails(text)
		for _, e := range emails {
			if e == email {
				return true
			}
		}
		return false
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestExtractSalary_Property(t *testing.T) {
	// Simple property: extracted salary should be a substring of original text
	f := func(text string) bool {
		salary := parse.ExtractSalary(text)
		if salary == "" {
			return true
		}
		// In Go, strings.Contains is enough check
		for i := 0; i < len(text)-len(salary)+1; i++ {
			if text[i:i+len(salary)] == salary {
				return true
			}
		}
		return false
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
