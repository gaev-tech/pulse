package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Sender interface {
	SendMagicLink(ctx context.Context, toEmail, link string) error
}

type ResendClient struct {
	apiKey    string
	fromEmail string
	client    *http.Client
}

func NewResend(apiKey, fromEmail string) *ResendClient {
	return &ResendClient{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		client:    &http.Client{},
	}
}

type resendRequest struct {
	From    string `json:"from"`
	To      []string `json:"to"`
	Subject string `json:"subject"`
	HTML    string `json:"html"`
}

func (resend *ResendClient) SendMagicLink(ctx context.Context, toEmail, link string) error {
	body := resendRequest{
		From:    resend.fromEmail,
		To:      []string{toEmail},
		Subject: "Your Pulse login link",
		HTML:    fmt.Sprintf(`<p>Click <a href="%s">here</a> to log in to Pulse. Link expires in 15 minutes.</p>`, link),
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+resend.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := resend.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("resend API error: status %d", resp.StatusCode)
	}

	return nil
}
