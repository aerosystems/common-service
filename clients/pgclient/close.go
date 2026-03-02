package pgclient

import "context"

// Close closes the connection pool.
func (c *Client) Close(_ context.Context) error {
	c.pool.Close()
	return nil
}
