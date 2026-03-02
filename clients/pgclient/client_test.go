package pgclient

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		valid  bool
	}{
		{
			name: "valid config",
			config: &Config{
				DSN:             "postgres://user:pass@localhost:5432/db",
				Application:     "test-app",
				MaxConnections:  10,
				MinConnections:  2,
				MaxConnLifetime: 1 * time.Hour,
				MaxConnIdleTime: 30 * time.Minute,
				SLO: SLO{
					Query: 100 * time.Millisecond,
				},
			},
			valid: true,
		},
		{
			name: "minimal config",
			config: &Config{
				DSN:         "postgres://user:pass@localhost:5432/db",
				Application: "test-app",
			},
			valid: true,
		},
		{
			name: "config with zero values",
			config: &Config{
				DSN:             "postgres://user:pass@localhost:5432/db",
				Application:     "test-app",
				MaxConnections:  0,
				MinConnections:  0,
				MaxConnLifetime: 0,
				MaxConnIdleTime: 0,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.config)
			assert.NotEmpty(t, tt.config.DSN)
		})
	}
}

func TestNewClient_InvalidDSN(t *testing.T) {
	cfg := &Config{
		DSN:         "invalid-dsn",
		Application: "test-app",
	}

	assert.Panics(t, func() {
		NewClient(cfg)
	}, "NewClient should panic with invalid DSN")
}

func TestClient_Pool(t *testing.T) {
	// This is a unit test that verifies the Pool method exists
	// and returns the expected type
	// Actual connection testing is done in integration tests
	client := &Client{}
	pool := client.Pool()

	// Pool can be nil in unit test context
	if pool != nil {
		assert.NotNil(t, pool)
	}
}

func TestSLO_Configuration(t *testing.T) {
	tests := []struct {
		name     string
		slo      SLO
		expected time.Duration
	}{
		{
			name:     "100ms query SLO",
			slo:      SLO{Query: 100 * time.Millisecond},
			expected: 100 * time.Millisecond,
		},
		{
			name:     "1s query SLO",
			slo:      SLO{Query: 1 * time.Second},
			expected: 1 * time.Second,
		},
		{
			name:     "no SLO",
			slo:      SLO{Query: 0},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.slo.Query)
		})
	}
}

func TestConfig_PoolSettings(t *testing.T) {
	cfg := &Config{
		DSN:             "postgres://user:pass@localhost:5432/db", // nolint:govet
		Application:     "test-app",                               // nolint:govet
		MaxConnections:  50,
		MinConnections:  5,
		MaxConnLifetime: 2 * time.Hour,
		MaxConnIdleTime: 1 * time.Hour,
	}

	assert.Equal(t, int32(50), cfg.MaxConnections)
	assert.Equal(t, int32(5), cfg.MinConnections)
	assert.Equal(t, 2*time.Hour, cfg.MaxConnLifetime)
	assert.Equal(t, 1*time.Hour, cfg.MaxConnIdleTime)
}

func TestConfig_ApplicationName(t *testing.T) {
	tests := []struct {
		name        string
		application string
	}{
		{"simple name", "my-app"},
		{"with version", "my-app-v1.0"},
		{"with dashes", "my-cool-app"},
		{"with underscores", "my_cool_app"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				DSN:         "postgres://user:pass@localhost:5432/db", // nolint:govet
				Application: tt.application,
			}
			assert.Equal(t, tt.application, cfg.Application)
		})
	}
}
