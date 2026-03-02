package pgclient

import "time"

// Config holds the configuration for the PostgreSQL client.
type Config struct {
	// DSN is the connection string for PostgreSQL.
	DSN string

	// Application name for the PostgreSQL client.
	Application string

	// Connection pool settings
	MaxConnections  int32
	MinConnections  int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration

	// SLO defines service level objectives for monitoring
	SLO SLO
}

// SLO defines service level objectives for query performance monitoring.
type SLO struct {
	// Query is the threshold for slow query warnings
	Query time.Duration
}
