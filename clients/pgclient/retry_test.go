package pgclient

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectWithRetry_Success_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pgContainer := setupPostgreSQLContainer(t)
	defer pgContainer.Cleanup(t)

	cfg := &Config{
		DSN:         pgContainer.DSN,
		Application: "test-app",
	}

	// This should succeed on first attempt
	client := NewClient(cfg)
	require.NotNil(t, client)
	defer client.Close(context.TODO())

	ctx := context.Background()
	err := client.HealthCheck(ctx)
	assert.NoError(t, err)
}

func TestConnectWithRetry_InvalidDSN(t *testing.T) {
	cfg := &Config{
		DSN:         "postgres://invalid:invalid@nonexistent:5432/testdb",
		Application: "test-app",
	}

	// Should panic after exhausting retries
	assert.Panics(t, func() {
		NewClient(cfg)
	}, "should panic after max retries with invalid DSN")
}

func TestNewClient_WithRetry_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pgContainer := setupPostgreSQLContainer(t)
	defer pgContainer.Cleanup(t)

	// Test that client creation succeeds with valid connection
	start := time.Now()

	cfg := &Config{
		DSN:         pgContainer.DSN,
		Application: "test-app",
	}

	client := NewClient(cfg)
	require.NotNil(t, client)
	defer client.Close(context.TODO())

	elapsed := time.Since(start)

	// Should connect quickly (within 5 seconds) since container is ready
	assert.Less(t, elapsed, 5*time.Second, "connection should be fast when DB is ready")
}

func TestClient_Reconnection_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pgContainer := setupPostgreSQLContainer(t)
	defer pgContainer.Cleanup(t)

	cfg := &Config{
		DSN:         pgContainer.DSN,
		Application: "test-app",
	}

	// Create first client
	client1 := NewClient(cfg)
	require.NotNil(t, client1)

	ctx := context.Background()
	err := client1.HealthCheck(ctx)
	require.NoError(t, err)

	// Close first client
	client1.Close(context.TODO())

	// Create second client - should work fine
	client2 := NewClient(cfg)
	require.NotNil(t, client2)
	defer client2.Close(context.TODO())

	err = client2.HealthCheck(ctx)
	assert.NoError(t, err)
}

func TestClient_ConnectionPoolRecovery_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pgContainer := setupPostgreSQLContainer(t)
	defer pgContainer.Cleanup(t)

	cfg := &Config{
		DSN:            pgContainer.DSN,
		Application:    "test-app",
		MaxConnections: 2,
		MinConnections: 1,
	}

	client := NewClient(cfg)
	require.NotNil(t, client)
	defer client.Close(context.TODO())

	ctx := context.Background()

	// Acquire connections and release them
	for i := 0; i < 5; i++ {
		var result int
		err := client.Pool().QueryRow(ctx, "SELECT $1::int", i).Scan(&result)
		require.NoError(t, err)
		assert.Equal(t, i, result)
	}

	// Pool should still be healthy
	err := client.HealthCheck(ctx)
	assert.NoError(t, err)
}
