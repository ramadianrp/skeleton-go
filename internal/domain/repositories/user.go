package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/ramadianrp/skeleton-go/internal/domain/entities"
)

var ErrConflict = errors.New("user already exists")

type User interface {
	Create(ctx context.Context, data entities.User) (entities.User, error)
	FindAll(ctx context.Context) ([]entities.User, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) (*UserRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("postgresql database is nil")
	}
	return &UserRepository{db: db}, nil
}

func (r *UserRepository) Create(ctx context.Context, user entities.User) (entities.User, error) {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (name, email) VALUES ($1, $2)
		 RETURNING id, name, email, created_at, updated_at, deleted_at`,
		user.Name, user.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		var postgresErr *pq.Error
		if errors.As(err, &postgresErr) && postgresErr.Code == "23505" {
			return entities.User{}, ErrConflict
		}
		return entities.User{}, fmt.Errorf("insert user: %w", err)
	}
	return user, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]entities.User, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, email, created_at, updated_at, deleted_at
		 FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	users := make([]entities.User, 0)
	for rows.Next() {
		var user entities.User
		var deletedAt sql.NullTime
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt, &deletedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		if deletedAt.Valid {
			user.DeletedAt = &deletedAt.Time
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}
	return users, nil
}

var _ User = (*UserRepository)(nil)
