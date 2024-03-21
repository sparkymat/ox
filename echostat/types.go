package echostat

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=EchoRouter --case=underscore --with-expecter
//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=Context --srcpkg=github.com/labstack/echo/v4 --structname=EchoContext --case=underscore --with-expecter

import "github.com/labstack/echo/v4"

type StatOptions struct {
	APIKey       string
	APIKeyHeader string
}

type EchoRouter interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

type StatsFunc func(c echo.Context) (map[string]any, error)
