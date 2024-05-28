package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"grpc/hellogrpc"
)

var addr = flag.String("addr", "localhost:50055", "http service address")

//go:generate go build -o ../bin/server
func main() {
	flag.Parse()
	var s *grpc.Server
	// Create the TLS credentials
	// {
	//	creds, err := credentials.NewServerTLSFromFile("../zchd.crt", "../ca.key")
	//	if err != nil {
	//		log.Fatalf("failed to new tls creds: %v", err)
	//	}
	//	s = grpc.NewServer(grpc.Creds(creds))
	// }
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 前置处理
		fmt.Printf("Before RPC handling. Info: %+v\n", info)
		// 调用方法
		resp, err := handler(ctx, req)
		fmt.Printf("After RPC handling. resp: %+v\n", resp)
		// 后置处理
		return resp, err
	}

	// Create the insecure server
	{
		s = grpc.NewServer(grpc.UnaryInterceptor(interceptor))
		hellogrpc.RegisterGreeterServer(s, &hellogrpc.Server{Addr: *addr})
	}

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
