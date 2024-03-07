package config

type Database struct {
	DSN         string `env:"MARS_POSTGRES_URL"`
	AutoMigrate bool   `env:"MARS_DATABASE_AUTO_MIGRATE" envDefault:"false"`
}
