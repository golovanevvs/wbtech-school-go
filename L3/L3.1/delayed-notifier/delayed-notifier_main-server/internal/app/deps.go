package app

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgEmail"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRabbitmq"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRedis"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgTelegram"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/service"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/transport"
	"github.com/wb-go/wbf/zlog"
)

type dependencies struct {
	rd *pkgRedis.Client
	rb *pkgRabbitmq.Client
	tg *pkgTelegram.Client
	em *pkgEmail.Client
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

func (b *dependencyBuilder) initLogger() error {
	// if err := pkgLogger.InitLogger(b.cfg.lg); err != nil {
	// 	return fmt.Errorf("error initialize logger: %w", err)
	// }

	// zlog.Logger.Debug().
	// 	Str("component", "app").
	// 	Str("log_level", b.cfg.lg.ConsoleLevel.String()).
	// 	Msg("logging level has been configured")

	// return nil

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

func (b *dependencyBuilder) initRedis() error {
	rd, err := pkgRedis.New(b.cfg.rd)
	if err != nil {
		return fmt.Errorf("error initialize Redis: %w", err)
	}
	zlog.Logger.Debug().Str("component", "app").Msg("Redis has been initialized")
	b.deps.rd = rd
	b.rm.addResource(resource{
		name:      "Redis client",
		closeFunc: func() error { return b.deps.rd.Close() },
	})
	return nil
}

func (b *dependencyBuilder) initRabbitMQ() error {
	rb, err := pkgRabbitmq.NewClient(b.cfg.rb)
	if err != nil {
		return fmt.Errorf("error initialize RabbitMQ client: %w", err)
	}
	zlog.Logger.Debug().Str("component", "app").Msg("RabbitMQ client has been initialized")
	b.deps.rb = rb
	b.rm.addResource(resource{
		name:      "RabbitMQ client",
		closeFunc: func() error { return b.deps.rb.Close() },
	})
	return nil
}

func (b *dependencyBuilder) initTelegram() error {
	tg, err := pkgTelegram.New(b.cfg.tg)
	if err != nil {
		return fmt.Errorf("error initialize telegram client: %w", err)
	}
	zlog.Logger.Debug().Str("component", "app").Msg("telegram client has been initialized")
	b.deps.tg = tg
	return nil
}

func (b *dependencyBuilder) initEmail() error {
	em, err := pkgEmail.New(b.cfg.em)
	if err != nil {
		return fmt.Errorf("error initialize email client: %w", err)
	}
	zlog.Logger.Debug().Str("component", "app").Msg("email client has been initialized")
	b.deps.em = em
	return nil
}

func (b *dependencyBuilder) initRepository() error {
	rp, err := repository.New(b.deps.rd)
	if err != nil {
		return fmt.Errorf("error initialize repository: %w", err)
	}
	zlog.Logger.Debug().Str("component", "app").Msg("repository has been initialized")
	b.deps.rp = rp
	return nil
}

func (b *dependencyBuilder) initService() {
	sv := service.New(b.cfg.sv, b.deps.rp, b.deps.rb, b.deps.tg, b.deps.em)
	zlog.Logger.Debug().Str("component", "app").Msg("service has been initialized")
	b.deps.sv = sv
}

func (b *dependencyBuilder) initTransport() {
	zlog.Logger.Debug().Str("component", "app").Msgf("%s transport has been initialized", pkgConst.Info)
	b.deps.tr = transport.New(b.cfg.tr, b.deps.sv)
}

func (b *dependencyBuilder) build() (*dependencies, *resourceManager, error) {
	if err := b.initLogger(); err != nil {
		return nil, b.rm, err
	}
	if err := b.initRedis(); err != nil {
		return nil, b.rm, err
	}
	if err := b.initRabbitMQ(); err != nil {
		return nil, b.rm, err
	}
	if err := b.initTelegram(); err != nil {
		return nil, b.rm, err
	}
	if err := b.initEmail(); err != nil {
		return nil, b.rm, err
	}
	if err := b.initRepository(); err != nil {
		return nil, b.rm, err
	}
	b.initService()
	b.initTransport()
	return b.deps, b.rm, nil
}
