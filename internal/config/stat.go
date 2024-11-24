package config

type Stat struct {
	LA       bool `envconfig:"STAT_LA" default:"true"`
	CPU      bool `envconfig:"STAT_CPU" default:"true"`
	DiskLoad bool `envconfig:"STAT_DISK_LOAD" default:"true"`
	DiskUse  bool `envconfig:"STAT_DISK_USE" default:"true"`
	NetStat  bool `envconfig:"STAT_NET_STAT" default:"true"`
	NetTop   bool `envconfig:"STAT_NET_TOP" default:"true"`
}
