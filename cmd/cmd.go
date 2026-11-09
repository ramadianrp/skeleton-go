package cmd

import (
	"fmt"
	"os"

	"github.com/ramadianrp/skeleton-go/internal/config"
	"github.com/ramadianrp/skeleton-go/internal/infrastructure/container"
	"github.com/ramadianrp/skeleton-go/internal/server"
)

func Run() error {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = ".env"
	}
	if err := config.Load(os.Getenv("APP_ENV"), configFile); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	dependencies, err := container.New()
	if err != nil {
		return fmt.Errorf("initialize dependencies: %w", err)
	}
	defer dependencies.Close()
	return server.StartService(dependencies)
}
