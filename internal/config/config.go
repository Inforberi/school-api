package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      App
	Log      Log
	Postgres Postgres
	// Redis    Redis    `env-prefix:""`
}

type Log struct {
	Level  string `env:"LOG_LEVEL" env-default:"info"`   // debug, info, warn, error
	Format string `env:"LOG_FORMAT" env-default:"json"`  // console | json
}

type App struct {
	Env  string `env:"APP_ENV"` // local|stage|prod
	HTTP struct {
		Addr              string        `env:"HTTP_ADDR" env-default:":8080"`
		ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" env-default:"5s"`
		ReadTimeout       time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"15s"`
		WriteTimeout      time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"15s"`
		IdleTimeout       time.Duration `env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
	} `env-prefix:""`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Port     int    `env:"POSTGRES_PORT" env-required:"true"`
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName   string `env:"POSTGRES_NAME_DB" env-required:"true"`
	SSLMode  string `env:"PG_SSL_MODE" env-default:"disable"`

	// Пул (опционально)
	MaxConns        int32         `env:"PG_MAX_CONNS" env-default:"20"`
	MinConns        int32         `env:"PG_MIN_CONNS" env-default:"2"`
	HealthTimeout   time.Duration `env:"PG_HEALTH_TIMEOUT" env-default:"3s"`
	MaxConnLifetime time.Duration `env:"PG_MAX_CONN_LIFETIME" env-default:"30m"`
	MaxConnIdleTime time.Duration `env:"PG_MAX_CONN_IDLE_TIME" env-default:"5m"`
}

// type Redis struct {
// 	Addr          string        `env:"REDIS_ADDR" env-required:"true"` // host:port
// 	Password      string        `env:"REDIS_PASSWORD" env-default:""`
// 	DB            int           `env:"REDIS_DB" env-default:"0"`
// 	DialTimeout   time.Duration `env:"REDIS_DIAL_TIMEOUT" env-default:"2s"`
// 	ReadTimeout   time.Duration `env:"REDIS_READ_TIMEOUT" env-default:"2s"`
// 	WriteTimeout  time.Duration `env:"REDIS_WRITE_TIMEOUT" env-default:"2s"`
// 	HealthTimeout time.Duration `env:"REDIS_HEALTH_TIMEOUT" env-default:"2s"`
// 	KeyPrefix     string        `env:"REDIS_KEY_PREFIX" env-default:"app"`
// }

func Load() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		return nil, fmt.Errorf("config: read env: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config: validate: %w", err)
	}
	return &cfg, nil
}

func (c *Config) Validate() error {
	var errs []error

	if c.Postgres.Port <= 0 || c.Postgres.Port > 65535 {
		errs = append(errs, errors.New("POSTGRES_PORT out of range"))
	}
	// if c.Redis.Addr == "" {
	// 	errs = append(errs, errors.New("REDIS_ADDR is required"))
	// }

	return errors.Join(errs...)
}

func (p Postgres) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode,
	)
}
