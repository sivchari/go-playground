package playground

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	FrontendURL = "https://go.dev/play"
	BackendURL  = "https://play.golang.org"
)

type Client struct {
	client  *http.Client
	version string
}

func NewClient(version string) *Client {
	return &Client{
		client:  &http.Client{},
		version: version,
	}
}

type RunResult struct {
	Errors     string
	Events     []*RunEvent
	IsTest     bool
	Status     int
	TestFailed bool
}

type RunEvent struct {
	Message string
	Kind    string
	Delay   time.Duration
}

func backendURL() string {
	if val, ok := os.LookupEnv("PLAYGROUND_BACKEND_URL"); ok {
		return val
	}
	return BackendURL
}

func frontendURL() string {
	if val, ok := os.LookupEnv("PLAYGROUND_FRONTEND_URL"); ok {
		return val
	}
	return FrontendURL
}

func (c *Client) Run(src []byte) (*RunResult, error) {
	p, err := url.JoinPath(backendURL(), "compile")
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(p)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	query.Set("version", c.version)
	body := string(src)
	query.Set("body", body)

	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result *RunResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

type FormatResult struct {
	Body  string
	Error string
}

func (c *Client) Format(src []byte, imports bool) (*FormatResult, error) {
	p, err := url.JoinPath(backendURL(), "fmt")
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(p)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	query.Set("imports", fmt.Sprintf("%t", imports))
	body := string(src)
	query.Set("body", body)
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result *FormatResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) Share(src []byte) (*url.URL, error) {
	p, err := url.JoinPath(backendURL(), "share")
	if err != nil {
		return nil, err
	}

	payload := bytes.NewReader(src)

	req, err := http.NewRequest(http.MethodPost, p, payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	shareURL, err := url.JoinPath(frontendURL(), "p", string(res))
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(shareURL)
	if err != nil {
		return nil, err
	}

	return u, nil
}

type VersionResult struct {
	Version string
	Release string
	Name    string
}

func (c *Client) Version() (*VersionResult, error) {
	p, err := url.JoinPath(backendURL(), "version")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, p, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result *VersionResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
