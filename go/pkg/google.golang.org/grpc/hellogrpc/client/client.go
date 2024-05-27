package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"grpc/hellogrpc"
)

func init() {
	resolver.Register(&MyResolveBuilder{})
	balancer.Register(newMyBalanceBuilder())
}
func main() {
	var conn *grpc.ClientConn
	var err error

	conn, err = grpc.Dial("myresolve:///mytarget",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, lbName)),
		// grpc.WithBalancerName("mybalancer"), // same as above
	)
	if err != nil {
		log.Fatalf("failed to Dial: %v", err)
	}
	c := hellogrpc.NewGreeterClient(conn)

	for i := 0; i < 5; i++ {
		r, err := c.SayHello(context.Background(), &hellogrpc.HelloRequest{Name: "zhangsan"})
		if err != nil {
			// status.Status
			// s, _ := status.FromError(err)
			log.Fatal(err)
		}
		log.Printf("Greeting: %s", r.GetMessage())
	}
}
