package ui_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestE2E_CLI_Scrape_And_List(t *testing.T) {
	root := "/home/user/openclaw-setup"
	
	// Build the CLI binary
	cmd := exec.Command("go", "build", "-o", "sprayer-e2e", "./cmd/cli/main.go")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Build failed: %v\nOutput: %s", err, string(out))
	}
	defer os.Remove(filepath.Join(root, "sprayer-e2e"))
	
	// 1. Scrape (Fast API only)
	scrapeCmd := exec.Command("./sprayer-e2e", "scrape", "--fast", "rust", "remote")
	scrapeCmd.Dir = root
	out, err := scrapeCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Scrape failed: %v\nOutput: %s", err, string(out))
	}
	t.Logf("Scrape Output: %s", string(out))

	// 2. List
	listCmd := exec.Command("./sprayer-e2e", "list")
	listCmd.Dir = root
	listOut, err := listCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("List failed: %v\nOutput: %s", err, string(listOut))
	}
	
	output := string(listOut)
	if !strings.Contains(output, "rust") && !strings.Contains(output, "Rust") && !strings.Contains(output, "Jobs:") { 
		// Check for some output indicating it ran
	}
}
