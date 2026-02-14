package parse

import "regexp"

// Common regex patterns for extracting structured data from job text.
var (
	EmailPattern    = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	SalaryPattern   = regexp.MustCompile(`(?i)(\$[\d,]+(?:k)?(?:\s*[-–]\s*\$[\d,]+(?:k)?)?(?:\s*(?:per\s+)?(?:year|yr|annually|pa|p\.a\.))?|[\d,]+(?:k)?\s*[-–]\s*[\d,]+(?:k)?\s*(?:USD|EUR|GBP|BRL))`)
	LocationPattern = regexp.MustCompile(`(?i)(remote|on-?site|hybrid|(?:[A-Z][a-z]+(?:\s[A-Z][a-z]+)*,\s*[A-Z]{2,}))`)
	URLPattern      = regexp.MustCompile(`https?://[^\s<>"]+`)
)
