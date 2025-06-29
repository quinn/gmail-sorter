package models

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
)

type OauthProvider string

const (
	ProviderGoogle  OauthProvider = "google"
	ProviderTodoist OauthProvider = "todoist"
)

type OauthConfigFile struct {
	Provider     string   `json:"provider"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`
	AuthURL      string   `json:"auth_url"`
	TokenURL     string   `json:"token_url"`
	RedirectURL  string   `json:"redirect_url"`
}

type OauthConfigMap map[OauthProvider]OauthConfigFile

func LoadOauthConfig(provider OauthProvider) (*oauth2.Config, error) {
	b, err := os.ReadFile("oauth_credentials.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read oauth_credentials.json: %w", err)
	}

	var configs []OauthConfigFile
	if err := json.Unmarshal(b, &configs); err != nil {
		return nil, fmt.Errorf("failed to parse oauth_credentials.json: %w", err)
	}

	for _, cfg := range configs {
		if cfg.Provider == string(provider) {
			return &oauth2.Config{
				ClientID:     cfg.ClientID,
				ClientSecret: cfg.ClientSecret,
				Scopes:       cfg.Scopes,
				Endpoint: oauth2.Endpoint{
					AuthURL:  cfg.AuthURL,
					TokenURL: cfg.TokenURL,
				},
				RedirectURL: cfg.RedirectURL,
			}, nil
		}
	}

	return nil, fmt.Errorf("provider %s not found in oauth_credentials.json", provider)
}
