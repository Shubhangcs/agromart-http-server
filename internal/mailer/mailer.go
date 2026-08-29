// Package mailer sends transactional email through Resend (https://resend.com).
// Requires RESEND_API_KEY and EMAIL_FROM (e.g. "South Canara Agro Mart <noreply@southcanaraagromart.com>").
// When either is missing a log-only mailer is used so local development never needs an API key.
package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/shubhangcs/agromart-server/internal/env"
)

const resendURL = "https://api.resend.com/emails"

type Mailer interface {
	Send(ctx context.Context, to, subject, textBody, htmlBody string) error
}

// New returns a Resend mailer when configured, otherwise a logging mailer.
func New(logger *slog.Logger) Mailer {
	apiKey := env.GetString("RESEND_API_KEY", "")
	from := env.GetString("EMAIL_FROM", "")
	if apiKey == "" || from == "" {
		logger.Warn("RESEND_API_KEY / EMAIL_FROM not set: emails will be logged, not sent")
		return &logMailer{logger: logger}
	}
	return &resendMailer{apiKey: apiKey, from: from, logger: logger, client: &http.Client{Timeout: 15 * time.Second}}
}

type resendMailer struct {
	apiKey string
	from   string
	logger *slog.Logger
	client *http.Client
}

func (m *resendMailer) Send(ctx context.Context, to, subject, textBody, htmlBody string) error {
	body, _ := json.Marshal(map[string]any{
		"from": m.from, "to": []string{to}, "subject": subject, "text": textBody, "html": htmlBody,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")
	res, err := m.client.Do(req)
	if err != nil {
		m.logger.Error("resend request failed", "to", to, "error", err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
		m.logger.Error("resend rejected email", "to", to, "status", res.StatusCode, "response", string(b))
		return fmt.Errorf("resend: status %d", res.StatusCode)
	}
	return nil
}

type logMailer struct{ logger *slog.Logger }

// Captured is the most recent email handled by the log-only mailer (for local inspection and tests).
type Captured struct{ To, Subject, Text string }

var (
	captureMu sync.Mutex
	captured  []Captured
)

// LastCaptured returns the newest captured email for the address, if any.
func LastCaptured(to string) (Captured, bool) {
	captureMu.Lock()
	defer captureMu.Unlock()
	for i := len(captured) - 1; i >= 0; i-- {
		if captured[i].To == to {
			return captured[i], true
		}
	}
	return Captured{}, false
}

func (m *logMailer) Send(_ context.Context, to, subject, textBody, _ string) error {
	m.logger.Info("email (not sent: mailer unconfigured)", "to", to, "subject", subject, "body", textBody)
	captureMu.Lock()
	captured = append(captured, Captured{To: to, Subject: subject, Text: textBody})
	if len(captured) > 50 {
		captured = captured[len(captured)-50:]
	}
	captureMu.Unlock()
	return nil
}
