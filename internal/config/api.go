package config

type API struct {
	ServiceHost string `env:"MARS_SERVICE_HOST" envDefault:"0.0.0.0"`
	ServicePort int    `env:"MARS_SERVICE_PORT" envDefault:"8080"`
}
