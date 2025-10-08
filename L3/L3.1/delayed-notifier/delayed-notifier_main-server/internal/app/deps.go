package app

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport"
	"github.com/wb-go/wbf/zlog"
)

type dependencies struct {
	tr *transport.Transport
	rp *repository.Repository
	sv *service.Service
}

type dependencyBuilder struct {
	cfg  *appConfig
	rm   *resourceManager
	deps *dependencies
}

func newDependencyBuilder(cfg *appConfig) *dependencyBuilder {
	return &dependencyBuilder{
		cfg:  cfg,
		rm:   &resourceManager{},
		deps: &dependencies{},
	}
}

func (b *dependencyBuilder) withLogger() error {
	err := zlog.SetLevel(b.cfg.lg.Level)

	if err != nil {
		zlog.Logger.Error().Err(err).Str("component", "logger").Msg("error set log level")
		return fmt.Errorf("error set log level: %w", err)
	}

	zlog.Logger.Info().
		Str("component", "logger").
		Str("log_level", zlog.Logger.GetLevel().String()).
		Msg("logging level has been configure")

	return nil
}

func (b *dependencyBuilder) withRepository() error {
	rp, err := repository.New(b.cfg.rp)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("component", "repository").Msg("error create repository")
		return fmt.Errorf("error create repository: %w", err)
	}
	zlog.Logger.Info().Str("component", "repository").Msg("repository has been create")
	b.deps.rp = rp

	return nil
}

func (b *dependencyBuilder) withService() {
	sv := service.New(b.deps.rp)
	b.deps.sv = sv
}

func (b *dependencyBuilder) withTransport() {
	b.deps.tr = transport.New(b.cfg.tr)
}

func (b *dependencyBuilder) build() (*dependencies, error) {
	if err := b.withLogger(); err != nil {
		return nil, err
	}

	b.withTransport()

	if err := b.withRepository(); err != nil {
		return nil, err
	}

	b.withService()

	return b.deps, nil
}
