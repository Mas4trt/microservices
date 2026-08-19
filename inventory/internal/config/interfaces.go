package config

type InventoryGRPCConfig interface {
	Address() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}
