package postgres

type Config struct {
	Master *DsnConfig
	Slave1 *DsnConfig
	Slave2 *DsnConfig
}

type DsnConfig struct {
	Port     int
	Host     string
	User     string
	Password string
}

func NewConfig() *Config {
	return &Config{
		Master: &DsnConfig{},
		Slave1: &DsnConfig{},
		Slave2: &DsnConfig{},
	}
}
