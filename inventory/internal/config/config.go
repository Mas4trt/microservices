package config

import (
	"os"

	"github.com/Mas4trt/microservices/inventory/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	InternalGRPC InventoryGRPCConfig
	Logger       LoggerConfig
	Mongo        MongoConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	internalGRPCCfg, err := env.NewInternalGRPCConfig()
	if err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		InternalGRPC: internalGRPCCfg,
		Logger:       loggerCfg,
		Mongo:        mongoCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
