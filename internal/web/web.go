package web

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/quinn/gmail-sorter/internal/web/handlers"
	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/pkg/core"
	"go.quinn.io/ccf/assets"
	"google.golang.org/api/gmail/v1"
)

//go:embed public
var assetsFS embed.FS

var messages []*gmail.Message

func NewServer(spec *core.Spec) (*echo.Echo, error) {
	e := echo.New()
	e.Use(echomiddleware.Logger())

	// Attach the fingerprinted assets.
	assets.Attach(
		e,
		"public",              // URL prefix -> /public
		"internal/web/public", // local directory path
		assetsFS,              // embedded FS
		os.Getenv("USE_EMBEDDED_ASSETS") == "true",
	)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		slog.Error("error", "err", err)
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Fetch the list of messages for navigation
	listRes, err := spec.GmailService().Users.Messages.List("me").MaxResults(50).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	messages = listRes.Messages

	e.Use(middleware.Messages(&messages))
	e.Use(middleware.Gmail(spec.GmailService()))

	e.GET("/healthz", handlers.Health)
	e.GET("/oauth/start", handlers.OauthStart)
	e.GET("/oauth/callback", handlers.OauthCallback)

	for _, action := range models.Actions {
		switch action.Method {
		case "GET":
			e.GET(action.Path, action.Handler())
		case "POST":
			e.POST(action.Path, action.Handler())
		default:
			slog.Error("unknown action method", "method", action.Method)
		}
	}

	return e, nil
}
