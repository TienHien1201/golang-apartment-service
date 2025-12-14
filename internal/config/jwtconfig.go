// internal/config/jwtconfig.go
package config

import (
	"fmt"
	"time"
)

type JWTConfig struct {
	AccessSecret  string        `mapstructure:"access_secret" env:"JWT_ACCESS_SECRET" validate:"required,min=32"`
	AccessExpire  time.Duration `mapstructure:"access_expire" env:"JWT_ACCESS_EXPIRE" envDefault:"168h"` // 7 ngày
	RefreshSecret string        `mapstructure:"refresh_secret" env:"JWT_REFRESH_SECRET" validate:"required,min=32"`
	RefreshExpire time.Duration `mapstructure:"refresh_expire" env:"JWT_REFRESH_EXPIRE" envDefault:"720h"` // 30 ngày
	TokenType     string        `mapstructure:"token_type" env:"JWT_TOKEN_TYPE" envDefault:"Bearer"`
}

func (c JWTConfig) Validate() error {
	if len(c.AccessSecret) < 32 {
		return fmt.Errorf("JWT_ACCESS_SECRET must be at least 32 characters")
	}
	if len(c.RefreshSecret) < 32 {
		return fmt.Errorf("JWT_REFRESH_SECRET must be at least 32 characters")
	}
	if c.AccessExpire <= 0 {
		return fmt.Errorf("JWT_ACCESS_EXPIRE must be positive")
	}
	if c.RefreshExpire <= 0 {
		return fmt.Errorf("JWT_REFRESH_EXPIRE must be positive")
	}
	return nil
}
