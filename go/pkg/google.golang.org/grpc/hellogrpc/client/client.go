package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/credentials/insecure"
	"grpc/hellogrpc"
)

func init() {
	// resolver.Register(&MyResolveBuilder{})
	balancer.Register(newMyBalanceBuilder())
}
func main() {
	var ctx = context.Background()
	var conn *grpc.ClientConn
	var err error

	conn, err = grpc.DialContext(ctx, "myresolve:///mytarget",
		grpc.WithResolvers(&MyResolveBuilder{}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, lbName)),
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
