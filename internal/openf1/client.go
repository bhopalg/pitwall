package openf1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func New() *Client {
	return &Client{
		baseURL: "https://api.openf1.org/v1",
		http:    &http.Client{},
	}
}

func (c *Client) Get(ctx context.Context, path string, q url.Values, out any) error {
	u := c.baseURL + path
	if len(q) > 0 {
		u += "?" + q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opend1: %s returned %d", path, resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
