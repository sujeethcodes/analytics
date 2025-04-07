package connectors

import (
	"log"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func LoadEnv() {
	log.Println("Creating env load instance")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load env file. Err: %s", err)
	}
}

func LoadLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Could not create logging instance. Err: %s", err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
	// defer global()
	logger.Info("Created logging instance")

	return logger
}
