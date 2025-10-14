package app

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/rabbitmq"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/telegram"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport"
	"github.com/wb-go/wbf/zlog"
)

type dependencies struct {
	rd *pkgRedis.Client
	rb *rabbitmq.Client
	tg *telegram.Client
	rp *repository.Repository
	sv *service.Service
	tr *transport.Transport
}

type dependencyBuilder struct {
	cfg  *Config
	rm   *resourceManager
	deps *dependencies
}

func newDependencyBuilder(cfg *Config) *dependencyBuilder {
	return &dependencyBuilder{
		cfg:  cfg,
		rm:   &resourceManager{},
		deps: &dependencies{},
	}
}

func (b *dependencyBuilder) Close() error {
	return b.rm.closeAll()
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

func (b *dependencyBuilder) withRedis() error {
	rd, err := pkgRedis.New(b.cfg.rd)
	if err != nil {
		return fmt.Errorf("error initialize Redis: %w", err)
	}
	zlog.Logger.Debug().Str("component", "app").Msg("Redis has been initialized")
	b.deps.rd = rd
	b.rm.addResource(func() error { return b.deps.rd.Close() })

	return nil
}

func (b *dependencyBuilder) withRabbitMQ() error {
	rb, err := rabbitmq.NewClient(b.cfg.rb)
	if err != nil {
		return fmt.Errorf("error initialize RabbitMQ client: %w", err)
	}
	zlog.Logger.Debug().Str("component", "app").Msg("RabbitMQ client has been initialized")
	b.deps.rb = rb
	b.rm.addResource(func() error { return b.deps.rb.Close() })

	return nil
}

func (b *dependencyBuilder) WithTelegram() error {
	tg, err := telegram.New(b.cfg.tg)
	if err != nil {
		return fmt.Errorf("error initialize telegram client: %w", err)
	}

	zlog.Logger.Debug().Str("component", "app").Msg("telegram client has been initialized")
	b.deps.tg = tg

	return nil
}

func (b *dependencyBuilder) withRepository() error {
	rp, err := repository.New(b.cfg.rp, b.deps.rd)
	if err != nil {
		return fmt.Errorf("error initialize repository: %w", err)
	}
	zlog.Logger.Debug().Str("component", "app").Msg("repository has been initialized")
	b.deps.rp = rp

	return nil
}

func (b *dependencyBuilder) withService() {
	sv := service.New(b.deps.rp, b.deps.rb, b.deps.tg, b.deps.rd)
	zlog.Logger.Debug().Str("component", "app").Msg("service has been initialized")
	b.deps.sv = sv
}

func (b *dependencyBuilder) withTransport() {
	zlog.Logger.Debug().Str("component", "app").Msg("transport has been initialized")
	b.deps.tr = transport.New(b.cfg.tr, b.deps.sv)
}

func (b *dependencyBuilder) build() (*dependencies, *resourceManager, error) {
	if err := b.withLogger(); err != nil {
		return nil, b.rm, err
	}
	if err := b.withRedis(); err != nil {
		return nil, b.rm, err
	}
	if err := b.withRabbitMQ(); err != nil {
		return nil, b.rm, err
	}
	if err := b.WithTelegram(); err != nil {
		return nil, b.rm, err
	}
	if err := b.withRepository(); err != nil {
		return nil, b.rm, err
	}
	b.withService()
	b.withTransport()
	return b.deps, b.rm, nil
}
