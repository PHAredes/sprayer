package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

func (c *CLI) handleSetup() {
	var (
		smtpHost string = os.Getenv("SPRAYER_SMTP_HOST")
		smtpPort string = os.Getenv("SPRAYER_SMTP_PORT")
		smtpUser string = os.Getenv("SPRAYER_SMTP_USER")
		smtpPass string = os.Getenv("SPRAYER_SMTP_PASS")
		smtpFrom string = os.Getenv("SPRAYER_SMTP_FROM")
		llmKey   string = os.Getenv("SPRAYER_LLM_KEY")
		llmURL   string = os.Getenv("SPRAYER_LLM_URL")
		llmModel string = os.Getenv("SPRAYER_LLM_MODEL")
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Sprayer Setup").
				Description("Configure your email (SMTP) and LLM settings."),
				
			huh.NewInput().
				Title("SMTP Host").
				Value(&smtpHost).
				Placeholder("smtp.gmail.com"),

			huh.NewInput().
				Title("SMTP Port").
				Value(&smtpPort).
				Placeholder("587"),

			huh.NewInput().
				Title("SMTP User").
				Value(&smtpUser).
				Placeholder("me@example.com"),
			
			huh.NewInput().
				Title("SMTP Password").
				Value(&smtpPass).
				EchoMode(huh.EchoModePassword),

			huh.NewInput().
				Title("SMTP From Address").
				Value(&smtpFrom).
				Placeholder("Valid Name <me@example.com>"),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("LLM API Key").
				Value(&llmKey).
				EchoMode(huh.EchoModePassword),
			
			huh.NewInput().
				Title("LLM Base URL").
				Value(&llmURL).
				Placeholder("https://api.openai.com/v1"),
				
			huh.NewInput().
				Title("LLM Model").
				Value(&llmModel).
				Placeholder("gpt-4o"),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println("Setup cancelled.")
		return
	}

	content := fmt.Sprintf(`SPRAYER_SMTP_HOST=%s
SPRAYER_SMTP_PORT=%s
SPRAYER_SMTP_USER=%s
SPRAYER_SMTP_PASS=%s
SPRAYER_SMTP_FROM=%s
SPRAYER_LLM_KEY=%s
SPRAYER_LLM_URL=%s
SPRAYER_LLM_MODEL=%s
`, smtpHost, smtpPort, smtpUser, smtpPass, smtpFrom, llmKey, llmURL, llmModel)

	err = os.WriteFile(".env", []byte(content), 0600)
	if err != nil {
		fmt.Printf("Error configuring .env: %v\n", err)
		return
	}

	fmt.Println("Configuration saved to .env")
}
