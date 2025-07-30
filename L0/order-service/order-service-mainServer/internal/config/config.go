package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	vip    *viper.Viper
	Server server
	Logger logger
}

type server struct {
	Addr string
}

type logger struct {
	LogLevel string
}

func New() *Config {
	vip := viper.New()

	return &Config{
		vip: vip,
	}
}

func (c *Config) Load(pathConfigFile string, pathEnvFile string, envPrefix string) error {
	godotenv.Load(pathEnvFile)

	c.vip.AutomaticEnv()

	if envPrefix != "" {
		c.vip.SetEnvPrefix(envPrefix)
	}

	c.vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	c.vip.SetConfigFile(pathConfigFile)
	err := c.vip.ReadInConfig()
	if err != nil {
		return err
	}

	c.Server.Addr = c.GetString("server.addr")

	c.Logger.LogLevel = c.GetString("logging.level")

	return nil
}

func (c *Config) GetString(key string) string {
	return c.vip.GetString(key)
}

func (c *Config) GetInt(key string) int {
	return c.vip.GetInt(key)
}

func (c *Config) GetBool(key string) bool {
	return c.vip.GetBool(key)
}
