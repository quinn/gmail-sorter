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

func NewServer(spec *core.Spec) (*echo.Echo, error) {
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

	h, err := handlers.NewHandler(spec)
	if err != nil {
		return nil, err
	}

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		slog.Error("error", "err", err)
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	e.GET("/healthz", h.HealthCheck)
	e.GET("/emails", h.EmailsHandler)
	e.GET("/emails/:id", h.EmailHandler)
	e.GET("/emails/:id/group", h.GroupEmailHandler)
	e.POST("/emails/:id/skip", h.SkipEmailHandler)
	e.POST("/emails/:id/delete", h.DeleteEmailHandler)
	e.GET("/oauth/start", h.OauthStartHandler)
	e.GET("/oauth/callback", h.OauthCallbackHandler)
	e.GET("/", h.IndexHandler)

	return e, nil
}
