package openf1

import (
	"context"
	"net/url"
)

func (c *Client) GetSession(ctx context.Context) (*Session, error) {
	q := url.Values{}
	q.Set("country_name", "Belgium")
	q.Set("session_name", "Sprint")
	q.Set("year", "2023")

	var sessions []Session
	if err := c.Get(ctx, "/sessions", q, &sessions); err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return nil, nil
	}

	return &sessions[0], nil
}
