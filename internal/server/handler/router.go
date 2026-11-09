package handler

import "github.com/labstack/echo/v4"

func (h *Handler) RegisterRoutes(app *echo.Echo) {
	app.GET("/", h.HealthCheck.Check)
	app.GET("/health", h.HealthCheck.Check)
	if h.User == nil {
		return
	}
	users := app.Group("/v1/users", h.BasicAuth)
	users.GET("", h.User.FindAll)
	users.POST("", h.User.Create)
}
