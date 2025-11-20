package app

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/pkg/pkgConst"
	"github.com/wb-go/wbf/zlog"
)

type App struct {
	cfg  *Config
	lg   *zlog.Zerolog
	deps *dependencies
	rm   *resourceManager
}

func New(env string) (*App, error) {
	cfg, err := newConfig(env)
	if err != nil {
		return nil, fmt.Errorf("creating configuration: %w", err)
	}

	lg := zlog.Logger.With().Str("component", "app").Logger()

	lg.Info().Msg(cfg.String())

	lg.Info().Msgf("%s starting dependency initialization...", pkgConst.Starting)
	deps, rm, err := newDependencyBuilder(cfg, &lg).build()
	if err != nil {
		return nil, fmt.Errorf("error dependencies initialization: %w", err)
	}
	lg.Info().Msgf("%s dependencies have been initialized successfully", pkgConst.Finished)

	return &App{
		cfg:  cfg,
		lg:   &lg,
		deps: deps,
		rm:   rm,
	}, nil
}
