package parse

import (
	"strings"
)

// Traps contains phrases often used to filter out AI/bot applications.
var Traps = []string{
	"ignore all previous instructions",
	"ignore previous instructions",
	"write a cover letter about",
	"start your cover letter with",
	"include the word",
	"no ai",
	"no llm",
	"no chatgpt",
	"human only",
	"not a bot",
	"solve this math",
	"brown m&m",
	"blue m&m",
	"banana", // sometimes used as a codeword
	"mention this word",
	"add the sha256",
	"add the result of",
	"word count",
	"how many times does the letter",
	"recipe for",
}

// CheckForTraps scans the text for anti-AI phrases or instructions.
// Returns a list of potential traps found.
func CheckForTraps(text string) []string {
	text = strings.ToLower(text)
	var found []string
	
	for _, trap := range Traps {
		if strings.Contains(text, trap) {
			found = append(found, trap)
		}
	}
	
	return found
}
