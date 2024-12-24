package main

import (
	"context"
	"os"
	"runtime/debug"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/mahsandr/arman-challenge/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"sync"
)

var once sync.Once
var serviceConfig *config.ServiceConfig

func main() {
	var err error
	once.Do(func() {
		serviceConfig, err = config.GetConfig()
		if err != nil {
			panic(err)
		}
	})
	switch os.Args[1] {
	case "segmentation":
		RunSegmentationService()
	case "estimation":
		RunEstimationService()
	default:
		panic("invalid " + os.Args[1])
	}
}
func checkErr(logger *zap.Logger, msg string, err error) {
	if err != nil {
		logger.Fatal(msg,
			zap.Error(err))
		return
	}
}
func grpcServer(logger *zap.Logger) *grpc.Server {
	unaryServerOptions := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.UnaryServerInterceptor(logger),
		grpc_zap.PayloadUnaryServerInterceptor(logger, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
			return os.Getenv("Log_REQUEST") == "true"
		}),
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("stack trace from panic " + string(debug.Stack()))
			return status.Errorf(codes.Internal, "%v", p)
		})),
	}
	return grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryServerOptions...)),
	)
}
