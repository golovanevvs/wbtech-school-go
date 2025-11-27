package rpPostgres

import "github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.7/warehouse-control/warehouse-control_main-server/internal/pkg/pkgPostgres"

// RpPostgres implements the repository interface for PostgreSQL
type RpPostgres struct {
	db *pkgPostgres.Postgres
}

// NewPostgresRepository creates a new instance of PostgresRepository
func New(db *pkgPostgres.Postgres) *RpPostgres {
	return &RpPostgres{db: db}
}

// Close closes the database connection
func (rp *RpPostgres) Close() error {
	return rp.db.Close()
}

// User returns the user repository
func (rp *RpPostgres) User() *UserRepository {
	return NewUserRepository(rp.db)
}

// RefreshToken returns the refresh token repository
func (rp *RpPostgres) RefreshToken() *RefreshTokenRepository {
	return NewRefreshTokenRepository(rp.db)
}
