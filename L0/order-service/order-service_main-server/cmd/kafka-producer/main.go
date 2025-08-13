package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service_main-server/internal/config"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service_main-server/internal/kafka"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service_main-server/internal/logger/zlog"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service_main-server/internal/model"
)

func main() {
	t := 5

	zlog.Init()

	zlog.Logger.Info().Msg("kafka-producer starting")

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
			zlog.Logger.Error().Err(err).Msg("failed to load default configuration")
			time.Sleep(time.Duration(t) * time.Second)
			os.Exit(1)
		}
	}
	zlog.Logger.Info().Msg("configuration loaded successfully")

	logLevelStr := cfg.Logger.LogLevel
	logLevel, err := zlog.ParseLogLevel(logLevelStr)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to parse log level")
		time.Sleep(time.Duration(t) * time.Second)
		os.Exit(1)
	}
	zlog.Logger.Info().Str("logLevel", logLevel.String()).Msg("logging level")
	zlog.Logger = zlog.Logger.Level(logLevel)

	k, err := kafka.New(&cfg.Kafka, &zlog.Logger)
	if err != nil {
		time.Sleep(time.Duration(t) * time.Second)
		os.Exit(1)
	}

	sp, err := kafka.NewSyncProducer(k)
	if err != nil {
		time.Sleep(time.Duration(t) * time.Second)
		os.Exit(1)
	}
	defer sp.Close()

	dc, err := time.Parse(time.RFC3339Nano, "2021-11-26T06:22:19Z")
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("error time parse")
		time.Sleep(time.Duration(t) * time.Second)
		os.Exit(1)
	}

	var ctx context.Context
	var cancel context.CancelFunc

	for i := range 100 {
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		var items []model.Item
		goodsTotal := 0
		rj := rand.Intn(30)
		for range rj {
			price := rand.Intn(10000)
			sale := rand.Intn(50)
			totalPrice := price - price*sale/100
			goodsTotal += totalPrice
			item := model.Item{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       price,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        sale,
				Size:        "0",
				TotalPrice:  totalPrice,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			}
			items = append(items, item)
		}

		m := model.Order{
			OrderUID:    fmt.Sprintf("b563feb7b2b84b6test_%d", i),
			TrackNumber: "WBILMTESTTRACK",
			Entry:       "WBIL",
			Delivery: model.Delivery{
				Name:    "Test Testov",
				Phone:   "+9720000000",
				Zip:     "2639809",
				City:    "Kiryat Mozkin",
				Address: "Ploshad Mira 15",
				Region:  "Kraiot",
				Email:   fmt.Sprintf("test_%d@gmail.com", i),
			},
			Payment: model.Payment{
				Transaction:  fmt.Sprintf("b563feb7b2b84b6test_%d", i),
				RequestID:    "",
				Currency:     "USD",
				Provider:     "wbpay",
				Amount:       1817,
				PaymentDt:    1637907727,
				Bank:         "alpha",
				DeliveryCost: 1500,
				GoodsTotal:   goodsTotal,
				CustomFee:    0,
			},
			Items:             items,
			Locale:            "en",
			InternalSignature: "",
			CustomerID:        "test",
			DeliveryService:   "meest",
			Shardkey:          "9",
			SmID:              99,
			DateCreated:       dc,
			OofShard:          "1",
		}

		mJSON, err := json.MarshalIndent(m, "", "\t")
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("error encoding json")
			time.Sleep(time.Duration(t) * time.Second)
			os.Exit(1)
		}

		zlog.Logger.Trace().Str("order_uid", m.OrderUID).Msg("sending")
		sp.SendSync(ctx, cfg.Kafka.Topic, nil, sarama.ByteEncoder(mJSON), nil)

		time.Sleep(time.Duration(t) * time.Second)
	}

	cancel()
}
