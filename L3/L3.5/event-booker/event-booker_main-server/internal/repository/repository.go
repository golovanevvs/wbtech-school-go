package repository

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgPostgres"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/repository/rpPostgres"
)

type Repository struct {
	*rpPostgres.RpPostgres
	*rpPostgres.UserRepository
	*rpPostgres.EventRepository
	*rpPostgres.BookingRepository
	*rpPostgres.RefreshTokenRepository
	*rpPostgres.NoticeRepository
}

func New(pg *pkgPostgres.Postgres) (*Repository, error) {
	userRp := rpPostgres.NewUserRepository(pg)
	return &Repository{
		RpPostgres:             rpPostgres.New(pg),
		UserRepository:         userRp,
		EventRepository:        rpPostgres.NewEventRepository(pg),
		BookingRepository:      rpPostgres.NewBookingRepository(pg),
		RefreshTokenRepository: rpPostgres.NewRefreshTokenRepository(pg),
		NoticeRepository:       rpPostgres.NewNoticeRepository(userRp),
	}, nil

}
