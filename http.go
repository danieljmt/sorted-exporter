package sorted

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	defaultClientTImeout    = 30 * time.Second
	defaultTransportTimeout = 5 * time.Second
	defaultIdleTimeout      = 5 * time.Second
)

type Client struct {
	c       *http.Client
	url     string
	headers map[string]string
}

func NewClient(baseURL string, opt ...ClientOpt) (*Client, error) {
	c := &Client{
		headers: map[string]string{
			"Content-Type": "application/json",
		},
		url: baseURL,
		c: &http.Client{
			Timeout: defaultClientTImeout,
			Transport: &http.Transport{
				Dial:                (&net.Dialer{Timeout: defaultTransportTimeout}).Dial,
				TLSHandshakeTimeout: defaultTransportTimeout,
				IdleConnTimeout:     defaultIdleTimeout,
			},
		},
	}

	return c, nil
}

type ClientOpt func(*Client) error

func ClientOptDebug() ClientOpt {
	return func(c *Client) error {
		return nil
	}
}

func (c *Client) SetHeader(k, v string) {
	c.headers[k] = v
}

func (c *Client) Do(ctx context.Context, method, path string, payload, response interface{}) error {
	var body bytes.Buffer
	if payload != nil {
		if err := json.NewEncoder(&body).Encode(payload); err != nil {
			return fmt.Errorf("failed to encode body: %w", err)
		}
	}

	req, err := http.NewRequest(method, c.url+path, &body)
	if err != nil {
		return fmt.Errorf("invalid http request: %w", err)
	}
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("failed executing http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 == 2 {
		if response == nil {
			return nil
		}

		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
		return nil
	}

	return fmt.Errorf("http request returned error code %d", resp.StatusCode)
}
