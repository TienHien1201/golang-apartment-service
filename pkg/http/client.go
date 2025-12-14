package xhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
	MethodPatch  = "PATCH"
)

// ClientOption defines a function type for configuring HTTPClient.
type ClientOption func(*HTTPClient)

// ClientRequestOptions holds the options for making an HTTP request.
type ClientRequestOptions struct {
	Method      string
	URL         string
	Headers     map[string]string
	QueryParams map[string][]string
	Body        interface{}
}

// HTTPClient represents an HTTP client with a configurable timeout.
type HTTPClient struct {
	timeout    time.Duration
	httpClient *http.Client
}

// NewHTTPClient creates a new HTTPClient with the given options.
func NewHTTPClient(opts ...ClientOption) *HTTPClient {
	client := &HTTPClient{
		timeout: 20 * time.Second,
	}

	for _, opt := range opts {
		opt(client)
	}

	client.httpClient = &http.Client{Timeout: client.timeout}
	return client
}

// SendRequest sends an HTTP request based on the given options and returns the response.
func (c *HTTPClient) SendRequest(opts *ClientRequestOptions) (*http.Response, error) {
	req, err := c.buildRequest(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// SendAndParse sends an HTTP request and parses the response into the given responseBody.
func (c *HTTPClient) SendAndParse(opts *ClientRequestOptions, responseBody interface{}) error {
	resp, err := c.SendRequest(opts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		return fmt.Errorf("unexpected status code [%d]: %s", resp.StatusCode, bodyBytes)
	}

	if responseBody == nil {
		return nil
	}

	switch v := responseBody.(type) {
	case *[]byte:
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		*v = bodyBytes
	case io.Writer:
		_, err := io.Copy(v, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to copy response body: %w", err)
		}
	default:
		if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}
	}

	return nil
}

// HTTPClient returns the underlying http.Client.
func (c *HTTPClient) HTTPClient() *http.Client {
	return c.httpClient
}

// buildRequest builds an HTTP request based on the given options.
func (c *HTTPClient) buildRequest(opts *ClientRequestOptions) (*http.Request, error) {
	body, err := c.createRequestBody(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %w", err)
	}

	req, err := http.NewRequest(opts.Method, opts.URL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.addQueryParams(req, opts.QueryParams)
	c.addHeaders(req, opts.Headers)

	return req, nil
}

func (c *HTTPClient) createRequestBody(opts *ClientRequestOptions) (io.Reader, error) {
	switch v := opts.Body.(type) {
	case *[]byte:
		return bytes.NewBuffer(*v), nil
	case io.Reader:
		return v, nil
	default:
		if formData, ok := opts.Body.(map[string]string); ok && opts.Headers["Content-Type"] == "application/x-www-form-urlencoded" {
			values := url.Values{}
			for key, value := range formData {
				values.Set(key, value)
			}
			return strings.NewReader(values.Encode()), nil
		}
		if opts.Headers["Content-Type"] == "multipart/form-data" {
			return v.(io.Reader), nil
		}
		jsonBody, err := json.Marshal(opts.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		return bytes.NewBuffer(jsonBody), nil
	}
}

func (c *HTTPClient) addQueryParams(req *http.Request, queryParams map[string][]string) {
	if len(queryParams) > 0 {
		query := req.URL.Query()
		for key, values := range queryParams {
			for _, value := range values {
				query.Add(key, value)
			}
		}
		req.URL.RawQuery = query.Encode()
	}
}

func (c *HTTPClient) addHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
}

// WithTimeout sets the timeout for the HTTPClient.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *HTTPClient) {
		c.timeout = timeout
	}
}
