package main

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"microservice/internal/app/exibillia"
	"microservice/internal/config"
	"microservice/internal/database"
	"microservice/internal/server"
	"net"

	v1 "microservice/pkg/v1"
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

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		zap.S().Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	v1.RegisterExibilliaServiceServer(s, srv)

	zap.S().Info("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		zap.S().Fatalf("failed to serve: %v", err)
	}
}
