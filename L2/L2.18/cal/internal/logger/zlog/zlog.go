package zlog

import (
	"os"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func Init() {
	Logger = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Level(zerolog.TraceLevel).
		Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05",
		})
}

func SetLevel(logLevelStr string) error {
	logLevel, err := zerolog.ParseLevel(logLevelStr)
	if err != nil {
		return err
	}

	Logger.Info().Str("logLevel", logLevel.String()).Msg("logging level")
	Logger = Logger.Level(logLevel)

	return nil
}
