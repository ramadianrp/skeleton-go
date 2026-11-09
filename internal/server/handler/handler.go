package handler

import (
	"github.com/ramadianrp/skeleton-go/internal/usecase/healthcheck"
	"github.com/ramadianrp/skeleton-go/internal/usecase/user"
)

type Handler struct {
	HealthCheck *HealthCheckHandler
	User        *UserHandler
	authUser    string
	authPass    string
}

func New(healthService healthcheck.Service, userService user.Service, authUser, authPass string) *Handler {
	h := &Handler{
		HealthCheck: NewHealthCheckHandler(healthService),
		authUser:    authUser,
		authPass:    authPass,
	}
	if userService != nil {
		h.User = NewUserHandler(userService)
	}
	return h
}
