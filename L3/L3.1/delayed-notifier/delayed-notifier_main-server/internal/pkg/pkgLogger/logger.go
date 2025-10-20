package pkgLogger

import (
	"os"

	"github.com/natefinch/lumberjack/v3"
	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/zlog"
)

type levelFilterWriter struct {
	Writer   zerolog.LevelWriter
	MinLevel zerolog.Level
}

func (w levelFilterWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level < w.MinLevel {
		return len(p), nil
	}
	return w.Writer.WriteLevel(level, p)
}

func (w levelFilterWriter) Write(p []byte) (n int, err error) {
	return w.WriteLevel(zerolog.InfoLevel, p)
}

func InitLogger(cfg *Config) error {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}

	consoleLevelWriter := levelFilterWriter{
		Writer:   zerolog.LevelWriterAdapter{Writer: consoleWriter},
		MinLevel: cfg.ConsoleLevel,
	}

	var fileLevelWriter levelFilterWriter
	if cfg.EnableFile {
		roller, err := lumberjack.NewRoller(
			cfg.FilePath,
			cfg.MaxSizeB,
			&lumberjack.Options{
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
				Compress:   cfg.Compress,
			},
		)
		if err != nil {
			return err
		}
		fileLevelWriter = levelFilterWriter{
			Writer:   zerolog.LevelWriterAdapter{Writer: roller},
			MinLevel: cfg.FileLevel,
		}
	}

	var multi zerolog.LevelWriter
	if cfg.EnableFile {
		multi = zerolog.MultiLevelWriter(consoleLevelWriter, fileLevelWriter)
	} else {
		multi = consoleLevelWriter
	}

	logger := zerolog.New(multi).With().Timestamp().Logger()
	zlog.Logger = logger

	zlog.Logger.Info().
		Str("file", cfg.FilePath).
		Str("level_console", cfg.ConsoleLevel.String()).
		Str("level_file", cfg.FileLevel.String()).
		Msg("logger initialized")

	return nil
}
