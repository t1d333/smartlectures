package main

import (
	"net"

	storage "github.com/t1d333/smartlectures/internal/storage"
	delivery "github.com/t1d333/smartlectures/internal/storage/delivery/grpc"
	"github.com/t1d333/smartlectures/internal/storage/repository/elasticsearch"
	"github.com/t1d333/smartlectures/internal/storage/service"
	"github.com/t1d333/smartlectures/pkg/logger"
	"google.golang.org/grpc"
)

func main() {
	logger := logger.NewLogger()
	_, err := elasticsearch.NewRepository(logger)
	if err != nil {
		logger.Fatalf("failed to create elasticsearch repository: %s", err)
	}

	logger.Info("creating repository...")
	rep, err := elasticsearch.NewRepository(logger)
	if err != nil {
		logger.Fatalf("failed to create repository %v", err)
	}

	logger.Info("repository created succsessfully")

	service := service.NewService(logger, rep)

	del := delivery.NewDelivery(logger, service)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	logger.Info("starting grpc server on port 50051...")
	grpcServer := grpc.NewServer(opts...)
	storage.RegisterStorageServer(grpcServer, del)
	grpcServer.Serve(lis)
}
