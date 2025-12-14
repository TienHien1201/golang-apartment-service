package auth

import "time"

type Config struct {
	AccessSecret  string
	AccessExpire  time.Duration
	RefreshSecret string
	RefreshExpire time.Duration
}
