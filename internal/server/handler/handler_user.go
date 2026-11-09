package handler

import (
	"errors"
	"net/http"
	"net/mail"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ramadianrp/skeleton-go/internal/domain/repositories"
	"github.com/ramadianrp/skeleton-go/internal/pkg/constants"
	"github.com/ramadianrp/skeleton-go/internal/usecase/user"
)

type UserHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Create(c echo.Context) error {
	var request user.CreateRequest
	if err := c.Bind(&request); err != nil {
		return writeError(c, http.StatusBadRequest, err)
	}
	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(request.Email)
	if len([]rune(request.Name)) < 3 || len([]rune(request.Name)) > 32 || !validEmail(request.Email) {
		return c.JSON(http.StatusBadRequest, constants.DefaultResponse{
			Status: constants.StatusBadRequest, Message: constants.MessageBadRequest,
			Data: struct{}{}, Errors: []string{"name must be 3-32 characters and email must be valid"},
		})
	}

	response, err := h.service.Create(c.Request().Context(), request)
	if errors.Is(err, repositories.ErrConflict) {
		return writeError(c, http.StatusConflict, err)
	}
	if err != nil {
		return writeError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, response)
}

func (h *UserHandler) FindAll(c echo.Context) error {
	response, err := h.service.FindAll(c.Request().Context(), user.FindAllRequest{})
	if err != nil {
		return writeError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, response)
}

func validEmail(value string) bool {
	if len(value) == 0 || len(value) > 32 {
		return false
	}
	address, err := mail.ParseAddress(value)
	return err == nil && address.Address == value
}

func writeError(c echo.Context, status int, err error) error {
	code := constants.StatusFailed
	message := constants.MessageFailed
	if status == http.StatusBadRequest {
		code, message = constants.StatusBadRequest, constants.MessageBadRequest
	} else if status == http.StatusConflict {
		code, message = constants.StatusConflict, constants.MessageConflict
	} else if status == http.StatusUnauthorized {
		code, message = constants.StatusUnauthorized, constants.MessageUnauthorized
	}
	return c.JSON(status, constants.DefaultResponse{
		Status: code, Message: message, Data: struct{}{}, Errors: []string{err.Error()},
	})
}
