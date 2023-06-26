package main

import (
	"go.uber.org/zap"
	"microservice/internal/app/exibillia"
	"microservice/internal/config"
	"microservice/internal/database"
	"microservice/internal/server"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	zap.S().Info("bootstrapping microservice")

	err := config.InitConfig()
	if err != nil {
		zap.S().Fatalw("failed to init config", "error", err)
	}
	zap.S().Info("config loaded")

	db, err := database.NewConnection(config.GetDB())
	if err != nil {
		zap.S().Fatalw("failed to connect to database", "error", err)
	}
	zap.S().Info("db connected")

	exibilliaService := exibillia.NewService(exibillia.NewStorage(db))

	srv := server.NewServer(exibilliaService)

	zap.S().Info("running app")

	zap.S().Error(srv.Listen())
}
