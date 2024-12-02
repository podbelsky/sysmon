package config

import (
	"time"
)

type Time struct {
	Snap  time.Duration `envconfig:"SNAP_TIME" default:"5s"`
	Clean time.Duration `envconfig:"CLEAN_TIME" default:"1m"`
	Store time.Duration `envconfig:"STORE_TIME" default:"1h"`
}
