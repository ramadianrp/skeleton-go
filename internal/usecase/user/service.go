package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/ramadianrp/skeleton-go/internal/domain/entities"
	"github.com/ramadianrp/skeleton-go/internal/domain/repositories"
	"github.com/ramadianrp/skeleton-go/internal/pkg/constants"
)

type Service interface {
	Create(ctx context.Context, req CreateRequest) (CreateResponse, error)
	FindAll(ctx context.Context, req FindAllRequest) (FindAllResponse, error)
}

type service struct {
	users repositories.User
}

func NewService(users repositories.User) Service {
	return &service{users: users}
}

func (s *service) Create(ctx context.Context, req CreateRequest) (CreateResponse, error) {
	created, err := s.users.Create(ctx, entities.User{
		Name:  strings.TrimSpace(req.Name),
		Email: strings.TrimSpace(req.Email),
	})
	if err != nil {
		return CreateResponse{}, fmt.Errorf("create user: %w", err)
	}
	return constants.DefaultResponse{
		Status: constants.StatusSuccess, Message: constants.MessageSuccess,
		Data: created, Errors: []string{},
	}, nil
}

func (s *service) FindAll(ctx context.Context, _ FindAllRequest) (FindAllResponse, error) {
	users, err := s.users.FindAll(ctx)
	if err != nil {
		return FindAllResponse{}, fmt.Errorf("find users: %w", err)
	}
	return constants.DefaultResponse{
		Status: constants.StatusSuccess, Message: constants.MessageSuccess,
		Data: users, Errors: []string{},
	}, nil
}
