package rpPostgres

import "github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.6/sales-tracker/sales-tracker_main-server/internal/pkg/pkgPostgres"

// RpPostgres implements the repository interface for PostgreSQL
type RpPostgres struct {
	db *pkgPostgres.Postgres
}

// NewPostgresRepository creates a new instance of PostgresRepository
func New(db *pkgPostgres.Postgres) *RpPostgres {
	return &RpPostgres{db: db}
}

// Close closes the database connection
func (rp *RpPostgres) CloseDB() error {
	return rp.db.Close()
}
