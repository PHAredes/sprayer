package llm

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// LoadPrompt reads a prompt file from the prompts/ directory and interpolates variables.
// Variables use {{name}} syntax.
func LoadPrompt(name string, vars map[string]string) (string, error) {
	content, err := readPromptFile(name)
	if err != nil {
		return "", err
	}

	return Interpolate(content, vars), nil
}

// Interpolate replaces {{key}} placeholders with values from vars.
func Interpolate(template string, vars map[string]string) string {
	for k, v := range vars {
		template = strings.ReplaceAll(template, "{{"+k+"}}", v)
	}
	return template
}

func readPromptFile(name string) (string, error) {
	// Find prompts dir relative to the project root.
	// Try: ./prompts/, then relative to this source file.
	candidates := []string{
		filepath.Join("prompts", name+".txt"),
	}

	// Also try relative to the binary location.
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "prompts", name+".txt"))
	}

	// Also try relative to source (for dev).
	_, thisFile, _, ok := runtime.Caller(0)
	if ok {
		projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
		candidates = append(candidates, filepath.Join(projectRoot, "prompts", name+".txt"))
	}

	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err == nil {
			return string(data), nil
		}
	}

	return "", fmt.Errorf("prompt file not found: %s.txt", name)
}
