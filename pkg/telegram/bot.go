// Package telegram provides a simple bridge for Telegram bot integration.
package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"ClosedWheeler/pkg/utils"
)

// Bot handles communication with the Telegram API
type Bot struct {
	token   string
	chatID  int64
	client  *http.Client
	baseURL string
}

// NewBot creates a new Telegram bot instance
func NewBot(token string, chatID int64) *Bot {
	return &Bot{
		token:   token,
		chatID:  chatID,
		client:  &http.Client{Timeout: 45 * time.Second}, // Increased for long polling
		baseURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
	}
}

// SendMessage sends a text message to the configured default chat ID
func (b *Bot) SendMessage(text string) error {
	return b.SendMessageToChat(b.chatID, text)
}

// SendMessageToChat sends a text message to a specific chat ID
func (b *Bot) SendMessageToChat(chatID int64, text string) error {
	if b.token == "" || chatID == 0 {
		return nil // Not configured or invalid ID
	}

	url := fmt.Sprintf("%s/sendMessage", b.baseURL)
	payload := map[string]any{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "Markdown",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	operation := func() error {
		resp, err := b.client.Post(url, "application/json", bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			apiErr := fmt.Errorf("telegram API error: %s", resp.Status)
			if utils.IsRetryableError(resp.StatusCode) {
				return apiErr
			}
			return apiErr
		}

		return nil
	}

	return utils.ExecuteWithRetry(operation, utils.DefaultRetryConfig())
}

// SendMessageWithButtons sends a text message with inline buttons
func (b *Bot) SendMessageWithButtons(chatID int64, text string, buttons [][]InlineButton) error {
	if b.token == "" || chatID == 0 {
		return nil
	}

	url := fmt.Sprintf("%s/sendMessage", b.baseURL)
	payload := map[string]any{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "Markdown",
		"reply_markup": map[string]any{
			"inline_keyboard": buttons,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	operation := func() error {
		resp, err := b.client.Post(url, "application/json", bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			// Read error body for debugging
			errorBody, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("telegram API error (status %d): %s", resp.StatusCode, string(errorBody))
		}
		return nil
	}

	return utils.ExecuteWithRetry(operation, utils.DefaultRetryConfig())
}

// AnswerCallbackQuery acknowledges a callback query to stop the "loading" state on Telegram
func (b *Bot) AnswerCallbackQuery(callbackQueryID string, text string) error {
	if b.token == "" {
		return nil
	}

	url := fmt.Sprintf("%s/answerCallbackQuery", b.baseURL)
	payload := map[string]any{
		"callback_query_id": callbackQueryID,
		"text":              text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := b.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// GetUpdates checks for new messages
func (b *Bot) GetUpdates(offset int64) ([]Update, error) {
	url := fmt.Sprintf("%s/getUpdates?offset=%d&timeout=30", b.baseURL, offset)
	var result struct {
		OK     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	operation := func() error {
		resp, err := b.client.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("telegram API error: %s", resp.Status)
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}
		return nil
	}

	if err := utils.ExecuteWithRetry(operation, utils.DefaultRetryConfig()); err != nil {
		return nil, err
	}

	return result.Result, nil
}

// InlineButton represents an inline keyboard button
type InlineButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

// Update represents a Telegram update
type Update struct {
	UpdateID      int64          `json:"update_id"`
	Message       *Message       `json:"message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

// Message represents a Telegram message
type Message struct {
	Text string `json:"text"`
	Chat struct {
		ID int64 `json:"id"`
	} `json:"chat"`
}

// CallbackQuery represents an incoming callback query from an inline keyboard
type CallbackQuery struct {
	ID      string   `json:"id"`
	Data    string   `json:"data"`
	Message *Message `json:"message"`
}
