package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// OAuthCredentials holds OAuth 2.0 tokens for Anthropic.
type OAuthCredentials struct {
	Provider     string `json:"provider"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"` // unix milliseconds
}

// IsExpired returns true if the access token is expired (with 5-minute buffer).
func (c *OAuthCredentials) IsExpired() bool {
	if c == nil {
		return true
	}
	return time.Now().UnixMilli() >= c.ExpiresAt
}

// NeedsRefresh returns true if the token is expired or within 5 minutes of expiry.
func (c *OAuthCredentials) NeedsRefresh() bool {
	if c == nil {
		return true
	}
	buffer := int64(5 * 60 * 1000) // 5 minutes in ms
	return time.Now().UnixMilli() >= (c.ExpiresAt - buffer)
}

// ExpiresIn returns how long until the token expires.
func (c *OAuthCredentials) ExpiresIn() time.Duration {
	if c == nil {
		return 0
	}
	ms := c.ExpiresAt - time.Now().UnixMilli()
	if ms < 0 {
		return 0
	}
	return time.Duration(ms) * time.Millisecond
}

// oauthPath returns the path to the OAuth credentials file.
func oauthPath() string {
	return filepath.Join(".agi", "oauth.json")
}

// LoadOAuth loads OAuth credentials from .agi/oauth.json.
// Returns nil (no error) if the file doesn't exist.
func LoadOAuth() (*OAuthCredentials, error) {
	data, err := os.ReadFile(oauthPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var creds OAuthCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}
	return &creds, nil
}

// SaveOAuth saves OAuth credentials to .agi/oauth.json.
func SaveOAuth(creds *OAuthCredentials) error {
	if err := os.MkdirAll(".agi", 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(oauthPath(), data, 0600)
}
