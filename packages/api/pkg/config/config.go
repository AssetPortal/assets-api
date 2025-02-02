package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/subosito/gotenv"
)

type Configuration struct {
	ServiceAddress        string                `env:"SERVICE_ADDRESS" envDefault:":8000"`
	HTTPTimeout           time.Duration         `env:"HTTP_TIMEOUT" envDefault:"20s"`
	MaxRequestsPerSecond  int                   `env:"MAX_REQUESTS_PER_SECOND" envDefault:"3"`
	LogLevel              string                `env:"LOG_LEVEL" envDefault:"warn"`
	TokenExpiration       time.Duration         `env:"TOKEN_EXPIRATION" envDefault:"5m"`
	AuthConfiguration     AuthConfiguration     `envPrefix:"AUTH_"`
	BucketConfiguration   BucketConfiguration   `envPrefix:"BUCKET_"`
	DatabaseConfiguration DatabaseConfiguration `envPrefix:"DATABASE_"`
}
type DatabaseConfiguration struct {
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS" envDefault:"10"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS" envDefault:"5"`
	ConnMaxIdleTime time.Duration `env:"CONN_MAX_IDLE_TIME" envDefault:"5m"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME" envDefault:"30m"`
	URL             string        `env:"URL"`
}
type AuthConfiguration struct {
	APIURL      string        `env:"API_URL"`
	HTTPTimeout time.Duration `env:"HTTP_TIMEOUT" envDefault:"20s"`
	Enabled     bool          `env:"ENABLED" envDefault:"true"`
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
