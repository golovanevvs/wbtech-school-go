package app

import (
	"fmt"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgRedis"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/service"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.5/shortener-opt/shortener-opt_main-server/internal/transport"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

type dependencies struct {
	rs *pkgRetry.Retry
	rd *pkgRedis.Client
	pg *pkgPostgres.Postgres
	rp *repository.Repository
	sv *service.Service
	tr *transport.Transport
}

type dependencyBuilder struct {
	cfg  *Config
	lg   *zlog.Zerolog
	rm   *resourceManager
	deps *dependencies
}

func newDependencyBuilder(cfg *Config, lg *zlog.Zerolog) *dependencyBuilder {
	return &dependencyBuilder{
		cfg:  cfg,
		lg:   lg,
		rm:   &resourceManager{},
		deps: &dependencies{},
	}
}

func (b *dependencyBuilder) Close() error {
	return b.rm.closeAll()
}

func (b *dependencyBuilder) initLogger() error {
	err := zlog.SetLevel(b.cfg.lg.Level)

	if err != nil {
		return fmt.Errorf("set log level: %w", err)
	}

	b.lg = &zlog.Logger

	b.lg.Debug().
		Str("log_level", zlog.Logger.GetLevel().String()).
		Msgf("%s logging level has been configured", pkgConst.Info)

	return nil
}

func (b *dependencyBuilder) InitRetry() {
	b.deps.rs = pkgRetry.New(b.cfg.rs)
}

func (b *dependencyBuilder) initRedis() error {
	fn := func() error {
		rd, err := pkgRedis.New(b.cfg.rd)
		if err != nil {
			b.lg.Warn().Err(err).Int("port", b.cfg.rd.Port).Msgf("%s failed to initialize Redis", pkgConst.Warn)
			return err
		}
		b.deps.rd = rd
		return nil
	}
	if err := retry.Do(fn, retry.Strategy(*b.deps.rs)); err != nil {
		return fmt.Errorf("initialize Redis, port: %d: %w", b.cfg.rd.Port, err)
	}

	b.lg.Debug().Msgf("%s Redis has been initialized", pkgConst.Info)
	b.rm.addResource(resource{
		name:      "Redis client",
		closeFunc: func() error { return b.deps.rd.Close() },
	})
	return nil
}

func (b *dependencyBuilder) InitPostgres() error {
	fn := func() error {
		pg, err := pkgPostgres.New(b.cfg.pg)
		if err != nil {
			b.lg.Warn().Err(err).Int("port", b.cfg.pg.Master.Port).Msgf("%s failed to initialize Postgres", pkgConst.Warn)
			return err
		}
		b.deps.pg = pg
		return nil
	}
	if err := retry.Do(fn, retry.Strategy(*b.deps.rs)); err != nil {
		return fmt.Errorf("initialize Postgres, port: %d: %w", b.cfg.pg.Master.Port, err)
	}

	b.lg.Debug().Msgf("%s Postgres has been initialized", pkgConst.Info)
	b.rm.addResource(resource{
		name:      "Postgres",
		closeFunc: func() error { return b.deps.pg.Close() },
	})
	return nil
}

func (b *dependencyBuilder) initRepository() error {
	rp, err := repository.New(b.deps.pg, b.deps.rs)
	if err != nil {
		return fmt.Errorf("initialize repository: %w", err)
	}
	b.lg.Debug().Msgf("%s repository has been initialized", pkgConst.Info)
	b.deps.rp = rp
	return nil
}

func (b *dependencyBuilder) initService() {
	sv := service.New(b.deps.rs, b.deps.rp)
	b.lg.Debug().Msgf("%s service has been initialized", pkgConst.Info)
	b.deps.sv = sv
}

func (b *dependencyBuilder) initTransport() {
	b.lg.Debug().Msgf("%s transport has been initialized", pkgConst.Info)
	b.deps.tr = transport.New(b.cfg.tr, b.deps.rs, b.deps.sv)
}

func (b *dependencyBuilder) build() (*dependencies, *resourceManager, error) {
	if err := b.initLogger(); err != nil {
		return nil, b.rm, err
	}
	b.InitRetry()
	if err := b.initRedis(); err != nil {
		return nil, b.rm, err
	}
	if err := b.InitPostgres(); err != nil {
		return nil, b.rm, err
	}
	if err := b.initRepository(); err != nil {
		return nil, b.rm, err
	}
	b.initService()
	b.initTransport()
	return b.deps, b.rm, nil
}
