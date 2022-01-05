package main

import (
	"context"
	"fmt"
	"log"
	"net"

	hello "github.com/zeimedee/goGrpc/hello"
	"google.golang.org/grpc"
)

type server struct {
	hello.UnimplementedHelloServiceServer
}

func (*server) Hello(ctx context.Context, request *hello.HelloRequest) (*hello.HelloResponse, error) {
	name := request.Name
	response := &hello.HelloResponse{
		Greeting: "Hello " + name,
	}
	return response, nil
}

func main() {
	address := "0.0.0.0:50051"

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Server is listening on: %v...", address)

	s := grpc.NewServer()
	hello.RegisterHelloServiceServer(s, &server{})

	s.Serve(lis)
}
