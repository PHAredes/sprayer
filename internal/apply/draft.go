package apply

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sprayer/internal/job"
	"sprayer/internal/profile"
)

// Draft generates a Maildir-format email draft file for mu4e.
func Draft(j job.Job, p profile.Profile, subject, body string) (string, error) {
	maildirPath := filepath.Join(os.Getenv("HOME"), "Maildir", "drafts", "new")
	if err := os.MkdirAll(maildirPath, 0755); err != nil {
		return "", fmt.Errorf("create drafts dir: %w", err)
	}

	// Determine recipient
	to := j.Email
	if to == "" {
		return "", fmt.Errorf("no email address for job %s", j.ID)
	}

	filename := fmt.Sprintf("%d.sprayer.%s", time.Now().Unix(), sanitize(j.ID))
	draftPath := filepath.Join(maildirPath, filename)

	// Try to attach CV PDF
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

	// Build the email
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

// findPDF looks for a .pdf file alongside or derived from the given tex path.
func findPDF(texPath string) string {
	if texPath == "" {
		return ""
	}
	// Try the same name with .pdf extension
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
