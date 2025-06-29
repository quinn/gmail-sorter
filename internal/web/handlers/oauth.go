package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"

	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
	"github.com/quinn/gmail-sorter/pkg/db"
)

// OauthStart redirects user to the provider's OAuth2 consent screen
func OauthStart(c echo.Context) error {
	provider := c.Param("provider")
	config, err := models.LoadOauthConfig(models.OauthProvider(provider))
	if err != nil {
		return err
	}
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusFound, url)
}

// OauthCallback handles the OAuth2 callback and saves token in SQLite via GORM
func OauthCallback(c echo.Context) error {
	provider := c.Param("provider")
	config, err := models.LoadOauthConfig(models.OauthProvider(provider))
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

	email, err := getEmailFromToken(config, tok)
	if err != nil {
		return fmt.Errorf("failed to get email: %w", err)
	}

	dbConn := middleware.GetDB(c)
	if dbConn == nil {
		return fmt.Errorf("failed to get db from context")
	}
	tokenJSON, err := json.Marshal(tok)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}
	acct := db.OAuthAccount{
		Provider:  provider,
		Email:     email,
		TokenJSON: string(tokenJSON),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	// Upsert (replace if exists)
	if err := dbConn.UpsertOAuthAccount(&acct); err != nil {
		return fmt.Errorf("db error: %w", err)
	}
	return c.Redirect(http.StatusFound, "/accounts")
}

// getEmailFromToken fetches the user's email using the token and config
func getEmailFromToken(config *oauth2.Config, token *oauth2.Token) (string, error) {
	client := config.Client(context.Background(), token)
	service, err := gmail.New(client)
	if err != nil {
		return "", err
	}
	profile, err := service.Users.GetProfile("me").Do()
	if err != nil {
		return "", err
	}
	return profile.EmailAddress, nil
}

// Accounts lists all OAuth accounts
func Accounts(c echo.Context) error {
	dbConn := middleware.GetDB(c)
	if dbConn == nil {
		return c.String(http.StatusInternalServerError, "failed to get db from context")
	}
	var accounts []db.OAuthAccount
	err := dbConn.Find(&accounts).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to list accounts")
	}
	return pages.Accounts(accounts).Render(c.Request().Context(), c.Response().Writer)
}

// CreateAccount handles POST /accounts
func CreateAccount(c echo.Context) error {
	dbConn := middleware.GetDB(c)
	if dbConn == nil {
		return fmt.Errorf("failed to get db from context")
	}
	var acct db.OAuthAccount
	if err := c.Bind(&acct); err != nil {
		return err
	}
	acct.CreatedAt = time.Now().Unix()
	acct.UpdatedAt = time.Now().Unix()
	if err := dbConn.Create(&acct).Error; err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}
	return c.Redirect(http.StatusSeeOther, "/accounts")
}

// GetAccount handles GET /accounts/:id
func GetAccount(c echo.Context) error {
	dbConn := middleware.GetDB(c)
	if dbConn == nil {
		return fmt.Errorf("failed to get db from context")
	}
	id := c.Param("id")
	acct, err := dbConn.GetOAuthAccountByID(id)
	if err != nil {
		return err
	}
	return pages.AccountForm(acct).Render(c.Request().Context(), c.Response().Writer)
}

// UpdateAccount handles PUT /accounts/:id
func UpdateAccount(c echo.Context) error {
	dbConn := middleware.GetDB(c)
	if dbConn == nil {
		return fmt.Errorf("failed to get db from context")
	}
	id := c.Param("id")
	acct, err := dbConn.GetOAuthAccountByID(id)
	if err != nil {
		return err
	}
	if err := c.Bind(acct); err != nil {
		return err
	}
	acct.UpdatedAt = time.Now().Unix()
	if err := dbConn.UpdateOAuthAccount(acct); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}
	return c.Redirect(http.StatusSeeOther, "/accounts")
}

// DeleteAccount handles DELETE /accounts/:id
func DeleteAccount(c echo.Context) error {
	dbConn := middleware.GetDB(c)
	if dbConn == nil {
		return fmt.Errorf("failed to get db from context")
	}
	id := c.Param("id")
	if err := dbConn.DeleteOAuthAccount(id); err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}
	return c.Redirect(http.StatusSeeOther, "/accounts")
}
