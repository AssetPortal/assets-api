package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/subosito/gotenv"
)

type Configuration struct {
	ServiceAddress       string              `env:"SERVICE_ADDRESS" envDefault:":8000"`
	HTTPTimeout          time.Duration       `env:"HTTP_TIMEOUT" envDefault:"20s"`
	RWDBURL              string              `env:"RW_DB_URL" envDefault:"10"`
	MaxRequestsPerSecond int                 `env:"MAX_REQUESTS_PER_SECOND" envDefault:"10"`
	LogLevel             string              `env:"LOG_LEVEL" envDefault:"warn"`
	TokenExpiration      time.Duration       `env:"TOKEN_EXPIRATION" envDefault:"5m"`
	AuthConfiguration    AuthConfiguration   `envPrefix:"AUTH_"`
	BucketConfiguration  BucketConfiguration `envPrefix:"BUCKET_"`
}

type AuthConfiguration struct {
	APIURL      string        `env:"API_URL"`
	HTTPTimeout time.Duration `env:"HTTP_TIMEOUT" envDefault:"20s"`
}

type BucketConfiguration struct {
	AccessKey string `env:"ACCESS_KEY" envDefault:"test"`
	SecretKey string `env:"SECRET_KEY" envDefault:"test"`
	Session   string `env:"SESSION"`
	Region    string `env:"REGION"`
	Name      string `env:"NAME"`
}

func MustGetConfig() (*Configuration, error) {
	_ = gotenv.Load()
	cfg := &Configuration{}
	err := env.Parse(cfg)
	return cfg, err
}
