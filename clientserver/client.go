package main

import (
	"context"
	"fmt"
	"log"

	hello "github.com/zeimedee/goGrpc/hello"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello client ....")
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	client := hello.NewHelloServiceClient(cc)
	request := &hello.HelloRequest{Name: "Abu"}

	resp, _ := client.Hello(context.Background(), request)
	fmt.Printf("Receive => [%v}", resp.Greeting)
}
