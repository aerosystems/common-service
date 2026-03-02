package pgclient

import "context"

// HealthCheck checks the health of the database connection.
func (c *Client) HealthCheck(ctx context.Context) error {
	return c.pool.Ping(ctx)
}
