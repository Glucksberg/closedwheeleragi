package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"ClosedWheeler/pkg/config"
)

// OpenAI Codex CLI OAuth constants.
const (
	OpenAIOAuthClientID     = "app_EMoamEEZ73f0CkXaXp7hrann"
	OpenAIOAuthAuthorizeURL = "https://auth.openai.com/oauth/authorize"
	OpenAIOAuthTokenURL     = "https://auth.openai.com/oauth/token"
	OpenAIOAuthRedirectURI  = "http://localhost:1455/auth/callback"
	OpenAIOAuthScopes       = "openid profile email offline_access"
	OpenAIOAuthAudience     = "https://api.openai.com/v1"
	OpenAIOAuthCallbackPort = 1455
)

// BuildOpenAIAuthURL constructs the OpenAI OAuth authorization URL using PKCE.
func BuildOpenAIAuthURL(challenge, state string) string {
	q := url.Values{}
	q.Set("client_id", OpenAIOAuthClientID)
	q.Set("response_type", "code")
	q.Set("redirect_uri", OpenAIOAuthRedirectURI)
	q.Set("scope", OpenAIOAuthScopes)
	q.Set("audience", OpenAIOAuthAudience)
	q.Set("code_challenge", challenge)
	q.Set("code_challenge_method", "S256")
	q.Set("state", state)
	return OpenAIOAuthAuthorizeURL + "?" + q.Encode()
}

// OpenAICallbackResult is the result from the localhost callback server.
type OpenAICallbackResult struct {
	Code  string
	State string
	Err   error
}

// StartOpenAICallbackServer starts a temporary HTTP server on localhost:1455
// to capture the OAuth callback. It returns a channel that receives the result.
// The server auto-shuts down after receiving the callback or when ctx is cancelled.
func StartOpenAICallbackServer(ctx context.Context) (<-chan OpenAICallbackResult, error) {
	resultCh := make(chan OpenAICallbackResult, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		errParam := r.URL.Query().Get("error")

		if errParam != "" {
			errDesc := r.URL.Query().Get("error_description")
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "<html><body><h2>Login failed</h2><p>%s: %s</p><p>You can close this tab.</p></body></html>", errParam, errDesc)
			resultCh <- OpenAICallbackResult{Err: fmt.Errorf("%s: %s", errParam, errDesc)}
			return
		}

		if code == "" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, "<html><body><h2>Missing code</h2><p>No authorization code received.</p></body></html>")
			resultCh <- OpenAICallbackResult{Err: fmt.Errorf("no authorization code in callback")}
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "<html><body><h2>Login successful!</h2><p>You can close this tab and return to the terminal.</p></body></html>")
		resultCh <- OpenAICallbackResult{Code: code, State: state}
	})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", OpenAIOAuthCallbackPort))
	if err != nil {
		return nil, fmt.Errorf("failed to start callback server on port %d: %w", OpenAIOAuthCallbackPort, err)
	}

	server := &http.Server{Handler: mux}

	// Run server in background
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			resultCh <- OpenAICallbackResult{Err: fmt.Errorf("callback server error: %w", err)}
		}
	}()

	// Auto-shutdown on context cancel or after result received
	go func() {
		select {
		case <-ctx.Done():
		case <-resultCh:
			// Give the HTTP response time to flush
			time.Sleep(500 * time.Millisecond)
		}
		server.Close()
	}()

	return resultCh, nil
}

// ExchangeOpenAICode exchanges an authorization code for OpenAI OAuth tokens.
func ExchangeOpenAICode(code, verifier string) (*config.OAuthCredentials, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", OpenAIOAuthClientID)
	form.Set("code", code)
	form.Set("redirect_uri", OpenAIOAuthRedirectURI)
	form.Set("code_verifier", verifier)

	resp, err := oauthHTTPClient.Post(OpenAIOAuthTokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
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

	expiresAt := time.Now().UnixMilli() + tokenResp.ExpiresIn*1000
	return &config.OAuthCredentials{
		Provider:     "openai",
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// RefreshOpenAIToken refreshes an expired OpenAI OAuth token.
func RefreshOpenAIToken(refreshToken string) (*config.OAuthCredentials, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", OpenAIOAuthClientID)
	form.Set("refresh_token", refreshToken)

	resp, err := oauthHTTPClient.Post(OpenAIOAuthTokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
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

	expiresAt := time.Now().UnixMilli() + tokenResp.ExpiresIn*1000
	return &config.OAuthCredentials{
		Provider:     "openai",
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}
