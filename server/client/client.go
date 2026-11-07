package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	base    *url.URL
	token   string
	http    *http.Client
	OnToken func(string) error
}

func New(endpoint, token string, timeout time.Duration) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil || u == nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, fmt.Errorf("server must be an http(s) URL without credentials, query or fragment")
	}

	if timeout <= 0 {
		return nil, fmt.Errorf("timeout must be positive")
	}
	return &Client{base: u, token: token, http: &http.Client{Timeout: timeout, CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }}}, nil
}

func (c *Client) Get(ctx context.Context, path string, query url.Values) (json.RawMessage, error) {
	if c.token == "" {
		return nil, fmt.Errorf("run gshark-cli auth login first")
	}
	return c.request(ctx, http.MethodGet, path, query, nil)
}

func (c *Client) Mutate(ctx context.Context, method, path string, payload interface{}) (json.RawMessage, error) {
	if c.token == "" {
		return nil, fmt.Errorf("run gshark-cli auth login first")
	}
	if method != http.MethodPost && method != http.MethodPut && method != http.MethodDelete {
		return nil, fmt.Errorf("unsupported mutation method")
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return c.request(ctx, method, path, nil, body)
}

func (c *Client) Post(ctx context.Context, path string, payload interface{}) (json.RawMessage, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return c.request(ctx, http.MethodPost, path, nil, body)
}

func (c *Client) request(ctx context.Context, method, path string, query url.Values, body []byte) (json.RawMessage, error) {
	u := *c.base
	u.Path = strings.TrimRight(u.Path, "/") + "/" + strings.TrimLeft(path, "/")
	u.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-token", c.token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed; check server connectivity: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned HTTP %d; check URL and credentials", res.StatusCode)
	}
	const maxBody = 8 << 20
	responseBody, err := io.ReadAll(io.LimitReader(res.Body, maxBody+1))
	if err != nil {
		return nil, err
	}
	if len(responseBody) > maxBody {
		return nil, fmt.Errorf("response exceeds 8 MiB; reduce --limit")
	}
	var envelope struct {
		Code *int            `json:"code"`
		Data json.RawMessage `json:"data"`
		Msg  string          `json:"msg"`
	}
	if err := json.Unmarshal(responseBody, &envelope); err != nil || envelope.Code == nil || len(envelope.Data) == 0 {
		return nil, fmt.Errorf("invalid GShark response; check --server (include /api when using the frontend proxy)")
	}
	if *envelope.Code != 0 {
		return nil, fmt.Errorf("GShark API error: %s", envelope.Msg)
	}
	if token := res.Header.Get("new-token"); token != "" && c.OnToken != nil {
		if err := c.OnToken(token); err != nil {
			return nil, fmt.Errorf("save refreshed credentials: %w", err)
		}
		c.token = token
	}
	return envelope.Data, nil
}
