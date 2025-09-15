package app

import (
	"fmt"
	"os"
	"time"

	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/config"
	"github.com/golovanevvs/wbtech-school-go/L2/L2.18/cal/internal/logger/zlog"
)

func Run() {
	zlog.Init()
	zlog.Logger.Info().Msg("starting cal...")

	envFile := ".env"
	envPath := fmt.Sprintf("./%s", envFile)

	cfg := config.New()
	err := cfg.Load(envPath)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to load configuration")
		time.Sleep(time.Duration(5) * time.Second)
		os.Exit(1)
	}
	zlog.Logger.Info().Msg("configuration loaded successfully")

}
