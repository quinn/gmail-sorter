package web

import (
	"embed"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/quinn/gmail-sorter/internal/web/handlers"
	"github.com/quinn/gmail-sorter/pkg/core"
	"go.quinn.io/ccf/assets"
)

//go:embed public
var assetsFS embed.FS

func NewServer(spec *core.Spec) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())

	// Attach the fingerprinted assets.
	assets.Attach(
		e,
		"public",              // URL prefix -> /public
		"internal/web/public", // local directory path
		assetsFS,              // embedded FS
		os.Getenv("USE_EMBEDDED_ASSETS") == "true",
	)

	h := handlers.NewHandler(spec)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		slog.Error("error", "err", err)
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	e.GET("/healthz", h.HealthCheck)
	e.GET("/api/emails", h.EmailsHandler)
	e.GET("/oauth/start", h.OauthStartHandler)
	e.GET("/oauth/callback", h.OauthCallbackHandler)

	return e
}
