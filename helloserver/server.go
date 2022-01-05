package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	hello "github.com/zeimedee/goGrpc/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	address := ":8080"

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Server is listening on: %v...", address)
	//create new grpc server
	s := grpc.NewServer()
	//attach service to server
	hello.RegisterHelloServiceServer(s, &server{})

	//server grpc server, use goroutine since .Serve is a blocking process to allow other processes to run
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	//create a client connection to grpc server
	//this is where grpc-gateway proxies the requests

	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// register greeter
	err = hello.RegisterHelloServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register service", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
