package main

import (
	"context"
	"microservice/internal/app/jarklin"
	"microservice/internal/config"
	infq "microservice/internal/infra/queue"
	"microservice/internal/stubs"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	zap.S().Info("bootstrapping microservice")

	ctx, cancel := context.WithCancel(context.Background())

	err := config.InitConfig()
	if err != nil {
		zap.S().Fatalw("failed to init config", "error", err)
	}
	zap.S().Info("config loaded")

	pr, err := infq.NewProducer(config.GetKafkaBrokers())
	if err != nil {
		zap.S().Fatalw("failed to create kafka producer", "error", err)
	}
	zap.S().Info("kafka producer created")

	jarklinConsumer, err := infq.NewConsumer(ctx, config.GetKafkaBrokers(), "stubs", stubs.NewJarklinConsumerHandler(pr))
	if err != nil {
		zap.S().Fatalw("failed to create Jarklin consumer", "error", err)
	}
	zap.S().Info("Jarklin consumer created")

	jarklinConsumer.Run(ctx, []string{jarklin.TopicJarklinEvents})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	<-sigs
	cancel()
	jarklinConsumer.Stop()
}
