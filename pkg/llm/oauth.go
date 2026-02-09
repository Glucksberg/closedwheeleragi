package llm

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"ClosedWheeler/pkg/config"
)

// OAuth constants (from openclaw/pi-ai reference implementation).
const (
	OAuthClientID     = "9d1c250a-e61b-44d9-88ed-5944d1962f5e"
	OAuthAuthorizeURL = "https://claude.ai/oauth/authorize"
	OAuthTokenURL     = "https://console.anthropic.com/v1/oauth/token"
	OAuthRedirectURI  = "https://console.anthropic.com/oauth/code/callback"
	OAuthScopes       = "org:create_api_key user:profile user:inference"
)

// oauthHTTPClient has a sensible timeout so token calls don't hang forever.
var oauthHTTPClient = &http.Client{Timeout: 15 * time.Second}

// OAuthTokenResponse is the JSON response from the token endpoint.
type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
	TokenType    string `json:"token_type"`
}

// GeneratePKCE creates a code_verifier and its S256 code_challenge.
func GeneratePKCE() (verifier, challenge string, err error) {
	// 32 random bytes → base64url verifier
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	verifier = base64.RawURLEncoding.EncodeToString(buf)

	// SHA-256 of verifier → base64url challenge
	h := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(h[:])
	return verifier, challenge, nil
}

// BuildAuthURL constructs the full Anthropic OAuth authorization URL.
// The verifier is included as the state parameter for CSRF validation.
func BuildAuthURL(challenge, verifier string) string {
	params := url.Values{
		"code":                  {"true"},
		"client_id":             {OAuthClientID},
		"response_type":         {"code"},
		"redirect_uri":          {OAuthRedirectURI},
		"scope":                 {OAuthScopes},
		"code_challenge":        {challenge},
		"code_challenge_method": {"S256"},
		"state":                 {verifier},
	}
	return OAuthAuthorizeURL + "?" + params.Encode()
}

// ExchangeCode exchanges an authorization code for OAuth tokens.
// The authCode should be in "code#state" format as returned by the redirect.
// The expectedState is validated against the state portion to prevent CSRF.
func ExchangeCode(authCode, verifier, expectedState string) (*config.OAuthCredentials, error) {
	parts := strings.SplitN(authCode, "#", 2)
	code := parts[0]

	state := ""
	if len(parts) == 2 {
		state = parts[1]
	}

	// Validate state matches what we sent (CSRF protection)
	if expectedState != "" && state != expectedState {
		return nil, fmt.Errorf("OAuth state mismatch (possible CSRF attack)")
	}

	body, err := json.Marshal(map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     OAuthClientID,
		"code":          code,
		"state":         state,
		"redirect_uri":  OAuthRedirectURI,
		"code_verifier": verifier,
	})
	if err != nil {
		return nil, err
	}

	resp, err := oauthHTTPClient.Post(OAuthTokenURL, "application/json", strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("token exchange request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, truncateError(respBody))
	}

	var tokenResp OAuthTokenResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	// Store true expiry time (buffer is applied only in NeedsRefresh())
	expiresAt := time.Now().UnixMilli() + tokenResp.ExpiresIn*1000

	return &config.OAuthCredentials{
		Provider:     "anthropic",
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// RefreshOAuthToken refreshes an expired OAuth token.
func RefreshOAuthToken(refreshToken string) (*config.OAuthCredentials, error) {
	body, err := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     OAuthClientID,
		"refresh_token": refreshToken,
	})
	if err != nil {
		return nil, err
	}

	resp, err := oauthHTTPClient.Post(OAuthTokenURL, "application/json", strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("token refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read refresh response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh failed (status %d): %s", resp.StatusCode, truncateError(respBody))
	}

	var tokenResp OAuthTokenResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse refresh response: %w", err)
	}

	// Store true expiry time (buffer is applied only in NeedsRefresh())
	expiresAt := time.Now().UnixMilli() + tokenResp.ExpiresIn*1000

	return &config.OAuthCredentials{
		Provider:     "anthropic",
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// truncateError limits error body size to avoid leaking sensitive details.
func truncateError(body []byte) string {
	s := string(body)
	if len(s) > 200 {
		return s[:200] + "..."
	}
	return s
}
