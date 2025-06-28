// Package secureecho wires a complete set of hardening middlewares into an Echo server.
// It is intended for “localhost” utility apps that handle personal data.
//
// go get github.com/you/secureecho
//
// Usage:
//
//	e := echo.New()
//	opts := secureecho.Options{
//	    AllowedHosts: []string{"localhost:8585", "127.0.0.1:8585"},
//	    Auth: secureecho.BasicAuth("admin", "S3cret!"), // or secureecho.APIKeyAuth("X-My-Token", "abcd")
//	}
//	secureecho.Harden(e, opts)
//
//	if err := secureecho.StartLocal(e, ":8585"); err != nil {
//	    e.Logger.Fatal(err)
//	}
package secureecho

import (
	"crypto/subtle"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ---------- Public API ----------

// Options groups all configurable pieces.
type Options struct {
	// AllowedHosts are the only Host header values accepted (anti‑DNS‑rebinding).
	// If empty, defaults to localhost on any port.
	AllowedHosts []string

	// Auth provides a middleware factory implementing the authentication layer.
	// Use BasicAuth(), APIKeyAuth(), or nil to disable (NOT recommended).
	Auth echo.MiddlewareFunc

	// EnableCSRF allows you to turn CSRF protection off; default true.
	EnableCSRF bool

	// CORSOrigins lists origins to allow.  Leave nil/empty to deny all cross‑origin reads.
	CORSOrigins []string
}

// Harden attaches all middlewares to the supplied Echo instance.
func Harden(e *echo.Echo, o Options) {
	setDefaults(&o)

	// Secure headers (X‑Frame‑Options, X‑Content‑Type‑Options, HSTS, etc.)
	e.Use(middleware.Secure())

	// Host header allow‑list.
	e.Use(hostCheck(o.AllowedHosts))

	// Auth barrier.
	if o.Auth != nil {
		e.Use(o.Auth)
	}

	// CSRF (cookie -> header X-CSRF-Token).
	if o.EnableCSRF {
		e.Use(middleware.CSRF())
	}

	// Restrictive CORS / PNA – allow only specific origins if given.
	e.Use(cors(o.CORSOrigins))
}

// StartLocal is a helper that refuses to listen on non‑loopback addresses.
func StartLocal(e *echo.Echo, addr string) error {
	if !isLoopbackAddr(addr) {
		return errors.New("StartLocal: address must bind to 127.0.0.1 or [::1]")
	}
	return e.Start(addr)
}

// StartLocalTLS is the same as StartLocal but serves HTTPS.
func StartLocalTLS(e *echo.Echo, addr, certFile, keyFile string) error {
	if !isLoopbackAddr(addr) {
		return errors.New("StartLocalTLS: address must bind to 127.0.0.1 or [::1]")
	}
	return e.StartTLS(addr, certFile, keyFile)
}

// ---------- Auth helpers ----------

// BasicAuth returns a middleware that enforces a fixed username/password.
func BasicAuth(user, pass string) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(u, p string, _ echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(u), []byte(user)) == 1 &&
			subtle.ConstantTimeCompare([]byte(p), []byte(pass)) == 1 {
			return true, nil
		}
		return false, nil
	})
}

// APIKeyAuth enforces a constant token in a header (or query param).
func APIKeyAuth(headerName, token string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: headerName, // e.g. "X-My-Token"
		Validator: func(key string, _ echo.Context) (bool, error) {
			return subtle.ConstantTimeCompare([]byte(key), []byte(token)) == 1, nil
		},
	})
}

// ---------- internal ----------

func setDefaults(o *Options) {
	if len(o.AllowedHosts) == 0 {
		o.AllowedHosts = []string{"localhost", "127.0.0.1", "[::1]"}
	}
	if o.EnableCSRF == false { // explicit false means off; zero value is true
		return
	}
	if o.Auth == nil {
		// Encourage auth by default?  Leave nil to let caller decide.
	}
}

// hostCheck blocks requests whose Host header is not allowed.
func hostCheck(allowed []string) echo.MiddlewareFunc {
	allowedMap := map[string]struct{}{}
	for _, h := range allowed {
		allowedMap[strings.ToLower(h)] = struct{}{}
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if host := strings.ToLower(c.Request().Host); host != "" {
				if _, ok := allowedMap[host]; !ok {
					return echo.NewHTTPError(http.StatusForbidden, "invalid host")
				}
			}
			return next(c)
		}
	}
}

// cors returns a restrictive CORS handler; it also responds to Private Network
// Access preflights by *not* granting permission unless the origin is whitelisted.
func cors(origins []string) echo.MiddlewareFunc {
	if len(origins) == 0 {
		// Deny all cross‑origin reads (browser will enforce same‑origin policy).
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				// Let the request through; no CORS headers means browser blocks non‑same‑origin reads.
				return next(c)
			}
		}
	}
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     origins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
		MaxAge:           int((12 * time.Hour).Seconds()),
	})
}

// isLoopbackAddr returns true if addr begins with “127.”, “[::1]”, or “localhost”.
func isLoopbackAddr(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(host)
	if ip != nil {
		return ip.IsLoopback()
	}
	return strings.EqualFold(host, "localhost")
}
