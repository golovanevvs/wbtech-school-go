package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	vip        *viper.Viper
	Server     Server
	Logger     Logger
	Repository Repository
	Handler    Handler
	Kafka      Kafka
}

type Server struct {
	Addr string
}

type Logger struct {
	LogLevel string
}

type Repository struct {
	Postgres postgres
}

type postgres struct {
	DatabaseDSN     string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type Handler struct {
	GinMode string
}

type Kafka struct {
	Brokers             []string
	ClientID            string
	ConsumerGroup       string
	Version             string
	RetryMax            int
	RequiredAcks        int    // NoResponse=0, WaitForLocal=1, WaitForAll=-1
	Partitioner         string //"roundrobin", "hash"
	EnableReturnSuccess bool
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

	c.Server.Addr = c.vip.GetString("server.addr")

	c.Logger.LogLevel = c.vip.GetString("logging.level")

	postgresHost := c.vip.GetString("postgres.host")
	postgresPort := c.vip.GetString("postgres.port")
	postgresUser := c.vip.GetString("postgres.user")
	postgresPassword := c.vip.GetString("postgres.password")
	postgresDB := c.vip.GetString("postgres.db")
	c.Repository.Postgres.DatabaseDSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", postgresHost, postgresPort, postgresUser, postgresPassword, postgresDB)
	c.Repository.Postgres.MaxOpenConns = c.vip.GetInt("postgres.max_open_conns")
	c.Repository.Postgres.MaxIdleConns = c.vip.GetInt("postgres.max_idle_conns")
	c.Repository.Postgres.ConnMaxLifetime = c.vip.GetDuration("postgres.conn_max_lifetime")

	c.Handler.GinMode = c.vip.GetString("gin.mode")

	c.Kafka.Brokers = c.vip.GetStringSlice("kafka.brokers")
	c.Kafka.ClientID = c.vip.GetString("kafka.client_id")
	c.Kafka.ConsumerGroup = c.vip.GetString("kafka.consumer_group")
	c.Kafka.Version = c.vip.GetString("kafka.version")
	c.Kafka.RetryMax = c.vip.GetInt("kafka.retry_max")
	c.Kafka.RequiredAcks = c.vip.GetInt("kafka.required_acks")
	c.Kafka.Partitioner = c.vip.GetString("kafka.partitioner")
	c.Kafka.EnableReturnSuccess = c.vip.GetBool(("kafka.enable_return_success"))

	return nil
}
