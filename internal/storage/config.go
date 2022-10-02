package storage

// Config is the configuration for storage storage
type Config struct {
	Db DBConfig
}

// DefaultConfig returns a default configuration for the db.
func DefaultConfig() Config {
	return Config{
		Db: DBConfig{
			Addr:           "localhost:27017",
			MaxConnections: 10,
			DBName:         "reporting",
		},
	}
}
