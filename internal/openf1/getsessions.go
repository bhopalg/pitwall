package openf1

import (
	"context"
	"net/url"
)

func (c *Client) GetSessions(ctx context.Context, country_name, year string) (*[]Session, error) {
	q := url.Values{}
	q.Set("country_name", country_name)
	q.Set("year", year)

	var sessions []Session
	if err := c.Get(ctx, "/sessions", q, &sessions); err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return nil, nil
	}

	return &sessions, nil
}
