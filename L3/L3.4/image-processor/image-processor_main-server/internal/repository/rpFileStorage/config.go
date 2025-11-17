package rpFileStorage

import (
	"fmt"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	Dir string
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Dir: cfg.GetString("app.repository.file_storage.dir"),
	}
}

func (c Config) String() string {
	return fmt.Sprintf(`file storage: %s`,
		c.Dir)
}
