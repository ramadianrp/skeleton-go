package user

import "github.com/ramadianrp/skeleton-go/internal/pkg/constants"

type CreateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type FindAllRequest struct{}

type CreateResponse = constants.DefaultResponse
type FindAllResponse = constants.DefaultResponse
