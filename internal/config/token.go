package config

import "time"

type TokenConfig struct {
	AccessSecret  string
	AccessExpire  time.Duration
	RefreshSecret string
	RefreshExpire time.Duration
}
