package main

import (
	"context"
	"os"
	"time"

	"github.com/podbelsky/sysmon/cmd"
	"github.com/podbelsky/sysmon/internal/config"
	"github.com/rs/zerolog"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	logLevel, err := conf.LogLevel()
	if err != nil {
		panic(err)
	}

	log := zerolog.New(os.Stderr).
		Level(logLevel).
		Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly}).
		With().Timestamp(). //Caller().
		Str("service.name", conf.App.Name).
		Logger()

	ctx := log.WithContext(context.Background())

	log.Info().Msg("the application is launching")

	exitCode := 0

	err = cmd.Run(ctx, conf)
	if err != nil {
		log.Err(err).Send()

		exitCode = 1
	}

	os.Exit(exitCode)
}
