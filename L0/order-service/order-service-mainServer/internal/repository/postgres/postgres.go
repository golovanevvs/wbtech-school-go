package postgres

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type Postgres struct {
	db     *sqlx.DB
	logger *zerolog.Logger
}

func New(ctx context.Context, cfgRp *config.Repository, logger *zerolog.Logger) (*Postgres, error) {
	log := logger.With().Str("component", "postgres").Logger()

	log.Info().Msg("сonnecting to PostgreSQL")

	db, err := sqlx.ConnectContext(ctx, "pgx", cfgRp.Postgres.DatabaseDSN)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to PostgreSQL")
		return nil, err
	}

	db.SetMaxOpenConns(cfgRp.Postgres.MaxOpenConns)
	db.SetMaxIdleConns(cfgRp.Postgres.MaxIdleConns)
	db.SetConnMaxLifetime(cfgRp.Postgres.ConnMaxLifetime)

	log.Info().
		Int("max_open_conns", cfgRp.Postgres.MaxOpenConns).
		Int("max_idle_conns", cfgRp.Postgres.MaxIdleConns).
		Dur("conn_max_lifetime", cfgRp.Postgres.ConnMaxLifetime).
		Msg("database connection pool configured")

	log.Info().Msg("successfully connected to PostgreSQL")

	return &Postgres{
		db:     db,
		logger: &log,
	}, nil
}
