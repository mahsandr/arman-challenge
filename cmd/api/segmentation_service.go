package main

import (
	"net"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"

	pb "github.com/mahsandr/arman-challenge/api/proto/generated/segmentation/v1"
	"github.com/mahsandr/arman-challenge/config"
	segmentationgrpc "github.com/mahsandr/arman-challenge/internal/delivery/grpc/segmentation"
	"github.com/mahsandr/arman-challenge/internal/infrastructure/clickhouse"
	"github.com/mahsandr/arman-challenge/internal/infrastructure/kafka"
	"github.com/mahsandr/arman-challenge/internal/usecases/segmentation"
	"github.com/vahid-sohrabloo/chconn/v2/chpool"
)

func RunSegmentationService() {
	logger, _ := zap.NewProduction()

	kafkaEnvs, err := config.GetKafkaConfig()
	checkErr(logger, "error on getting kafka config", err)

	// Initialize repositories
	kafkaProducer, err := kafka.NewKafkaProducer(logger, kafkaEnvs.Brokers, kafkaEnvs.Topic, kafkaEnvs.BatchBytes, kafkaEnvs.FlushInterval)
	checkErr(logger, "error on init kafka", err)

	conn, err := chpool.New(serviceConfig.ClickHouseAddr)
	checkErr(logger, "error on connect to ClickHouse", err)

	clickhouseRepo, err := clickhouse.NewClickHouseRepository(conn, serviceConfig.UserSegmentTableName, serviceConfig.SegmentsViewName)
	checkErr(logger, "error on init clickhouse", err)
	// Initialize services
	segmentService := segmentation.NewService(
		clickhouseRepo,
		kafkaProducer,
		logger,
	)

	// Initialize gRPC server
	grpcHandler := segmentationgrpc.NewHandler(logger, segmentService)
	server := grpcServer(logger)
	pb.RegisterSegmentationServiceServer(server, grpcHandler)

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
	<-ch

	logger.Info("Stopping the server")
	server.Stop()
	logger.Info("Closing the listener")
	lis.Close()
	logger.Info("Server Shutdown")
	segmentService.Stop()
}
