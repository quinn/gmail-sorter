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
	"google.golang.org/api/gmail/v1"
)

var (
	// Scopes for Gmail API
	oauthScopes = []string{
		gmail.GmailReadonlyScope,
		gmail.GmailModifyScope,
		gmail.GmailSettingsBasicScope,
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

// OauthStart redirects user to the provider's OAuth2 consent screen
func OauthStart(c echo.Context) error {
	provider := c.Param("provider")
	config, err := LoadOauthConfig(OauthProvider(provider))
	if err != nil {
		return err
	}
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusFound, url)
}

// OauthCallback handles the OAuth2 callback and saves token.json per provider
func OauthCallback(c echo.Context) error {
	provider := c.Param("provider")
	config, err := LoadOauthConfig(OauthProvider(provider))
	if err != nil {
		return err
	}
	code := c.QueryParam("code")
	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing code param")
	}
	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		return fmt.Errorf("token exchange failed: %w", err)
	}
	tokenFile := provider + "_token.json"
	f, err := os.Create(tokenFile)
	if err != nil {
		return fmt.Errorf("could not create %s: %w", tokenFile, err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(tok); err != nil {
		return fmt.Errorf("could not encode token: %w", err)
	}
	return c.Redirect(http.StatusFound, "/")
}

func Accounts(c echo.Context) error {

}
