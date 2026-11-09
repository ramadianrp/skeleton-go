package container

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ramadianrp/skeleton-go/internal/config"
	"github.com/ramadianrp/skeleton-go/internal/domain/repositories"
	"github.com/ramadianrp/skeleton-go/internal/infrastructure/postgresql"
	appLog "github.com/ramadianrp/skeleton-go/internal/pkg/log"
	"github.com/ramadianrp/skeleton-go/internal/server/handler"
	"github.com/ramadianrp/skeleton-go/internal/usecase/healthcheck"
	"github.com/ramadianrp/skeleton-go/internal/usecase/user"
)

type Container struct {
	Config             *config.DefaultConfig
	PostgresqlDB       *config.PostgreSQLDB
	DB                 *sql.DB
	Logger             *slog.Logger
	UserRepository     repositories.User
	HealthCheckService healthcheck.Service
	UserService        user.Service
	Handler            *handler.Handler
}

func New() (*Container, error) {
	appLog.New()
	dependencies := &Container{
		Config:             newAppConfig(),
		PostgresqlDB:       newPostgresqlConfig(),
		Logger:             appLog.SystemLogger(),
		HealthCheckService: healthcheck.NewService(),
	}

	// Database
	database, err := connectPostgresql(dependencies)
	if err != nil {
		return nil, err
	}
	dependencies.DB = database

	// Repositories
	dependencies.UserRepository, err = newUserRepository(dependencies)
	if err != nil {
		_ = dependencies.Close()
		return nil, err
	}

	// Services
	dependencies.UserService = newUserService(dependencies.UserRepository)

	// HTTP handlers
	dependencies.Handler = newHandler(dependencies)

	if err := dependencies.Validate(); err != nil {
		_ = dependencies.Close()
		return nil, err
	}
	return dependencies, nil
}

func (c *Container) Validate() error {
	if c == nil {
		return errors.New("container is nil")
	}
	if c.Config == nil {
		return errors.New("application config is nil")
	}
	if c.PostgresqlDB == nil {
		return errors.New("PostgreSQL config is nil")
	}
	if c.Logger == nil {
		return errors.New("logger is nil")
	}
	if c.HealthCheckService == nil {
		return errors.New("health check service is nil")
	}
	if c.Handler == nil {
		return errors.New("HTTP handler is nil")
	}
	if c.PostgresqlDB.Host != "" && (c.DB == nil || c.UserService == nil) {
		return errors.New("PostgreSQL dependencies are incomplete")
	}
	return nil
}

func newAppConfig() *config.DefaultConfig {
	return &config.DefaultConfig{
		Apps: config.Apps{
			Name:     config.GetString("apps.name"),
			Version:  config.GetString("apps.version"),
			Address:  config.GetString("apps.address"),
			HttpPort: config.GetString("apps.httpPort"),
		},
	}
}

func newPostgresqlConfig() *config.PostgreSQLDB {
	return &config.PostgreSQLDB{
		Username:           config.GetString("postgresql.skeleton.username"),
		Password:           config.GetString("postgresql.skeleton.password"),
		Name:               config.GetString("postgresql.skeleton.db"),
		Schema:             config.GetString("postgresql.skeleton.schema"),
		Host:               config.GetString("postgresql.skeleton.host"),
		Port:               config.GetInt("postgresql.skeleton.port"),
		SSLMode:            config.GetString("postgresql.skeleton.sslmode"),
		MinIdleConnections: config.GetInt("postgresql.skeleton.minidleconnections"),
		MaxOpenConnections: config.GetInt("postgresql.skeleton.maxopenconnections"),
		MaxLifetime:        config.GetInt("postgresql.skeleton.maxlifetime"),
		LogMode:            config.GetBool("postgresql.skeleton.logmode"),
	}
}

func connectPostgresql(c *Container) (*sql.DB, error) {
	cfg := c.PostgresqlDB
	if cfg.Host == "" {
		c.Logger.Warn("PostgreSQL connection skipped", "reason", "host is not configured")
		return nil, nil
	}

	database, err := postgresql.NewDB(*cfg)
	if err != nil {
		c.Logger.Error("failed to connect to PostgreSQL database",
			"database", cfg.Name,
			"host", cfg.Host,
			"error", err,
		)
		return nil, err
	}
	c.Logger.Info("connected to PostgreSQL database",
		"database", cfg.Name,
		"host", cfg.Host,
	)
	return database, nil
}

func newUserRepository(c *Container) (repositories.User, error) {
	if c.DB == nil {
		return nil, nil
	}

	username := config.GetString("basicAuth.v1.username")
	password := config.GetString("basicAuth.v1.password")
	if username == "" || password == "" {
		err := errors.New("basicAuth.v1 username and password are required when PostgreSQL is configured")
		c.Logger.Error("application configuration invalid", "error", err)
		return nil, err
	}

	repository, err := repositories.NewUserRepository(c.DB)
	if err != nil {
		return nil, fmt.Errorf("create user repository: %w", err)
	}
	return repository, nil
}

func newUserService(repository repositories.User) user.Service {
	if repository == nil {
		return nil
	}
	return user.NewService(repository)
}

func newHandler(dependencies *Container) *handler.Handler {
	return handler.New(
		dependencies.HealthCheckService,
		dependencies.UserService,
		config.GetString("basicAuth.v1.username"),
		config.GetString("basicAuth.v1.password"),
	)
}

func (c *Container) Close() error {
	if c == nil || c.DB == nil {
		return nil
	}
	return c.DB.Close()
}
