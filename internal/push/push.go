// Package push delivers notifications through the Expo Push API.
// Set PUSH_DRY_RUN=true to log instead of send (local development / tests).
package push

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/shubhangcs/agromart-server/internal/env"
)

const expoPushURL = "https://exp.host/--/api/v2/push/send"

// TokenStore is the persistence the notifier needs.
type TokenStore interface {
	GetTokens(userID string) ([]string, error)
	Delete(token string) error
}

// Notifier fans a notification out to all of a user's devices, asynchronously.
type Notifier interface {
	Notify(userID, title, body string, data map[string]string)
}

type expoMessage struct {
	To       string            `json:"to"`
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	Data     map[string]string `json:"data,omitempty"`
	Sound    string            `json:"sound"`
	Priority string            `json:"priority"`
}

type expoTicket struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Details *struct {
		Error string `json:"error"`
	} `json:"details"`
}

type ExpoNotifier struct {
	store  TokenStore
	logger *slog.Logger
	client *http.Client
	dryRun bool
	token  string // optional Expo access token
}

func NewExpoNotifier(store TokenStore, logger *slog.Logger) *ExpoNotifier {
	return &ExpoNotifier{
		store:  store,
		logger: logger,
		client: &http.Client{Timeout: 10 * time.Second},
		dryRun: strings.EqualFold(env.GetString("PUSH_DRY_RUN", "false"), "true"),
		token:  env.GetString("EXPO_ACCESS_TOKEN", ""),
	}
}

// IsExpoToken reports whether a string looks like an Expo push token.
func IsExpoToken(t string) bool {
	return (strings.HasPrefix(t, "ExponentPushToken[") || strings.HasPrefix(t, "ExpoPushToken[")) && strings.HasSuffix(t, "]")
}

func (n *ExpoNotifier) Notify(userID, title, body string, data map[string]string) {
	go func() {
		tokens, err := n.store.GetTokens(userID)
		if err != nil {
			n.logger.Error("push: load tokens", "user_id", userID, "error", err)
			return
		}
		if len(tokens) == 0 {
			return
		}
		if n.dryRun {
			n.logger.Info("push (dry run)", "user_id", userID, "devices", len(tokens), "title", title, "body", body, "data", data)
			return
		}
		msgs := make([]expoMessage, 0, len(tokens))
		for _, t := range tokens {
			msgs = append(msgs, expoMessage{To: t, Title: title, Body: body, Data: data, Sound: "default", Priority: "high"})
		}
		n.send(tokens, msgs)
	}()
}

func (n *ExpoNotifier) send(tokens []string, msgs []expoMessage) {
	payload, _ := json.Marshal(msgs)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, expoPushURL, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if n.token != "" {
		req.Header.Set("Authorization", "Bearer "+n.token)
	}
	res, err := n.client.Do(req)
	if err != nil {
		n.logger.Error("push: expo request", "error", err)
		return
	}
	defer res.Body.Close()
	var out struct {
		Data []expoTicket `json:"data"`
	}
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		n.logger.Error("push: decode expo response", "status", res.StatusCode, "error", err)
		return
	}
	for i, t := range out.Data {
		if t.Status == "ok" || i >= len(tokens) {
			continue
		}
		n.logger.Warn("push: ticket error", "token", tokens[i], "message", t.Message)
		if t.Details != nil && t.Details.Error == "DeviceNotRegistered" {
			_ = n.store.Delete(tokens[i])
		}
	}
}
