package config

import (
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	vip     *viper.Viper
	App     App
	Server  Server
	Logger  Logger
	Handler Handler
}

type App struct {
	DelayBeforeClosing time.Duration
}

type Server struct {
	Addr    string
	Timeout time.Duration
}

type Logger struct {
	LogLevel string
}

type Handler struct {
	GinMode string
}

func New(pathEnvFile string) *Config {
	vip := viper.New()

	cfg := &Config{
		vip: vip,
		App: App{
			DelayBeforeClosing: 5 * time.Second,
		},
		Server: Server{
			Timeout: 5 * time.Second,
		},
	}

	_ = godotenv.Load(pathEnvFile)

	pflag.StringP("server.addr", "a", "localhost:5470", "Server address")
	pflag.StringP("logging.level", "l", "debug", "Logging level")
	pflag.StringP("gin.mode", "m", "release", "Gin mode")
	pflag.Parse()

	cfg.vip.BindPFlags(pflag.CommandLine)

	cfg.vip.AutomaticEnv()
	cfg.vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfg.Server.Addr = cfg.vip.GetString("server.addr")
	cfg.Logger.LogLevel = cfg.vip.GetString("logging.level")
	cfg.Handler.GinMode = cfg.vip.GetString("gin.mode")

	return cfg
}
