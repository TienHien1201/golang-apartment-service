package xes

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
)

type Config struct {
	Addresses     []string
	Username      string
	Password      string
	APIKey        string
	Header        http.Header
	EnableLogging bool
}

type Client struct {
	Client *elasticsearch.Client
}

func NewClient(cfg *Config) (*Client, error) {
	esConfig := elasticsearch.Config{
		Addresses:     cfg.Addresses,
		Username:      cfg.Username,
		Password:      cfg.Password,
		APIKey:        cfg.APIKey,
		Header:        cfg.Header,
		MaxRetries:    5,
		RetryOnStatus: []int{502, 503, 504, 429},
	}

	var transport http.RoundTripper
	transport = &http.Transport{
		ResponseHeaderTimeout: time.Second * 2,
		DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
		MaxIdleConnsPerHost:   10,
	}

	if cfg.EnableLogging {
		transport = &LoggingTransport{Transport: transport}
	}

	esConfig.Transport = transport

	client, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	if _, err := client.Info(client.Info.WithContext(context.Background()), client.Info.WithHuman()); err != nil {
		return nil, fmt.Errorf("failed to ping elasticsearch: %w", err)
	}

	return &Client{Client: client}, nil
}

func (c Client) Close() error {
	return nil
}

// --- LoggingTransport ---

type LoggingTransport struct {
	Transport http.RoundTripper
}

func (lt *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Printf("ES REQUEST: %s %s\n", req.Method, req.URL)

	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		log.Printf("ES REQUEST BODY: %s\n", string(body))
		req.Body = io.NopCloser(bytes.NewReader(body))
	}

	resp, err := lt.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("ES RESPONSE BODY: %s\n", string(body))
		resp.Body = io.NopCloser(bytes.NewReader(body)) // Reset body để response không bị mất
	}

	return resp, nil
}
