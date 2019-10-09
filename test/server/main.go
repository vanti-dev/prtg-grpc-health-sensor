package main

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpc_health "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	failureMode = flag.Bool("fail", false, "Force NOT_SERVING status")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	defer server.GracefulStop()
	health := health.NewServer()
	defer health.Shutdown()

	grpc_health.RegisterHealthServer(server, health)

	if *failureMode {
		health.SetServingStatus("", grpc_health.HealthCheckResponse_NOT_SERVING)
	}

	err = server.Serve(lis)
	log.Fatalf("Error: %v", err)
}
