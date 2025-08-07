package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	rediscache "github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/cache/redis"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/config"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/kafka"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/logger/zlog"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/transport/http/handler"
	"github.com/rs/zerolog"
)

type app struct {
	config  *config.Config
	handler *handler.Handler
	logger  *zerolog.Logger
	rp      *repository.Repository
}

func Run() {
	zlog.Init()
	zlog.Logger.Info().Msg("starting order-service-mainServer...")

	configFile := "config.yaml"
	envFile := ".env"
	configPath := fmt.Sprintf("./config/%s", configFile)
	configDefaultFile := "default-config.yaml"
	configDefaultPath := fmt.Sprintf("./config/%s", configDefaultFile)
	envPath := fmt.Sprintf("./%s", envFile)

	zlog.Logger.Info().Str("file", configFile).Msg("loading configuration...")

	cfg := config.New()
	err := cfg.Load(configPath, envPath, "")
	if err != nil {
		zlog.Logger.Error().Err(err).Str("file", configFile).Msg("failed to load configuration")
		zlog.Logger.Warn().Str("file", configDefaultFile).Msg("loading default configuration...")
		err := cfg.Load(configDefaultPath, envPath, "")
		if err != nil {
			zlog.Logger.Fatal().Err(err).Msg("failed to load default configuration")
		}
	}
	zlog.Logger.Info().Msg("configuration loaded successfully")

	logLevelStr := cfg.Logger.LogLevel
	logLevel, err := zlog.ParseLogLevel(logLevelStr)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to parse log level")
	}
	zlog.Logger.Info().Str("logLevel", logLevel.String()).Msg("logging level")
	zlog.Logger = zlog.Logger.Level(logLevel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rp, err := repository.New(ctx, &cfg.Repository, &zlog.Logger)
	if err != nil {
		os.Exit(1)
	}

	rd, err := rediscache.New(&cfg.RedisCache, &zlog.Logger)
	if err != nil {
		os.Exit(1)
	}

	h := handler.New(&cfg.Handler, &zlog.Logger, rp, rd)

	h.InitRoutes()

	go func() {
		zlog.Logger.Info().Str("address", cfg.Server.Addr).Msg("starting HTTP server")
		if err := h.Run(cfg.Server.Addr); err != nil {
			zlog.Logger.Error().Err(err).Msg("failed to start HTTP server")
			cancel()
		}
	}()

	k, err := kafka.New(&cfg.Kafka, &zlog.Logger)
	if err != nil {
		os.Exit(1)
	}

	topics := []string{cfg.Kafka.Topic}

	kafkaHandler := func(ctx context.Context, msg *sarama.ConsumerMessage) error {
		var order model.Order
		err := json.Unmarshal(msg.Value, &order)
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("error decoding message to model")
		}

		rp.AddOrder(ctx, order)

		return nil
	}

	consumerGroup, err := kafka.NewConsumerGroup(k, topics, kafkaHandler)
	if err != nil {
		os.Exit(1)
	}
	defer consumerGroup.Close()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	if err := consumerGroup.Start(ctx); err != nil {
		os.Exit(1)
	}

	<-ctx.Done()
	zlog.Logger.Warn().Msg("received shutdown signal")

	zlog.Logger.Warn().Msg("shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: h.Router,
	}

	rp.Close()

	if err := srv.Shutdown(ctx); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("server forced to shutdown")
	}

	zlog.Logger.Info().Msg("server exited")
	time.Sleep(5 * time.Second)
}
