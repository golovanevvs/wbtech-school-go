package rpPostgres

import "github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.3/calendar/calendar_main-server/internal/pkg/pkgPostgres"

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
