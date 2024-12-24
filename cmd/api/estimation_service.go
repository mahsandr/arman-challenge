package main

import (
	"context"
	"net"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"

	pb "github.com/mahsandr/arman-challenge/api/proto/generated/estimation/v1"
	"github.com/mahsandr/arman-challenge/config"
	estimationgrpc "github.com/mahsandr/arman-challenge/internal/delivery/grpc/estimation"
	"github.com/mahsandr/arman-challenge/internal/infrastructure/clickhouse"
	"github.com/mahsandr/arman-challenge/internal/infrastructure/kafka"
	"github.com/mahsandr/arman-challenge/internal/usecases/estimation"
	"github.com/vahid-sohrabloo/chconn/v2/chpool"
)

func RunEstimationService() {
	logger, _ := zap.NewProduction()

	kafkaEnvs, err := config.GetKafkaConfig()
	checkErr(logger, "error on getting kafka config", err)

	// Initialize repositories
	conn, err := chpool.New(serviceConfig.ClickHouseAddr)
	checkErr(logger, "error on connect to ClickHouse", err)
	clickhouseRepo, err := clickhouse.NewClickHouseRepository(conn, serviceConfig.UserSegmentTableName, serviceConfig.SegmentsViewName)
	checkErr(logger, "error on init clickhouse", err)
	// Initialize consumer
	kafkaConsumer, err := kafka.NewKafkaConsumer(logger, kafkaEnvs.Brokers, kafkaEnvs.Topic,
		kafkaEnvs.GroupID, kafkaEnvs.MinBytes, kafkaEnvs.MaxBytes, kafkaEnvs.PollInterval, kafkaEnvs.ConsumerBuffer)
	checkErr(logger, "error on init kafka", err)
	consumer := estimation.NewSegmentConsumer(clickhouseRepo, kafkaConsumer, logger)
	ctx, cancelFunc := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		if err := consumer.StartConsuming(ctx); err != nil {
			logger.Error("error on start consuming", zap.Error(err))
			cancelFunc()
		}
	}(ctx)
	// Initialize services
	estimationService := estimation.NewService(
		clickhouseRepo,
		logger,
	)

	// Initialize gRPC server
	grpcHandler := estimationgrpc.NewHandler(logger, estimationService)
	server := grpcServer(logger)
	pb.RegisterEstimationServiceServer(server, grpcHandler)

	// reflection service on gRPC server.
	reflection.Register(server)

	lis, err := net.Listen("tcp", serviceConfig.Host+":"+serviceConfig.Port)
	checkErr(logger, "failed to listen", err)

	go func() {
		logger.Info("Server running ",
			zap.String("host", serviceConfig.Host),
			zap.String("port", serviceConfig.Port),
		)
		err := server.Serve(lis)
		checkErr(logger, "failed to serve", err)
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	select {
	case <-ctx.Done():
	case <-ch:
		cancelFunc()
	}

	consumer.Stop()
	logger.Info("Stopping the server")
	server.Stop()
	logger.Info("Closing the listener")
	lis.Close()
	logger.Info("Server Shutdown")
}
