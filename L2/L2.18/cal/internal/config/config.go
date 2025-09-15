package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	vip     *viper.Viper
	Server  Server
	Logger  Logger
	Handler Handler
}

type Server struct {
	Addr string
}

type Logger struct {
	LogLevel string
}

type Handler struct {
	GinMode string
}

func New() *Config {
	vip := viper.New()

	return &Config{
		vip: vip,
	}
}

func (c *Config) Load(pathEnvFile string) error {
	_ = godotenv.Load(pathEnvFile)

	pflag.StringP("server.addr", "a", "localhost:5470", "Server address")
	pflag.StringP("logging.level", "l", "info", "Logging level")
	pflag.StringP("gin.mode", "m", "release", "Gin mode")
	pflag.Parse()

	c.vip.BindPFlags(pflag.CommandLine)

	c.vip.AutomaticEnv()
	c.vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	c.Server.Addr = c.vip.GetString("server.addr")
	c.Logger.LogLevel = c.vip.GetString("logging.level")
	c.Handler.GinMode = c.vip.GetString("gin.mode")

	return nil
}
