package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	EnvLLMURL   = "SPRAYER_LLM_URL"
	EnvLLMKey   = "SPRAYER_LLM_KEY"
	EnvLLMModel = "SPRAYER_LLM_MODEL"
)

type Client struct {
	baseURL string
	apiKey  string
	model   string
	http    *http.Client
}

func NewClient() *Client {
	baseURL := os.Getenv(EnvLLMURL)
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	model := os.Getenv(EnvLLMModel)
	if model == "" {
		model = "kimi-k2"
	}
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  os.Getenv(EnvLLMKey),
		model:   model,
		http:    &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *Client) Available() bool {
	return c.apiKey != ""
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *Client) Complete(system, user string) (string, error) {
	if !c.Available() {
		return "", fmt.Errorf("LLM not configured: set SPRAYER_LLM_KEY")
	}

	req := chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("LLM request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result chatResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("LLM response parse error: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("LLM error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("LLM returned no choices")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}
