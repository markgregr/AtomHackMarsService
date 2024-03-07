package config

type App struct {
	ErrorLevel string `env:"MARS_ERROR_LEVEL" envDefault:"info"`

	API      API
	Database Database
	Minio    Minio
}
