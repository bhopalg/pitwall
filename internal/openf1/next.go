package openf1

import (
	"context"
	"net/url"
)

func (c *Client) Next(ctx context.Context) (*Session, error) {
	q := url.Values{}
	q.Set("session_key", "latest")

	var sessions []Session
	if err := c.Get(ctx, "/sessions", q, &sessions); err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return nil, nil
	}

	return &sessions[0], nil
}
