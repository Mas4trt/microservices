package config

import "time"

type OrderHTTPConfig interface {
	Address() string
	ReadTimeout() time.Duration
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
	MigrationDir() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type PaymentGRPCConfig interface {
	Address() string
}
