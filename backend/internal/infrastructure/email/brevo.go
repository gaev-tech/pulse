package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Sender sends transactional email.
type Sender interface {
	SendMagicLink(ctx context.Context, toEmail, link string) error
}

// BrevoClient sends email via the Brevo API.
type BrevoClient struct {
	apiKey    string
	fromEmail string
	client    *http.Client
}

// NewBrevo creates a new BrevoClient.
func NewBrevo(apiKey, fromEmail string) *BrevoClient {
	return &BrevoClient{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		client:    &http.Client{},
	}
}

type brevoSender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type brevoRecipient struct {
	Email string `json:"email"`
}

type brevoRequest struct {
	Sender      brevoSender      `json:"sender"`
	To          []brevoRecipient `json:"to"`
	Subject     string           `json:"subject"`
	HTMLContent string           `json:"htmlContent"`
}

// SendMagicLink sends a magic link email via Brevo.
func (client *BrevoClient) SendMagicLink(ctx context.Context, toEmail, link string) error {
	body := brevoRequest{
		Sender:      brevoSender{Name: "Pulse", Email: client.fromEmail},
		To:          []brevoRecipient{{Email: toEmail}},
		Subject:     "Your Pulse login link",
		HTMLContent: fmt.Sprintf(`<p>Click <a href="%s">here</a> to log in to Pulse. Link expires in 15 minutes.</p>`, link),
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.brevo.com/v3/smtp/email", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("api-key", client.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("brevo API error: status %d: %s", resp.StatusCode, bytes.TrimSpace(respBody))
	}

	return nil
}
