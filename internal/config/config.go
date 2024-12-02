package config

import (
	"os"

	"github.com/cockroachdb/errors"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type Config struct {
	App  App
	GRPC GRPC
	Stat Stat
	Time Time
}

type App struct {
	ENV      string `envconfig:"APP_ENV"            default:"local"`
	Name     string `envconfig:"APP_NAME"           default:"sysmon"`
	LogLevel string `envconfig:"LOG_LEVEL"          default:"debug"`
}

func Load() (Config, error) {
	cnf := Config{} //nolint:exhaustruct

	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return cnf, errors.Wrap(err, "read .env file")
	}

	if err := envconfig.Process("", &cnf); err != nil {
		return cnf, errors.Wrap(err, "read environment")
	}

	return cnf, nil
}

func (c *Config) LogLevel() (zerolog.Level, error) {
	lvl, err := zerolog.ParseLevel(c.App.LogLevel)
	if err != nil {
		return 0, errors.Wrapf(err, "loading log level from config value %q", c.App.LogLevel)
	}

	return lvl, nil
}
