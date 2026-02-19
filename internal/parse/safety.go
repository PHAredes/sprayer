package parse

import "strings"

func CheckForTraps(text string) []string {
	return findAll("TrapEntry", text)
}

func Sanitize(text string) string {
	traps := CheckForTraps(text)
	for _, trap := range traps {
		// Use a case-insensitive replacement if possible, 
		// but since we have the exact match from findAll, we can just replace it.
		// Note: findAll returns exactly what was matched.
		text = strings.ReplaceAll(text, trap, "[FLAGGED CONTENT REMOVED]")
	}
	return text
}
