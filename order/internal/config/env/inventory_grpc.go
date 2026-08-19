package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type internalGRPCEnvConfig struct {
	Host string `env:"INVENTORY_GRPC_HOST,required"`
	Port string `env:"INVENTORY_GRPC_PORT,required"`
}

type internalGRPCConfig struct {
	raw internalGRPCEnvConfig
}

func NewInternalGRPCConfig() (*internalGRPCConfig, error) {
	var raw internalGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &internalGRPCConfig{
		raw: raw,
	}, nil
}

func (cfg *internalGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
