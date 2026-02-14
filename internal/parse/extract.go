package parse

import "strings"

// ExtractEmails returns all email addresses found in text.
func ExtractEmails(text string) []string {
	return EmailPattern.FindAllString(text, -1)
}

// ExtractFirstEmail returns the first email found, or empty string.
func ExtractFirstEmail(text string) string {
	match := EmailPattern.FindString(text)
	return match
}

// ExtractSalary returns the first salary range found.
func ExtractSalary(text string) string {
	return SalaryPattern.FindString(text)
}

// ExtractLocations returns all location mentions.
func ExtractLocations(text string) []string {
	matches := LocationPattern.FindAllString(text, -1)
	var unique []string
	seen := make(map[string]bool)
	for _, m := range matches {
		m = strings.TrimSpace(m)
		lower := strings.ToLower(m)
		if !seen[lower] {
			seen[lower] = true
			unique = append(unique, m)
		}
	}
	return unique
}

// IsRemote checks if a job text indicates remote work.
func IsRemote(text string) bool {
	lower := strings.ToLower(text)
	return strings.Contains(lower, "remote")
}

// ExtractURLs returns all URLs found in text.
func ExtractURLs(text string) []string {
	return URLPattern.FindAllString(text, -1)
}
