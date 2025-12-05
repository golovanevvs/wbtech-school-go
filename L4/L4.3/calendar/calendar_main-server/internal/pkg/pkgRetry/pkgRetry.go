package pkgRetry

import (
	"github.com/wb-go/wbf/retry"
)

type Retry retry.Strategy

func New(cfg *Config) *Retry {
	return &Retry{
		Attempts: cfg.Attempts,
		Delay:    cfg.Delay,
		Backoff:  cfg.Backoff,
	}
}
