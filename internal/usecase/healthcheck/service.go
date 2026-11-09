package healthcheck

import (
	"context"
	"time"

	"github.com/ramadianrp/skeleton-go/internal/config"
)

type Service interface {
	HealthCheck(ctx context.Context) (Response, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) HealthCheck(context.Context) (Response, error) {
	version := config.GetString("apps.version")
	if version == "" {
		version = "dev"
	}
	return Response{
		Message:    "Server up and running",
		ServerTime: time.Now().Format(time.RFC1123),
		Version:    version,
	}, nil
}
