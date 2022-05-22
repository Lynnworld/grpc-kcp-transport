package main

import (
	"context"
	"fmt"
	kcp "github.com/Lynnworld/grpc-kcp-transport"
	demo "github.com/Lynnworld/grpc-kcp-transport/example/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	address := "localhost:8080"
	cfg := &kcp.Config{}
	cc, err := grpc.Dial(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		kcp.Dialer(cfg))
	if err != nil {
		panic(err)
	}
	client := demo.NewDemoClient(cc)
	resp, _ := client.Hello(context.Background(), &demo.HelloRequest{Name: "gRPC-KCP"})
	fmt.Println(resp.Message)

	// you can also create a generic client
	cc2, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client2 := demo.NewDemoClient(cc2)
	resp2, _ := client2.Hello(context.Background(), &demo.HelloRequest{Name: "gRPC-TCP"})
	fmt.Println(resp2.Message)

}
