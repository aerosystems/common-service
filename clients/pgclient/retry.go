package pgclient

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const maxPGConnectionRetries = 5

// connectWithRetry attempts to connect to PostgreSQL with exponential backoff.
func connectWithRetry(config *pgxpool.Config) (*pgxpool.Pool, error) {
	baseDelay := time.Second

	for attempt := 1; attempt <= maxPGConnectionRetries; attempt++ {
		pool, err := pgxpool.NewWithConfig(context.Background(), config)
		if err == nil {
			// Test the connection
			err = pool.Ping(context.Background())
		}
		if err == nil {
			return pool, nil
		}

		if attempt == maxPGConnectionRetries {
			return nil, fmt.Errorf("failed to connect to PostgreSQL after %d attempts: %w", attempt, err)
		}

		// Exponential backoff with jitter: 1s, 2s, 4s, 8s...
		sleep := baseDelay * time.Duration(1<<uint(attempt-1)) // nolint: gosec
		jitter := time.Duration(rand.Int63n(int64(sleep / 2))) // nolint: gosec
		sleep += jitter

		time.Sleep(sleep)
	}

	return nil, fmt.Errorf("failed to connect after %d retries", maxPGConnectionRetries)
}
