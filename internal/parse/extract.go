package parse

import "strings"

func ExtractEmails(text string) []string {
	return findAll("EmailEntry", text)
}

func ExtractFirstEmail(text string) string {
	return findFirst("EmailEntry", text)
}

func ExtractSalary(text string) string {
	return findFirst("SalaryEntry", text)
}

func ExtractLocations(text string) []string {
	matches := findAll("LocationEntry", text)
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

func IsRemote(text string) bool {
	return exists("Remote", text)
}

func ExtractURLs(text string) []string {
	return findAll("URLEntry", text)
}
