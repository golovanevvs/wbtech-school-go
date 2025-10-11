package telegram

import (
	"fmt"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	Token  string
	ChatID int64
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Token:  cfg.GetString("telegram.token"),
		ChatID: int64(cfg.GetInt("telegram.chat_id")),
	}
}

func (c Config) String() string {
	return fmt.Sprintf("telegram:\n \033[33mtoken: \033[0m\033[32m%s\033[0m, \033[33mchatID: \033[0m\033[32m%d\033[0m", c.Token, c.ChatID)
}
