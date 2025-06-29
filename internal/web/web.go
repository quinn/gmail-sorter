package web

import (
	"embed"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/quinn/gmail-sorter/internal/web/handlers"
	"github.com/quinn/gmail-sorter/internal/web/middleware"
	"github.com/quinn/gmail-sorter/internal/web/models"
	"github.com/quinn/gmail-sorter/internal/web/views/pages"
	"github.com/quinn/gmail-sorter/internal/web/views/ui"
	"github.com/quinn/gmail-sorter/pkg/gmailapi"
	"go.quinn.io/ccf/assets"
)

//go:embed public
var assetsFS embed.FS

func NewServer(api *gmailapi.GmailAPI) (*echo.Echo, error) {
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

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		statusCode := http.StatusInternalServerError
		if httpErr, ok := err.(*echo.HTTPError); ok {
			statusCode = httpErr.Code
		}

		var renderErr error
		if c.Request().Header.Get("HX-Request") == "true" {
			c.Response().Header().Set("HX-Retarget", "#flash")
			c.Response().Header().Set("HX-Reswap", "innerHTML")
			renderErr = ui.FlashMessage("error", err.Error()).Render(c.Request().Context(), c.Response().Writer)
		} else {
			c.Response().WriteHeader(statusCode)
			renderErr = pages.ErrorPage([]models.ActionLink{
				handlers.IndexAction.Link(),
			}, err.Error()).Render(c.Request().Context(), c.Response().Writer)
		}

		if renderErr != nil {
			_ = c.JSON(statusCode, err)
		}
	}

	e.Use(middleware.Gmail(api))
	e.Use(middleware.Echo)

	e.GET("/healthz", handlers.Health)
	e.GET("/oauth/:provider/start", handlers.OauthStart)
	e.GET("/oauth/:provider/callback", handlers.OauthCallback)

	e.GET("/test", func(c echo.Context) error {
		c.Response().Header().Set("HX-Trigger", "test-event")
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	for _, action := range models.Actions {
		switch action.Method {
		case "GET":
			e.GET(action.Path, action.Handler()).Name = action.ID
		case "POST":
			e.POST(action.Path, action.Handler()).Name = action.ID
		default:
			slog.Error("unknown action method", "method", action.Method)
		}
	}

	return e, nil
}
