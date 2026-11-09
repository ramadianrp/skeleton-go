package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/labstack/gommon/color"
	_ "github.com/lib/pq"
	"github.com/ramadianrp/skeleton-go/internal/config"
)

func NewDB(cfg config.PostgreSQLDB) (*sql.DB, error) {
	if cfg.Host == "" || cfg.Port <= 0 || cfg.Name == "" || cfg.Username == "" {
		return nil, errors.New("postgresql host, port, database name, and username are required")
	}

	connectionURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.Username, cfg.Password),
		Host:   net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Path:   cfg.Name,
	}
	query := url.Values{}
	if cfg.SSLMode != "" {
		query.Set("sslmode", cfg.SSLMode)
	}
	if cfg.Schema != "" {
		query.Set("search_path", cfg.Schema)
	}
	connectionURL.RawQuery = query.Encode()

	db, err := sql.Open("postgres", connectionURL.String())
	if err != nil {
		return nil, fmt.Errorf("open postgresql connection: %w", err)
	}
	if cfg.MaxOpenConnections > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConnections)
	}
	if cfg.MinIdleConnections > 0 {
		db.SetMaxIdleConns(cfg.MinIdleConnections)
	}
	if cfg.MaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("connect to postgresql: %w", err)
	}

	color.Println(color.Green(fmt.Sprintf("⇨ connected to postgres db on %s\n", cfg.Name)))

	return db, nil
}
