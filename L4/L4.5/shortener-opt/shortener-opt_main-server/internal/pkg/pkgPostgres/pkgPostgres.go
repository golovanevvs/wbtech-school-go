package pkgPostgres

import (
	"fmt"

	"github.com/wb-go/wbf/dbpg"
)

type Postgres struct {
	DB *dbpg.DB
}

func New(cfg *Config) (*Postgres, error) {
	masterDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Master.Host, cfg.Master.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	opts := &dbpg.Options{
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
	}

	db, err := dbpg.New(masterDSN, nil, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create DB instance: %w", err)
	}

	return &Postgres{
		DB: db,
	}, nil
}

func (p *Postgres) Close() error {
	if err := p.DB.Master.Close(); err != nil {
		return err
	}
	return nil
}
