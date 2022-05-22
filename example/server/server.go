package main

import (
	"context"
	kcp "github.com/Lynnworld/grpc-kcp-transport"
	demo "github.com/Lynnworld/grpc-kcp-transport/example/api"
	"google.golang.org/grpc"
	"net"
)

func main() {
	address := ":8080"
	config := &kcp.Config{}

	server := grpc.NewServer()
	demo.RegisterDemoServer(server, &serverImpl{})

	listener, err := kcp.Listener(address, config)
	if err != nil {
		panic(err)
	}
	go server.Serve(listener)
	// you can also use server.Serve(l) tcp on same address

	l2, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	server.Serve(l2)

}

type serverImpl struct {
}

func (s *serverImpl) Hello(ctx context.Context, in *demo.HelloRequest) (*demo.HelloReply, error) {
	return &demo.HelloReply{Message: "Hello " + in.Name}, nil
}
