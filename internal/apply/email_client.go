package apply

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/jordan-wright/email"
	"sprayer/internal/job"
	"sprayer/internal/profile"
)

type EmailDraft struct {
	ID              string
	From            string
	To              string
	Subject         string
	Body            string
	Attachments     []Attachment
	CreatedAt       time.Time
	DraftPath       string
	AttachmentPaths []string
}

type Attachment struct {
	Filename string
	Data     []byte
	MimeType string
}

type EmailClient interface {
	OpenDraft(draft *EmailDraft) error
	ListDrafts() ([]EmailDraft, error)
	Send(draft *EmailDraft) error
	CreateDraft(j job.Job, p profile.Profile, subject, body string) (*EmailDraft, error)
}

type EmailClientType string

const (
	ClientMu4e EmailClientType = "mu4e"
	ClientTUI  EmailClientType = "tui"
	ClientSMTP EmailClientType = "smtp"
)

func GetEmailClientType() EmailClientType {
	clientType := os.Getenv("SPRAYER_EMAIL_CLIENT")
	switch strings.ToLower(clientType) {
	case "mu4e":
		return ClientMu4e
	case "tui":
		return ClientTUI
	case "smtp":
		return ClientSMTP
	default:
		return ClientMu4e
	}
}

func NewEmailClient(clientType EmailClientType) EmailClient {
	switch clientType {
	case ClientMu4e:
		return NewMu4eClient()
	case ClientTUI:
		return NewTUIEmailClient()
	case ClientSMTP:
		return NewSMTPClient()
	default:
		return NewMu4eClient()
	}
}

type Mu4eClient struct {
	maildirPath string
	editor      string
}

func NewMu4eClient() *Mu4eClient {
	maildirPath := os.Getenv("SPRAYER_MAILDIR_PATH")
	if maildirPath == "" {
		maildirPath = filepath.Join(os.Getenv("HOME"), "Maildir")
	}

	editor := getEditor()

	return &Mu4eClient{
		maildirPath: maildirPath,
		editor:      editor,
	}
}

func (m *Mu4eClient) CreateDraft(j job.Job, p profile.Profile, subject, body string) (*EmailDraft, error) {
	to := j.Email
	if to == "" {
		return nil, fmt.Errorf("no email address for job %s", j.ID)
	}

	draft := &EmailDraft{
		ID:        fmt.Sprintf("%d.%s", time.Now().Unix(), sanitize(j.ID)),
		From:      p.ContactEmail,
		To:        to,
		Subject:   subject,
		Body:      body,
		CreatedAt: time.Now(),
	}

	cvPDF := findPDF(p.CVPath)
	if cvPDF != "" {
		draft.AttachmentPaths = []string{cvPDF}
	}

	return draft, nil
}

func (m *Mu4eClient) OpenDraft(draft *EmailDraft) error {
	draftsDir := filepath.Join(m.maildirPath, "drafts", "new")
	if err := os.MkdirAll(draftsDir, 0755); err != nil {
		return fmt.Errorf("create drafts dir: %w", err)
	}

	filename := fmt.Sprintf("%d.sprayer.%s", time.Now().Unix(), sanitize(draft.ID))
	draftPath := filepath.Join(draftsDir, filename)
	draft.DraftPath = draftPath

	content := m.buildEmailContent(draft)
	if err := os.WriteFile(draftPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("write draft: %w", err)
	}

	return m.openInEditor(draftPath)
}

func (m *Mu4eClient) openInEditor(draftPath string) error {
	editor := m.editor
	if editor == "" {
		editor = "emacs"
	}

	var cmd *exec.Cmd
	if strings.Contains(editor, "emacsclient") || editor == "emacs" {
		cmd = exec.Command(editor, "--eval",
			fmt.Sprintf("(find-file-other-window \"%s\")", draftPath))
	} else {
		cmd = exec.Command(editor, draftPath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}

	return nil
}

func (m *Mu4eClient) buildEmailContent(draft *EmailDraft) string {
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s\n", draft.From))
	msg.WriteString(fmt.Sprintf("To: %s\n", draft.To))
	msg.WriteString(fmt.Sprintf("Subject: %s\n", draft.Subject))
	msg.WriteString(fmt.Sprintf("Date: %s\n", draft.CreatedAt.Format(time.RFC1123Z)))
	msg.WriteString("MIME-Version: 1.0\n")

	if len(draft.AttachmentPaths) > 0 || len(draft.Attachments) > 0 {
		msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\n\n", boundary))
		msg.WriteString(fmt.Sprintf("--%s\n", boundary))
		msg.WriteString("Content-Type: text/plain; charset=utf-8\n\n")
		msg.WriteString(draft.Body)
		msg.WriteString("\n")

		for _, path := range draft.AttachmentPaths {
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			filename := filepath.Base(path)
			msg.WriteString(fmt.Sprintf("--%s\n", boundary))
			msg.WriteString("Content-Type: application/pdf\n")
			msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\n", filename))
			msg.WriteString("Content-Transfer-Encoding: base64\n\n")
			msg.WriteString(wrapBase64(base64.StdEncoding.EncodeToString(data)))
			msg.WriteString("\n")
		}

		for _, att := range draft.Attachments {
			msg.WriteString(fmt.Sprintf("--%s\n", boundary))
			msg.WriteString(fmt.Sprintf("Content-Type: %s\n", att.MimeType))
			msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\n", att.Filename))
			msg.WriteString("Content-Transfer-Encoding: base64\n\n")
			msg.WriteString(wrapBase64(base64.StdEncoding.EncodeToString(att.Data)))
			msg.WriteString("\n")
		}

		msg.WriteString(fmt.Sprintf("--%s--\n", boundary))
	} else {
		msg.WriteString("Content-Type: text/plain; charset=utf-8\n\n")
		msg.WriteString(draft.Body)
		msg.WriteString("\n")
	}

	return msg.String()
}

func (m *Mu4eClient) ListDrafts() ([]EmailDraft, error) {
	draftsDir := filepath.Join(m.maildirPath, "drafts", "new")
	entries, err := os.ReadDir(draftsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read drafts dir: %w", err)
	}

	var drafts []EmailDraft
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(draftsDir, entry.Name())
		draft, err := m.parseDraftFile(path)
		if err != nil {
			continue
		}
		drafts = append(drafts, *draft)
	}

	return drafts, nil
}

func (m *Mu4eClient) parseDraftFile(path string) (*EmailDraft, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	draft := &EmailDraft{
		DraftPath: path,
	}

	lines := strings.Split(string(content), "\n")
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
	draft.ID = filepath.Base(path)

	return draft, nil
}

func (m *Mu4eClient) Send(draft *EmailDraft) error {
	return sendViaSMTP(draft)
}

type TUIEmailClient struct{}

func NewTUIEmailClient() *TUIEmailClient {
	return &TUIEmailClient{}
}

func (t *TUIEmailClient) CreateDraft(j job.Job, p profile.Profile, subject, body string) (*EmailDraft, error) {
	to := j.Email
	if to == "" {
		return nil, fmt.Errorf("no email address for job %s", j.ID)
	}

	draft := &EmailDraft{
		ID:        fmt.Sprintf("%d.%s", time.Now().Unix(), sanitize(j.ID)),
		From:      p.ContactEmail,
		To:        to,
		Subject:   subject,
		Body:      body,
		CreatedAt: time.Now(),
	}

	cvPDF := findPDF(p.CVPath)
	if cvPDF != "" {
		draft.AttachmentPaths = []string{cvPDF}
	}

	return draft, nil
}

func (t *TUIEmailClient) OpenDraft(draft *EmailDraft) error {
	return fmt.Errorf("TUI editing should be handled by EmailComposer component")
}

func (t *TUIEmailClient) ListDrafts() ([]EmailDraft, error) {
	return nil, fmt.Errorf("draft listing not supported in TUI mode")
}

func (t *TUIEmailClient) Send(draft *EmailDraft) error {
	return sendViaSMTP(draft)
}

type SMTPClient struct{}

func NewSMTPClient() *SMTPClient {
	return &SMTPClient{}
}

func (s *SMTPClient) CreateDraft(j job.Job, p profile.Profile, subject, body string) (*EmailDraft, error) {
	to := j.Email
	if to == "" {
		return nil, fmt.Errorf("no email address for job %s", j.ID)
	}

	draft := &EmailDraft{
		ID:        fmt.Sprintf("%d.%s", time.Now().Unix(), sanitize(j.ID)),
		From:      p.ContactEmail,
		To:        to,
		Subject:   subject,
		Body:      body,
		CreatedAt: time.Now(),
	}

	cvPDF := findPDF(p.CVPath)
	if cvPDF != "" {
		draft.AttachmentPaths = []string{cvPDF}
	}

	return draft, nil
}

func (s *SMTPClient) OpenDraft(draft *EmailDraft) error {
	return nil
}

func (s *SMTPClient) ListDrafts() ([]EmailDraft, error) {
	return nil, fmt.Errorf("draft listing not supported in SMTP mode")
}

func (s *SMTPClient) Send(draft *EmailDraft) error {
	return sendViaSMTP(draft)
}

func sendViaSMTP(draft *EmailDraft) error {
	host := os.Getenv("SPRAYER_SMTP_HOST")
	port := os.Getenv("SPRAYER_SMTP_PORT")
	username := os.Getenv("SPRAYER_SMTP_USER")
	password := os.Getenv("SPRAYER_SMTP_PASS")
	from := os.Getenv("SPRAYER_SMTP_FROM")

	if host == "" || username == "" || password == "" {
		return fmt.Errorf("SMTP configuration missing (SPRAYER_SMTP_HOST, USER, PASS)")
	}
	if from == "" {
		from = draft.From
		if from == "" {
			from = username
		}
	}
	if port == "" {
		port = "587"
	}

	e := email.NewEmail()
	e.From = from
	e.To = []string{draft.To}
	e.Subject = draft.Subject
	e.Text = []byte(draft.Body)

	htmlBody := fmt.Sprintf("<html><body><pre style='font-family: sans-serif'>%s</pre></body></html>", draft.Body)
	e.HTML = []byte(htmlBody)

	for _, path := range draft.AttachmentPaths {
		if _, err := e.AttachFile(path); err != nil {
			return fmt.Errorf("attach file %s: %w", path, err)
		}
	}

	for _, att := range draft.Attachments {
		if _, err := e.Attach(bytes.NewReader(att.Data), att.Filename, att.MimeType); err != nil {
			return fmt.Errorf("attach %s: %w", att.Filename, err)
		}
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	auth := smtp.PlainAuth("", username, password, host)

	var err error
	if port == "465" {
		err = e.SendWithTLS(addr, auth, &tls.Config{ServerName: host})
	} else {
		err = e.Send(addr, auth)
	}

	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

func getEditor() string {
	if editor := os.Getenv("SPRAYER_EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	return "emacs"
}
