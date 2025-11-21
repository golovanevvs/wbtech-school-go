package service

import (
	"fmt"
	"time"

	"github.com/wb-go/wbf/config"
)

// Config service configuration
type Config struct {
	JWTSecret          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		JWTSecret:          cfg.GetString("app.service.jwt_secret"),
		AccessTokenExpiry:  cfg.GetDuration("app.service.jwt_access_token_expiry"),
		RefreshTokenExpiry: cfg.GetDuration("app.service.jwt_refresh_token_expiry"),
	}
}

func (c Config) String() string {
	var jwtSecret string
	if c.JWTSecret != "" {
		jwtSecret = "***"
	}
	return fmt.Sprintf(`service:
  %s: %s
  %s: %s
  %s: %s 
  `,
		"JWT secret", jwtSecret,
		"JWT access token expiry", c.AccessTokenExpiry,
		"JWT refresh token expiry", c.RefreshTokenExpiry,
	)
}
