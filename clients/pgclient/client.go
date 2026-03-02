package pgclient

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const applicationName = "application_name"

// Client wraps a PostgreSQL connection pool with additional functionality.
type Client struct {
	pool *pgxpool.Pool
}

// NewClient creates a new PostgreSQL client with the provided configuration.
func NewClient(cfg *Config) *Client {
	config, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		panic(fmt.Errorf("failed to parse config: %w", err))
	}

	// Apply connection pool settings
	if cfg.MaxConnections > 0 {
		config.MaxConns = cfg.MaxConnections
	}
	if cfg.MinConnections > 0 {
		config.MinConns = cfg.MinConnections
	}
	if cfg.MaxConnLifetime > 0 {
		config.MaxConnLifetime = cfg.MaxConnLifetime
	}
	if cfg.MaxConnIdleTime > 0 {
		config.MaxConnIdleTime = cfg.MaxConnIdleTime
	}

	// Set application name
	config.ConnConfig.RuntimeParams[applicationName] = cfg.Application

	// Connect with retry logic
	pool, err := connectWithRetry(config)
	if err != nil {
		panic(fmt.Errorf("failed to connect to PostgreSQL: %w", err))
	}

	return &Client{
		pool: pool,
	}
}

// Pool returns the underlying pgxpool.Pool for direct access.
func (c *Client) Pool() *pgxpool.Pool {
	return c.pool
}
