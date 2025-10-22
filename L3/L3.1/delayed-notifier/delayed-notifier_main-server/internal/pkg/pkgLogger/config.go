package pkgLogger

import (
	"fmt"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	Level string
	// EnableFile   bool
	// FilePath     string
	// MaxSizeB     int64
	// MaxBackups   int
	// MaxAge       time.Duration
	// Compress     bool
	// ConsoleLevel zerolog.Level
	// FileLevel    zerolog.Level
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{
		Level: cfg.GetString("logger.level"),
		// EnableFile:   cfg.GetBool("logger.enable_file"),
		// FilePath:     cfg.GetString("logger.file_path"),
		// MaxSizeB:     cfg.GetInt64("logger.max_size_b"),
		// MaxBackups:   cfg.GetInt("logger.max_backups"),
		// MaxAge:       cfg.GetDuration("logger.max_age"),
		// Compress:     cfg.GetBool("logger.compress"),
		// ConsoleLevel: parseLevel(cfg.GetString("logger.console_level")),
		// FileLevel:    parseLevel(cfg.GetString("logger.file_level")),
	}
}

// func parseLevel(level string) zerolog.Level {
// 	l, err := zerolog.ParseLevel(level)
// 	if err != nil {
// 		return zerolog.InfoLevel
// 	}
// 	return l
// }

//	func (c Config) String() string {
//		return fmt.Sprintf(`logger:
//	  %s: %s
//	  %s: %s
//	  %s: %s
//	  %s: %s
//	  %s: %s
//	  %s: %s
//	  %s: %s
//	  %s: %s`,
//			color.YellowString("enable file"), color.GreenString("%v", c.EnableFile),
//			color.YellowString("file path"), color.GreenString(c.FilePath),
//			color.YellowString("max size, B"), color.GreenString("%d", c.MaxSizeB),
//			color.YellowString("max backups"), color.GreenString("%d", c.MaxBackups),
//			color.YellowString("max age"), color.GreenString("%v", c.MaxAge),
//			color.YellowString("compress"), color.GreenString("%v", c.Compress),
//			color.YellowString("console level"), color.GreenString(c.ConsoleLevel.String()),
//			color.YellowString("file level"), color.GreenString(c.FileLevel.String()))
//	}
// func (c Config) String() string {
// 	return fmt.Sprintf(`logger:
//   %s: %v
//   %s: %s
//   %s: %d
//   %s: %d
//   %s: %v
//   %s: %v
//   %s: %s
//   %s: %s`,
// 		"enable file", c.EnableFile,
// 		"file path", c.FilePath,
// 		"max size, B", c.MaxSizeB,
// 		"max backups", c.MaxBackups,
// 		"max age", c.MaxAge,
// 		"compress", c.Compress,
// 		"console level", c.ConsoleLevel.String(),
// 		"file level", c.FileLevel.String())
// }

func (c Config) String() string {
	return fmt.Sprintf(`logger:
  %s: %s`,
		"level", c.Level)
}
