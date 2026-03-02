package pgclient

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgreSQL container configuration
const (
	postgresImage = "postgres:16-alpine"
	postgresUser  = "testuser"
	postgresPass  = "testpass"
	postgresDB    = "testdb"
)

// TestPostgreSQLContainer holds the container instance and connection details
type TestPostgreSQLContainer struct {
	container testcontainers.Container
	DSN       string
}

// setupPostgreSQLContainer creates and starts a PostgreSQL container for testing
func setupPostgreSQLContainer(t *testing.T) *TestPostgreSQLContainer {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        postgresImage,
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     postgresUser,
			"POSTGRES_PASSWORD": postgresPass,
			"POSTGRES_DB":       postgresDB,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
		).WithDeadline(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err, "failed to start PostgreSQL container")

	host, err := container.Host(ctx)
	require.NoError(t, err, "failed to get container host")

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err, "failed to get container port")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		postgresUser, postgresPass, host, port.Port(), postgresDB)

	return &TestPostgreSQLContainer{
		container: container,
		DSN:       dsn,
	}
}

// Cleanup terminates the container
func (tc *TestPostgreSQLContainer) Cleanup(t *testing.T) {
	ctx := context.Background()
	err := tc.container.Terminate(ctx)
	require.NoError(t, err, "failed to terminate container")
}

func TestNewClient_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pgContainer := setupPostgreSQLContainer(t)
	defer pgContainer.Cleanup(t)

	cfg := &Config{
		DSN:             pgContainer.DSN,
		Application:     "test-app",
		MaxConnections:  10,
		MinConnections:  2,
		MaxConnLifetime: 1 * time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
		SLO: SLO{
			Query: 100 * time.Millisecond,
		},
	}

	client := NewClient(cfg)
	require.NotNil(t, client)
	require.NotNil(t, client.Pool())
	defer client.Close(context.TODO())

	// Verify pool configuration
	pool := client.Pool()
	assert.NotNil(t, pool)
	assert.NotNil(t, pool.Config())
}

func TestClient_HealthCheck_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pgContainer := setupPostgreSQLContainer(t)
	defer pgContainer.Cleanup(t)

	cfg := &Config{
		DSN:         pgContainer.DSN,
		Application: "test-app",
	}

	client := NewClient(cfg)
	require.NotNil(t, client)
	defer client.Close(context.TODO())

	ctx := context.Background()
	err := client.HealthCheck(ctx)
	assert.NoError(t, err, "health check should pass")
}

func TestClient_HealthCheck_WithTimeout_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pgContainer := setupPostgreSQLContainer(t)
	defer pgContainer.Cleanup(t)

	cfg := &Config{
		DSN:         pgContainer.DSN,
		Application: "test-app",
	}

	client := NewClient(cfg)
	require.NotNil(t, client)
	defer client.Close(context.TODO())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.HealthCheck(ctx)
	assert.NoError(t, err, "health check should pass within timeout")
}

func TestClient_Query_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pgContainer := setupPostgreSQLContainer(t)
	defer pgContainer.Cleanup(t)

	cfg := &Config{
		DSN:         pgContainer.DSN,
		Application: "test-app",
	}

	client := NewClient(cfg)
	require.NotNil(t, client)
	defer client.Close(context.TODO())

	ctx := context.Background()

	// Test simple query
	var version string
	err := client.Pool().QueryRow(ctx, "SELECT version()").Scan(&version)
	require.NoError(t, err)
	assert.Contains(t, version, "PostgreSQL")
}

func TestClient_CreateTable_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pgContainer := setupPostgreSQLContainer(t)
	defer pgContainer.Cleanup(t)

	cfg := &Config{
		DSN:         pgContainer.DSN,
		Application: "test-app",
	}

	client := NewClient(cfg)
	require.NotNil(t, client)
	defer client.Close(context.TODO())

	ctx := context.Background()

	// Create a test table
	_, err := client.Pool().Exec(ctx, `
		CREATE TABLE test_users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	require.NoError(t, err)

	// Insert test data
	_, err = client.Pool().Exec(ctx,
		"INSERT INTO test_users (name, email) VALUES ($1, $2)",
		"John Doe", "john@example.com")
	require.NoError(t, err)

	// Query the data
	var name, email string
	err = client.Pool().QueryRow(ctx,
		"SELECT name, email FROM test_users WHERE email = $1",
		"john@example.com").Scan(&name, &email)
	require.NoError(t, err)
	assert.Equal(t, "John Doe", name)
	assert.Equal(t, "john@example.com", email)
}

func TestClient_Transaction_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pgContainer := setupPostgreSQLContainer(t)
	defer pgContainer.Cleanup(t)

	cfg := &Config{
		DSN:         pgContainer.DSN,
		Application: "test-app",
	}

	client := NewClient(cfg)
	require.NotNil(t, client)
	defer client.Close(context.TODO())

	ctx := context.Background()

	// Create a test table
	_, err := client.Pool().Exec(ctx, `
		CREATE TABLE test_accounts (
			id SERIAL PRIMARY KEY,
			balance INTEGER NOT NULL
		)
	`)
	require.NoError(t, err)

	// Test successful transaction
	tx, err := client.Pool().Begin(ctx)
	require.NoError(t, err)

	_, err = tx.Exec(ctx, "INSERT INTO test_accounts (balance) VALUES ($1)", 100)
	require.NoError(t, err)

	err = tx.Commit(ctx)
	require.NoError(t, err)

	// Verify data was committed
	var balance int
	err = client.Pool().QueryRow(ctx,
		"SELECT balance FROM test_accounts WHERE id = 1").Scan(&balance)
	require.NoError(t, err)
	assert.Equal(t, 100, balance)

	// Test rollback
	tx, err = client.Pool().Begin(ctx)
	require.NoError(t, err)

	_, err = tx.Exec(ctx, "UPDATE test_accounts SET balance = 200 WHERE id = 1")
	require.NoError(t, err)

	err = tx.Rollback(ctx)
	require.NoError(t, err)

	// Verify data was not changed
	err = client.Pool().QueryRow(ctx,
		"SELECT balance FROM test_accounts WHERE id = 1").Scan(&balance)
	require.NoError(t, err)
	assert.Equal(t, 100, balance, "balance should not have changed after rollback")
}

func TestClient_MultipleConnections_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pgContainer := setupPostgreSQLContainer(t)
	defer pgContainer.Cleanup(t)

	cfg := &Config{
		DSN:            pgContainer.DSN,
		Application:    "test-app",
		MaxConnections: 5,
		MinConnections: 2,
	}

	client := NewClient(cfg)
	require.NotNil(t, client)
	defer client.Close(context.TODO())

	ctx := context.Background()

	// Execute multiple concurrent queries
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			var result int
			err := client.Pool().QueryRow(ctx, "SELECT $1::int", id).Scan(&result)
			assert.NoError(t, err)
			assert.Equal(t, id, result)
			done <- true
		}(i)
	}

	// Wait for all queries to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestClient_Close_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pgContainer := setupPostgreSQLContainer(t)
	defer pgContainer.Cleanup(t)

	cfg := &Config{
		DSN:         pgContainer.DSN,
		Application: "test-app",
	}

	client := NewClient(cfg)
	require.NotNil(t, client)

	// Verify connection is working
	ctx := context.Background()
	err := client.HealthCheck(ctx)
	require.NoError(t, err)

	// Close the client
	client.Close(context.TODO())

	// Verify pool is closed (operations should fail)
	err = client.HealthCheck(ctx)
	assert.Error(t, err, "health check should fail after close")
}

func TestClient_ConfigWithSLO_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pgContainer := setupPostgreSQLContainer(t)
	defer pgContainer.Cleanup(t)

	cfg := &Config{
		DSN:         pgContainer.DSN,
		Application: "test-app",
		SLO: SLO{
			Query: 50 * time.Millisecond, // Very low threshold for testing
		},
	}

	client := NewClient(cfg)
	require.NotNil(t, client)
	defer client.Close(context.TODO())

	ctx := context.Background()

	// Execute a query - SLO tracer should be active
	var result int
	err := client.Pool().QueryRow(ctx, "SELECT 1").Scan(&result)
	require.NoError(t, err)
	assert.Equal(t, 1, result)
}

func TestClient_PoolConfiguration_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pgContainer := setupPostgreSQLContainer(t)
	defer pgContainer.Cleanup(t)

	tests := []struct {
		name           string
		maxConnections int32
		minConnections int32
	}{
		{
			name:           "default connections",
			maxConnections: 0,
			minConnections: 0,
		},
		{
			name:           "custom connections",
			maxConnections: 20,
			minConnections: 5,
		},
		{
			name:           "minimal connections",
			maxConnections: 2,
			minConnections: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				DSN:            pgContainer.DSN,
				Application:    "test-app",
				MaxConnections: tt.maxConnections,
				MinConnections: tt.minConnections,
			}

			client := NewClient(cfg)
			require.NotNil(t, client)
			defer client.Close(context.TODO())

			// Verify connection works
			ctx := context.Background()
			err := client.HealthCheck(ctx)
			assert.NoError(t, err)
		})
	}
}
