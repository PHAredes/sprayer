package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestE2E_CLI_Help(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	cmd := exec.Command("go", "run", "./cmd/sprayer")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "HOME="+os.TempDir())
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI help failed: %v\nOutput: %s", err, string(out))
	}

	output := string(out)
	if !strings.Contains(output, "scrape") || !strings.Contains(output, "list") {
		t.Errorf("Expected usage to contain 'scrape' and 'list', got: %s", output)
	}
}

func TestE2E_CLI_List_Empty(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	cmd := exec.Command("go", "run", "./cmd/sprayer", "list")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "HOME="+os.TempDir())
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI list failed: %v\nOutput: %s", err, string(out))
	}
}
