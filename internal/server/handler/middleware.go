package handler

import (
	"crypto/subtle"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ramadianrp/skeleton-go/internal/pkg/constants"
	appLog "github.com/ramadianrp/skeleton-go/internal/pkg/log"
)

func RequestContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := appLog.SetContextFromEchoRequest(c)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

func (h *Handler) BasicAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		username, password, ok := c.Request().BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(username), []byte(h.authUser)) != 1 ||
			subtle.ConstantTimeCompare([]byte(password), []byte(h.authPass)) != 1 {
			c.Response().Header().Set(echo.HeaderWWWAuthenticate, `Basic realm="personalities"`)
			return c.JSON(http.StatusUnauthorized, constants.DefaultResponse{
				Status: constants.StatusUnauthorized, Message: constants.MessageUnauthorized,
				Data: struct{}{}, Errors: []string{},
			})
		}
		return next(c)
	}
}
