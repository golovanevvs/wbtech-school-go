package repository

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/repository/rpPostgres"
)

type Repository struct {
	*rpPostgres.RpPostgres
	*rpPostgres.UserRepository
	*rpPostgres.RefreshTokenRepository
	*rpPostgres.ItemRepository
	*rpPostgres.ItemHistoryRepository
}

func New(pg *pkgPostgres.Postgres) (*Repository, error) {
	userRp := rpPostgres.NewUserRepository(pg)
	refreshTokenRp := rpPostgres.NewRefreshTokenRepository(pg)
	itemRp := rpPostgres.NewItemRepository(pg)
	itemHistoryRp := rpPostgres.NewItemHistoryRepository(pg)

	return &Repository{
		RpPostgres:             rpPostgres.New(pg),
		UserRepository:         userRp,
		RefreshTokenRepository: refreshTokenRp,
		ItemRepository:         itemRp,
		ItemHistoryRepository:  itemHistoryRp,
	}, nil

}
