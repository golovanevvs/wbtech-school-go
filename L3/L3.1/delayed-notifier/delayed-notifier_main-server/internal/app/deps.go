package app

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/rabbitmq"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport"
	"github.com/wb-go/wbf/zlog"
)

type dependencies struct {
	tr *transport.Transport
	rp *repository.Repository
	sv *service.Service
	rb *rabbitmq.Client
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
		return fmt.Errorf("error set log level: %w", err)
	}

	zlog.Logger.Debug().
		Str("component", "app").
		Str("log_level", zlog.Logger.GetLevel().String()).
		Msg("logging level has been configured")

	return nil
}

func (b *dependencyBuilder) withRepository() error {
	rp, err := repository.New(b.cfg.rp)
	if err != nil {
		return fmt.Errorf("error initialize repository: %w", err)
	}
	zlog.Logger.Debug().Str("component", "app").Msg("repository has been initialized")
	b.deps.rp = rp

	return nil
}

func (b *dependencyBuilder) withService() {
	sv := service.New(b.deps.rp, b.deps.rb)
	zlog.Logger.Debug().Str("component", "app").Msg("service has been initialized")
	b.deps.sv = sv
}

func (b *dependencyBuilder) withTransport() {
	zlog.Logger.Debug().Str("component", "app").Msg("transport has been initialized")
	b.deps.tr = transport.New(b.cfg.tr, b.deps.sv)
}

func (b *dependencyBuilder) withRabbitMQ() error {
	rb, err := rabbitmq.NewClient(*b.cfg.rb)
	if err != nil {
		return fmt.Errorf("error initialize RabbitMQ client: %w", err)
	}
	zlog.Logger.Debug().Str("component", "app").Msg("RabbitMQ client has been initialized")
	b.deps.rb = rb

	return nil
}

func (b *dependencyBuilder) build() (*dependencies, error) {
	if err := b.withLogger(); err != nil {
		return nil, err
	}

	if err := b.withRabbitMQ(); err != nil {
		return nil, err
	}

	if err := b.withRepository(); err != nil {
		return nil, err
	}

	b.withService()

	b.withTransport()

	return b.deps, nil
}
