package echostat

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupStats(router EchoRouter, statsFunc StatsFunc, opts StatOptions) {
	router.GET("/stats", func(c echo.Context) error {
		if opts.APIKey == "" {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "bad configuration"})
		}

		apiKeyHeader := opts.APIKeyHeader
		if apiKeyHeader == "" {
			apiKeyHeader = "X-API-Key" //nolint:gosec
		}

		apiKey := c.Request().Header.Get(apiKeyHeader)
		if apiKey != opts.APIKey {
			return c.JSON(http.StatusUnauthorized, map[string]any{"error": "not authorized"})
		}

		stats, err := statsFunc(c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, stats)
	})
}
