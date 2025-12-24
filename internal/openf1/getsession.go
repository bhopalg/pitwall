package openf1

import (
	"context"
	"net/url"
)

func (c *Client) GetSession(ctx context.Context, country_name, session_name, year string) (*Session, error) {
	q := url.Values{}
	q.Set("country_name", country_name)

	if session_name != "" {
		q.Set("session_name", session_name)
	}

	q.Set("year", year)

	var sessions []Session
	if err := c.Get(ctx, "/sessions", q, &sessions); err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return nil, nil
	}

	return &sessions[0], nil
}
