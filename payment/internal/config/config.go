package config

import (
	"os"

	"github.com/Mas4trt/microservices/payment/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	PaymentGRPC PaymentGRPCConfig
	Logger      LoggerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		PaymentGRPC: paymentGRPCCfg,
		Logger:      loggerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
