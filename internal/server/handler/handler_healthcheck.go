package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ramadianrp/skeleton-go/internal/usecase/healthcheck"
)

type HealthCheckHandler struct {
	service healthcheck.Service
}

func NewHealthCheckHandler(service healthcheck.Service) *HealthCheckHandler {
	return &HealthCheckHandler{service: service}
}

func (h *HealthCheckHandler) Check(c echo.Context) error {
	response, err := h.service.HealthCheck(c.Request().Context())
	if err != nil {
		return writeError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, response)
}
