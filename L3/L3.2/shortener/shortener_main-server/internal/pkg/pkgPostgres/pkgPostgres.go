package pkgPostgres

// func New(rd *pkgRedis.Client) (*Repository, error) {
// masterDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 	cfg.Postgres.Master.Host, cfg.Postgres.Master.Port, cfg.Postgres.Master.User, cfg.Postgres.Master.Password, cfg.Postgres.Master.DBName)

// slave1DSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 	cfg.Postgres.Slave1.Host, cfg.Postgres.Slave1.Port, cfg.Postgres.Slave1.User, cfg.Postgres.Slave1.Password, cfg.Postgres.Slave1.DBName)

// slave2DSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 	cfg.Postgres.Slave2.Host, cfg.Postgres.Slave2.Port, cfg.Postgres.Slave2.User, cfg.Postgres.Slave2.Password, cfg.Postgres.Slave2.DBName)

// opts := &dbpg.Options{
// 	MaxOpenConns:    cfg.Postgres.MaxOpenConns,
// 	MaxIdleConns:    cfg.Postgres.MaxIdleConns,
// 	ConnMaxLifetime: cfg.Postgres.ConnMaxLifetime,
// }

// db, err := dbpg.New(masterDSN, []string{slave1DSN, slave2DSN}, opts)
// if err != nil {
// 	return nil, fmt.Errorf("failed to create db: %w", err)
// }

// return &Repository{
// 	Postgres: db,
// }, nil
// return &Repository{
// postgres: postgres.New(),
// }, nil
// }
