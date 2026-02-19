package apply

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"sprayer/internal/job"
	"sprayer/internal/profile"
)

func Draft(j job.Job, p profile.Profile, subject, body string) (string, error) {
	maildirPath := filepath.Join(os.Getenv("HOME"), "Maildir", "drafts", "new")
	if err := os.MkdirAll(maildirPath, 0755); err != nil {
		return "", fmt.Errorf("create drafts dir: %w", err)
	}

	to := j.Email
	if to == "" {
		return "", fmt.Errorf("no email address for job %s", j.ID)
	}

	filename := fmt.Sprintf("%d.sprayer.%s", time.Now().Unix(), sanitize(j.ID))
	draftPath := filepath.Join(maildirPath, filename)

	var attachmentPart string
	cvPDF := findPDF(p.CVPath)
	if cvPDF != "" {
		pdfData, err := os.ReadFile(cvPDF)
		if err == nil {
			encoded := base64.StdEncoding.EncodeToString(pdfData)
			attachmentPart = fmt.Sprintf(`
--%s
Content-Type: application/pdf
Content-Disposition: attachment; filename="%s"
Content-Transfer-Encoding: base64

%s`, boundary, filepath.Base(cvPDF), wrapBase64(encoded))
		}
	}

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s\n", p.ContactEmail))
	msg.WriteString(fmt.Sprintf("To: %s\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\n", subject))
	msg.WriteString(fmt.Sprintf("Date: %s\n", time.Now().Format(time.RFC1123Z)))
	msg.WriteString("MIME-Version: 1.0\n")

	if attachmentPart != "" {
		msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\n\n", boundary))
		msg.WriteString(fmt.Sprintf("--%s\n", boundary))
		msg.WriteString("Content-Type: text/plain; charset=utf-8\n\n")
		msg.WriteString(body)
		msg.WriteString("\n")
		msg.WriteString(attachmentPart)
		msg.WriteString(fmt.Sprintf("\n--%s--\n", boundary))
	} else {
		msg.WriteString("Content-Type: text/plain; charset=utf-8\n\n")
		msg.WriteString(body)
		msg.WriteString("\n")
	}

	if err := os.WriteFile(draftPath, []byte(msg.String()), 0644); err != nil {
		return "", fmt.Errorf("write draft: %w", err)
	}

	return draftPath, nil
}

const boundary = "sprayer-boundary"

func findPDF(texPath string) string {
	if texPath == "" {
		return ""
	}
	pdf := strings.TrimSuffix(texPath, filepath.Ext(texPath)) + ".pdf"
	if _, err := os.Stat(pdf); err == nil {
		return pdf
	}
	return ""
}

func sanitize(s string) string {
	r := strings.NewReplacer("/", "_", " ", "_", ":", "_")
	return r.Replace(s)
}

func wrapBase64(s string) string {
	var out strings.Builder
	for i := 0; i < len(s); i += 76 {
		end := i + 76
		if end > len(s) {
			end = len(s)
		}
		out.WriteString(s[i:end])
		out.WriteString("\n")
	}
	return out.String()
}

func EditDraft(draftPath string) error {
	editor := getEditor()
	if editor == "" {
		editor = "emacs"
	}

	cmd := exec.Command(editor, draftPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}

	return nil
}

func EditAndSend(draftPath string) error {
	if err := EditDraft(draftPath); err != nil {
		return err
	}

	content, err := os.ReadFile(draftPath)
	if err != nil {
		return fmt.Errorf("read draft: %w", err)
	}

	draft, err := parseDraftContent(string(content), draftPath)
	if err != nil {
		return fmt.Errorf("parse draft: %w", err)
	}

	fmt.Printf("\nDraft ready to send:\n")
	fmt.Printf("  To: %s\n", draft.To)
	fmt.Printf("  Subject: %s\n", draft.Subject)
	fmt.Printf("\nSend? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response == "y" || response == "yes" {
		return sendViaSMTP(draft)
	}

	fmt.Println("Draft saved, not sent.")
	return nil
}

func parseDraftContent(content, draftPath string) (*EmailDraft, error) {
	draft := &EmailDraft{
		DraftPath: draftPath,
		ID:        filepath.Base(draftPath),
		CreatedAt: time.Now(),
	}

	lines := strings.Split(content, "\n")
	inBody := false
	var bodyLines []string

	for _, line := range lines {
		if inBody {
			bodyLines = append(bodyLines, line)
			continue
		}

		if line == "" {
			inBody = true
			continue
		}

		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue
		}

		switch strings.ToLower(parts[0]) {
		case "from":
			draft.From = parts[1]
		case "to":
			draft.To = parts[1]
		case "subject":
			draft.Subject = parts[1]
		case "date":
			draft.CreatedAt, _ = time.Parse(time.RFC1123Z, parts[1])
		}
	}

	draft.Body = strings.Join(bodyLines, "\n")
	return draft, nil
}
