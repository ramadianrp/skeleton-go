package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ramadianrp/skeleton-go/internal/infrastructure/container"
	appLog "github.com/ramadianrp/skeleton-go/internal/pkg/log"
	"github.com/ramadianrp/skeleton-go/internal/server/handler"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func New(handlers *handler.Handler) *echo.Echo {
	app := echo.New()
	app.HideBanner = true
	SetupMiddleware(app)
	handlers.RegisterRoutes(app)
	return app
}

func StartH2CServer(dependencies *container.Container) error {
	app := New(dependencies.Handler)
	address := net.JoinHostPort(dependencies.Config.Apps.Address, dependencies.Config.Apps.HttpPort)
	server := &http.Server{
		Addr:    address,
		Handler: h2c.NewHandler(app, &http2.Server{MaxConcurrentStreams: 500}),
	}
	appLog.Info(context.Background(), "h2c server started", "address", address)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("start h2c server: %w", err)
	}
	return nil
}
