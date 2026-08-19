package config

import (
	"os"

	"github.com/Mas4trt/microservices/order/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	OrderHTTP     OrderHTTPConfig
	Logger        LoggerConfig
	Postgres      PostgresConfig
	InventoryGRPC InventoryGRPCConfig
	PaymentGRPC   PaymentGRPCConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	orderHHTPCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	inventoryGRPCCfg, err := env.NewInternalGRPCConfig()
	if err != nil {
		return err
	}

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		OrderHTTP:     orderHHTPCfg,
		Logger:        loggerCfg,
		Postgres:      postgresCfg,
		InventoryGRPC: inventoryGRPCCfg,
		PaymentGRPC:   paymentGRPCCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
