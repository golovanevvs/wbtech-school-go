package pkgLogger

import (
	"os"
	"regexp"

	"github.com/natefinch/lumberjack/v3"
	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/zlog"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// levelFilterWriter фильтрует записи по уровню и опционально очищает ANSI-коды
type levelFilterWriter struct {
	Writer    zerolog.LevelWriter
	MinLevel  zerolog.Level
	StripANSI bool
}

func (w levelFilterWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level < w.MinLevel {
		return len(p), nil
	}

	data := p
	if w.StripANSI {
		data = ansiRegexp.ReplaceAll(p, []byte{})
	}

	return w.Writer.WriteLevel(level, data)
}

// Write реализует io.Writer для совместимости
func (w levelFilterWriter) Write(p []byte) (n int, err error) {
	return w.WriteLevel(zerolog.InfoLevel, p)
}

func InitLogger(cfg *Config) error {
	// =========================
	// stdout writer с цветами и timestamp
	// =========================
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}

	consoleLevelWriter := levelFilterWriter{
		Writer:    zerolog.LevelWriterAdapter{Writer: consoleWriter},
		MinLevel:  cfg.ConsoleLevel,
		StripANSI: false, // цвет оставляем
	}

	// =========================
	// writer для файла с ротацией через lumberjack
	// =========================
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
			Writer:    zerolog.LevelWriterAdapter{Writer: roller},
			MinLevel:  cfg.FileLevel,
			StripANSI: true, // удаляем ANSI-коды
		}
	}

	// =========================
	// MultiLevelWriter объединяет stdout и файл
	// =========================
	var multi zerolog.LevelWriter
	if cfg.EnableFile {
		multi = zerolog.MultiLevelWriter(consoleLevelWriter, fileLevelWriter)
	} else {
		multi = consoleLevelWriter
	}

	// =========================
	// создаем logger
	// =========================
	logger := zerolog.New(multi).With().Timestamp().Logger()
	zlog.Logger = logger

	zlog.Logger.Info().
		Str("file", cfg.FilePath).
		Str("level_console", cfg.ConsoleLevel.String()).
		Str("level_file", cfg.FileLevel.String()).
		Msg("logger initialized")

	return nil
}
