package parse

func findAll(rule string, text string) []string {
	var results []string
	for i := 0; i < len(text); {
		res, err := Parse("", []byte(text[i:]), Entrypoint(rule))
		if err == nil {
			s, ok := res.(string)
			if ok && s != "" {
				results = append(results, s)
				i += len(s)
				continue
			}
		}
		i++
	}
	return results
}

func findFirst(rule string, text string) string {
	for i := 0; i < len(text); i++ {
		res, err := Parse("", []byte(text[i:]), Entrypoint(rule))
		if err == nil {
			s, ok := res.(string)
			if ok && s != "" {
				return s
			}
		}
	}
	return ""
}

func exists(rule string, text string) bool {
	return findFirst(rule, text) != ""
}

func (c *current) textSlice() string {
	return string(c.text)
}
