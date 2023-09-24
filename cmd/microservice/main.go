package main

import (
	"context"
	"microservice/internal/app/acaer"
	"microservice/internal/app/exibillia"
	"microservice/internal/app/jarklin"
	"microservice/internal/app/looncan"
	"microservice/internal/config"
	"microservice/internal/cron"
	"microservice/internal/database"
	infq "microservice/internal/infra/queue"
	"microservice/internal/queue"
	"microservice/internal/server"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	v1 "microservice/pkg/v1"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	zap.S().Info("bootstrapping microservice")

	ctx := context.Background()

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

	pr, err := infq.NewProducer(config.GetKafkaBrokers())
	if err != nil {
		zap.S().Fatalw("failed to create kafka producer", "error", err)
	}
	zap.S().Info("kafka producer created")

	exibilliaService := exibillia.NewService(exibillia.NewStorage(db))
	looncanService := looncan.NewService(looncan.NewStorage(db))
	acaerService := acaer.NewService(acaer.NewStorage(db), looncanService)
	jarklinService := jarklin.NewService(jarklin.NewStorage(db), jarklin.NewSender(pr))

	jarklinConsumer, err := infq.NewConsumer(ctx, config.GetKafkaBrokers(), config.GetKafkaConsumerGroupName(), queue.NewHandlerJarklin(jarklinService))
	if err != nil {
		zap.S().Fatalw("failed to create Jarklin consumer", "error", err)
	}
	zap.S().Info("Jarklin consumer created")

	srv := server.NewServer(exibilliaService, looncanService, acaerService, jarklinService)

	zap.S().Info("running app")

	jarklinConsumer.Run(ctx, []string{jarklin.TopicJarklinEvents})

	cr := cron.NewCron(config.GetCronJobs(), jarklinService)
	cr.Start()

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		zap.S().Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	v1.RegisterExibilliaServiceServer(s, srv)
	v1.RegisterLooncanServiceServer(s, srv)
	v1.RegisterAcaerServiceServer(s, srv)
	v1.RegisterJarklinServiceServer(s, srv)
	reflection.Register(s)

	zap.S().Info("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		zap.S().Fatalf("failed to serve: %v", err)
	}
}
