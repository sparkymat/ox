package echostat_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sparkymat/ox/echostat"
	"github.com/sparkymat/ox/echostat/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSetup(t *testing.T) {
	t.Parallel()

	statsFunc := func(echo.Context) (map[string]any, error) {
		return map[string]any{
			"text": "Hello, World!",
		}, nil
	}
	failingStatsFunc := func(echo.Context) (map[string]any, error) {
		return nil, errors.New("something went wrong")
	}

	testCases := []struct {
		name      string
		router    echostat.EchoRouter
		statsFunc echostat.StatsFunc
		opts      echostat.StatOptions
	}{
		{
			name: "should fail without an api key",
			router: func() echostat.EchoRouter {
				m := &mocks.EchoRouter{}

				m.EXPECT().GET("/stats", mock.Anything, mock.Anything).
					Return(nil).
					Run(func(_ string, handlerFunc echo.HandlerFunc, _ ...echo.MiddlewareFunc) {
						ctx := &mocks.EchoContext{}

						ctx.EXPECT().JSON(http.StatusInternalServerError, mock.Anything).Return(nil)

						err := handlerFunc(ctx)
						assert.NoError(t, err)
					})

				return m
			}(),
			statsFunc: statsFunc,
		},
		{
			name: "should fail if api key does not match",
			router: func() echostat.EchoRouter {
				m := &mocks.EchoRouter{}

				m.EXPECT().GET("/stats", mock.Anything, mock.Anything).
					Return(nil).
					Run(func(_ string, handlerFunc echo.HandlerFunc, _ ...echo.MiddlewareFunc) {
						ctx := &mocks.EchoContext{}

						request := http.Request{}
						request.Header = make(http.Header)
						request.Header.Set("x-api-key", "wrong-api-token")

						ctx.EXPECT().Request().Return(&request)
						ctx.EXPECT().JSON(http.StatusUnauthorized, mock.Anything).Return(nil)

						err := handlerFunc(ctx)

						assert.NoError(t, err)
					})

				return m
			}(),
			statsFunc: statsFunc,
			opts: echostat.StatOptions{
				APIKey: "api-token",
			},
		},
		{
			name: "should fail if statusfn fails",
			router: func() echostat.EchoRouter {
				m := &mocks.EchoRouter{}

				m.EXPECT().GET("/stats", mock.Anything, mock.Anything).
					Return(nil).
					Run(func(_ string, handlerFunc echo.HandlerFunc, _ ...echo.MiddlewareFunc) {
						ctx := &mocks.EchoContext{}

						request := http.Request{}
						request.Header = make(http.Header)
						request.Header.Set("x-api-key", "api-token")

						ctx.EXPECT().Request().Return(&request)
						ctx.EXPECT().JSON(http.StatusInternalServerError, mock.Anything).Return(nil)

						err := handlerFunc(ctx)

						assert.NoError(t, err)
					})

				return m
			}(),
			statsFunc: failingStatsFunc,
			opts: echostat.StatOptions{
				APIKey: "api-token",
			},
		},
		{
			name: "should succeed if api key does match",
			router: func() echostat.EchoRouter {
				m := &mocks.EchoRouter{}

				m.EXPECT().GET("/stats", mock.Anything, mock.Anything).
					Return(nil).
					Run(func(_ string, handlerFunc echo.HandlerFunc, _ ...echo.MiddlewareFunc) {
						ctx := &mocks.EchoContext{}

						request := http.Request{}
						request.Header = make(http.Header)
						request.Header.Set("x-api-key", "api-token")

						ctx.EXPECT().Request().Return(&request)
						ctx.EXPECT().JSON(http.StatusOK, mock.Anything).Return(nil)

						err := handlerFunc(ctx)

						assert.NoError(t, err)
					})

				return m
			}(),
			statsFunc: statsFunc,
			opts: echostat.StatOptions{
				APIKey: "api-token",
			},
		},
		{
			name: "should succeed with custom api key header",
			router: func() echostat.EchoRouter {
				m := &mocks.EchoRouter{}

				m.EXPECT().GET("/stats", mock.Anything, mock.Anything).
					Return(nil).
					Run(func(_ string, handlerFunc echo.HandlerFunc, _ ...echo.MiddlewareFunc) {
						ctx := &mocks.EchoContext{}

						request := http.Request{}
						request.Header = make(http.Header)
						request.Header.Set("x-other-api-key", "api-token")

						ctx.EXPECT().Request().Return(&request)
						ctx.EXPECT().JSON(http.StatusOK, mock.Anything).Return(nil)

						err := handlerFunc(ctx)

						assert.NoError(t, err)
					})

				return m
			}(),
			statsFunc: statsFunc,
			opts: echostat.StatOptions{
				APIKey:       "api-token",
				APIKeyHeader: "x-other-api-key",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			echostat.SetupStats(tc.router, tc.statsFunc, tc.opts)
		})
	}
}
