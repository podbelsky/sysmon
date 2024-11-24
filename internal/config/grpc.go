package config

import (
	"net"
	"strconv"
)

type GRPC struct {
	NetworkType string `envconfig:"GRPC_NETWORK_TYPE" default:"tcp"`
	Host        string `envconfig:"GRPC_HOST" default:"localhost"`
	Port        int32  `envconfig:"GRPC_PORT" default:"8081"`
}

func (c *Config) GRPCAddr() string {
	return net.JoinHostPort(c.GRPC.Host, strconv.Itoa(int(c.GRPC.Port)))
}
