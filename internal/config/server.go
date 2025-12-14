package config

import "time"

type ServerConfig struct {
	HTTP HTTPConfig
}

type HTTPConfig struct {
	Host    string
	Port    int
	Timeout time.Duration
}
