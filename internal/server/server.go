package server

import (
	"errors"

	"github.com/ramadianrp/skeleton-go/internal/infrastructure/container"
	httpserver "github.com/ramadianrp/skeleton-go/internal/server/http"
)

func StartService(dependencies *container.Container) error {
	if dependencies == nil || dependencies.Config == nil || dependencies.Handler == nil {
		return errors.New("server dependencies are not initialized")
	}
	return httpserver.StartH2CServer(dependencies)
}
