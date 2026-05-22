package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ResendMailer struct {
	apiKey    string
	fromEmail string
	client    *http.Client
}

func NewResendMailer(apiKey, fromEmail string) *ResendMailer {
	return &ResendMailer{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		client:    &http.Client{},
	}
}

func (m *ResendMailer) Send(ctx context.Context, email Email) error {
	payload := map[string]string{
		"from":    m.fromEmail,
		"to":      email.To,
		"subject": email.Subject,
		"text":    email.Body,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("resend marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.resend.com/emails", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("resend request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resend: unexpected status %d", resp.StatusCode)
	}

	log.Printf("mailer: sent '%s' to %s", email.Subject, email.To)
	return nil
}
