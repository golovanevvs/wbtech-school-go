package repository

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service_main-server/internal/config"
	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service_main-server/internal/repository/postgres"
	"github.com/rs/zerolog"
)

type Repository struct {
	*postgres.Postgres
}

func New(ctx context.Context, cfgRp *config.Repository, logger *zerolog.Logger) (*Repository, error) {
	postgres, err := postgres.New(ctx, cfgRp, logger)
	if err != nil {
		return nil, err
	}

	return &Repository{
		postgres,
	}, nil
}
