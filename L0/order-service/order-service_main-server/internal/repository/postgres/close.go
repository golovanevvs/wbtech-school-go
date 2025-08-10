package postgres

func (p *Postgres) Close() error {
	p.logger.Info().Msg("closing database connection")
	if err := p.db.Close(); err != nil {
		p.logger.Error().Err(err).Msg("failed to close database connection")
		return err
	}
	p.logger.Info().Msg("database connection closed")
	return nil
}
