package web

import (
	"embed"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/quinn/gmail-sorter/internal/web/handlers"
	"go.quinn.io/ccf/assets"
)

//go:embed public
var assetsFS embed.FS

func NewServer() *echo.Echo {
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

	h := handlers.New()

	e.GET("/healthz", h.HealthCheck)

	return e
}
