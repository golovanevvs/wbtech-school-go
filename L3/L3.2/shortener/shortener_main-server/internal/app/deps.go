package app

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgEmail"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgRabbitmq"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgRedis"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgTelegram"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/repository"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/service"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/transport"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

type dependencies struct {
	rs *pkgRetry.Retry
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
		return pkgErrors.Wrap(err, "set log level")
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
	var rd *pkgRedis.Client
	var err error
	fn := func() error {
		rd, err = pkgRedis.New(b.cfg.rd)
		if err != nil {
			b.lg.Warn().Err(err).Int("port", b.cfg.rd.Port).Msgf("%s failed to initialize Redis", pkgConst.Warn)
			return err
		}
		return nil
	}
	if err := retry.Do(fn, retry.Strategy(*b.deps.rs)); err != nil {
		return pkgErrors.Wrapf(err, "initialize Redis, port: %d, attempts: %d", b.cfg.rd.Port, b.cfg.rs.Attempts)
	}

	b.lg.Debug().Msgf("%s Redis has been initialized", pkgConst.Info)
	b.deps.rd = rd
	b.rm.addResource(resource{
		name:      "Redis client",
		closeFunc: func() error { return b.deps.rd.Close() },
	})
	return nil
}

func (b *dependencyBuilder) initRepository() error {
	rp, err := repository.New(b.deps.rd)
	if err != nil {
		return pkgErrors.Wrap(err, "initialize repository")
	}
	b.lg.Debug().Msgf("%s repository has been initialized", pkgConst.Info)
	b.deps.rp = rp
	return nil
}

func (b *dependencyBuilder) initService() {
	sv := service.New(b.deps.rs, b.deps.rp, b.deps.rb, b.deps.tg, b.deps.em)
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
	if err := b.initRepository(); err != nil {
		return nil, b.rm, err
	}
	b.initService()
	b.initTransport()
	return b.deps, b.rm, nil
}
