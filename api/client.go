package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/hashicorp/go-cleanhttp"
	"golang.org/x/net/http2"
)

const (
	DefaultAddress = "https://localhost:6001/api"
	AuthHeader     = "Authorization"
)

// Config is used to configure the creation of the client.
type Config struct {
	AuthToken  string
	Address    string
	HttpClient *http.Client
}

type Client struct {
	addr    *url.URL
	cfg     *Config
	headers http.Header
}

func DefaultConfig() (*Config, error) {
	cfg := &Config{
		Address:    DefaultAddress,
		HttpClient: cleanhttp.DefaultPooledClient(),
	}
	transport := cfg.HttpClient.Transport.(*http.Transport)

	if err := http2.ConfigureTransport(transport); err != nil {
		return nil, err
	}
	return cfg, nil
}

func NewClient(c *Config) (*Client, error) {
	def, err := DefaultConfig()
	if err != nil {
		return nil, err
	}
	if c == nil {
		c = def
	}

	if c.AuthToken == "" {
		return nil, errors.New("missing auth token for openchangelog api client")
	}

	if c.Address == "" {
		c.Address = def.Address
	}
	if c.HttpClient == nil {
		c.HttpClient = def.HttpClient
	}
	if c.HttpClient.Transport == nil {
		c.HttpClient.Transport = def.HttpClient.Transport
	}

	u, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}

	client := &Client{
		addr:    u,
		cfg:     c,
		headers: make(http.Header),
	}
	client.headers[AuthHeader] = []string{fmt.Sprintf("Bearer %s", c.AuthToken)}
	return client, err
}

func (c *Client) NewRequest(ctx context.Context, method, requestPath string, body io.Reader) (*http.Request, error) {
	p, err := url.Parse(requestPath)
	if err != nil {
		return nil, err
	}

	url := &url.URL{
		User:     c.addr.User,
		Scheme:   c.addr.Scheme,
		Host:     c.addr.Host,
		Path:     path.Join(c.addr.Path, p.Path),
		RawQuery: p.RawQuery,
	}

	req, err := http.NewRequestWithContext(ctx, method, url.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header = c.headers
	return req, nil
}

func (c *Client) rawRequestWithContext(r *http.Request) (*Response, error) {
	httpClient := c.cfg.HttpClient

	var result *Response
	resp, err := httpClient.Do(r)
	if resp != nil {
		result = &Response{Response: resp}
	}
	if err != nil {
		return result, err
	}

	if err := result.Error(); err != nil {
		return result, err
	}

	return result, nil
}
