package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	// Scopes for Gmail API
	oauthScopes = []string{
		"https://www.googleapis.com/auth/gmail.readonly",
	}
)

// OauthConfig loads credentials.json and returns an *oauth2.Config
func OauthConfig() (*oauth2.Config, error) {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, err
	}
	config, err := google.ConfigFromJSON(b, oauthScopes...)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// OauthStartHandler redirects user to Google's OAuth2 consent screen
func (h *Handler) OauthStartHandler(c echo.Context) error {
	config, err := OauthConfig()
	if err != nil {
		return fmt.Errorf("failed to load OAuth config: %w", err)
	}
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusFound, url)
}

// OauthCallbackHandler handles the OAuth2 callback and saves token.json
func (h *Handler) OauthCallbackHandler(c echo.Context) error {
	config, err := OauthConfig()
	if err != nil {
		return fmt.Errorf("failed to load OAuth config: %w", err)
	}
	code := c.QueryParam("code")
	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing code param")
	}
	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		return fmt.Errorf("token exchange failed: %w", err)
	}
	f, err := os.Create("token.json")
	if err != nil {
		return fmt.Errorf("could not create token.json: %w", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(tok); err != nil {
		return fmt.Errorf("could not encode token: %w", err)
	}
	return c.String(http.StatusOK, "OAuth Success! You may close this window.")
}
