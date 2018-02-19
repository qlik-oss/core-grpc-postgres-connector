package main

import (
	"fmt"
	"log"
	"net"

	"github.com/qlik-ea/postgres-grpc-connector/postgres"
	qlik "github.com/qlik-ea/postgres-grpc-connector/qlik"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	srv := &server{make(map[string]*postgres.Reader)}
	qlik.RegisterConnectorServer(s, srv)

	// Register reflection service on gRPC server.
	reflection.Register(s)
	fmt.Println("Server started", port)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
