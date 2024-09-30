package main

import (
	grpcServerAdapter "aurma_product/internal/adapters/grpc/server"
	"aurma_product/internal/di"
	"aurma_product/internal/scheduler"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	container, err := di.NewContainer()
	if err != nil {
		log.Printf("Error: Failed to initialize container: %v", err)
		os.Exit(1)
	}
	defer container.Close()
	// Run scheduler
	go scheduler.Run(container)

	// Start gRPC server
	server := grpc.NewServer()
	grpcServerAdapter.Register(server, container)

	lis, err := net.Listen("tcp", ":"+container.Config.GrpcServerPort)
	if err != nil {
		log.Printf("Error: Failed to listen: %v", err)
		os.Exit(1)
	}

	go func() {
		log.Printf("Info: Starting gRPC server on port %s", container.Config.GrpcServerPort)
		if err := server.Serve(lis); err != nil {
			log.Printf("Error: Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Print("Info: Shutting down server and scheduler...")
	server.GracefulStop()

}
