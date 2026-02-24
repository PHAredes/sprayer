package apply

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"github.com/jordan-wright/email"
)

// SendDirect sends an email immediately using SMTP configuration.
// It mimics the behavior of tools like 'pop'.
func SendDirect(to, subject, body, attachmentPath string) error {
	host := os.Getenv("SPRAYER_SMTP_HOST")
	port := os.Getenv("SPRAYER_SMTP_PORT")
	username := os.Getenv("SPRAYER_SMTP_USER")
	password := os.Getenv("SPRAYER_SMTP_PASS")
	from := os.Getenv("SPRAYER_SMTP_FROM")

	if host == "" || username == "" || password == "" {
		return fmt.Errorf("SMTP configuration missing (SPRAYER_SMTP_HOST, USER, PASS)")
	}
	if from == "" {
		from = username
	}
	if port == "" {
		port = "587"
	}

	e := email.NewEmail()
	e.From = from
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(body)
	
	// Basic HTML conversion (wrapping body in pre/div)
	// In a real 'pop' like tool we would render markdown.
	htmlBody := fmt.Sprintf("<html><body><pre style='font-family: sans-serif'>%s</pre></body></html>", body)
	e.HTML = []byte(htmlBody)

	if attachmentPath != "" {
		if _, err := e.AttachFile(attachmentPath); err != nil {
			return fmt.Errorf("attach file: %w", err)
		}
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	auth := smtp.PlainAuth("", username, password, host)

	// Start TLS if port is 587 or 465
	var err error
	if port == "465" {
		// SSL/TLS
		err = e.SendWithTLS(addr, auth, &tls.Config{ServerName: host})
	} else {
		// StartTLS (587) or Plain (25)
		err = e.Send(addr, auth)
	}

	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}
