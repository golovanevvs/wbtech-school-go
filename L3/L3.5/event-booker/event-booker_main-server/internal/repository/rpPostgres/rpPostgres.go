package rpPostgres

import "github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.5/event-booker/event-booker_main-server/internal/pkg/pkgPostgres"

// RpPostgres implements the repository interface for PostgreSQL
type RpPostgres struct {
	db *pkgPostgres.Postgres
}

// NewPostgresRepository creates a new instance of PostgresRepository
func New(db *pkgPostgres.Postgres) *RpPostgres {
	return &RpPostgres{db: db}
}

// Close closes the database connection
func (r *RpPostgres) Close() error {
	return r.db.Close()
}

// User returns the user repository
func (r *RpPostgres) User() *UserRepository {
	return NewUserRepository(r.db)
}

// Event returns the event repository
func (r *RpPostgres) Event() *EventRepository {
	return NewEventRepository(r.db)
}

// Booking returns the booking repository
func (r *RpPostgres) Booking() *BookingRepository {
	return NewBookingRepository(r.db)
}

// RefreshToken returns the refresh token repository
func (r *RpPostgres) RefreshToken() *RefreshTokenRepository {
	return NewRefreshTokenRepository(r.db)
}
