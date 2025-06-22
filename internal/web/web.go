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

	e.GET("/healthz", h.Health)
	e.GET("/emails", h.Emails)
	e.GET("/emails/:id", h.Email)
	e.GET("/emails/:id/group", h.GroupEmail)
	e.GET("/emails/group-by/:type", h.GroupByEmail)
	e.POST("/emails/:id/skip", h.SkipEmail)
	e.POST("/emails/:id/delete", h.DeleteEmail)
	e.GET("/oauth/start", h.OauthStart)
	e.GET("/oauth/callback", h.OauthCallback)
	e.GET("/", h.Index)

	return e, nil
}
