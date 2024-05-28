package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

/*
Server discovery
*/
const scheme = "myresolve"

// MyResolveBuilder Implement ResolveBuilder
type MyResolveBuilder struct{}

func (m *MyResolveBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	if target.Endpoint() != "mytarget" {
		return nil, errors.New("only mytarget is supported")
	}
	r := &MyResolver{ctx: context.Background(), cc: cc}
	go r.update()
	return r, nil
}
func (m *MyResolveBuilder) Scheme() string { return scheme }

// MyResolver Implement Resolver
type MyResolver struct {
	ctx context.Context
	cc  resolver.ClientConn
}

func (r *MyResolver) ResolveNow(resolver.ResolveNowOptions) {}
func (r *MyResolver) Close()                                {}
func (r *MyResolver) update() {
	// 发现服务
	var state = resolver.State{
		Addresses: []resolver.Address{
			{Addr: "localhost:50055"},
			{Addr: "localhost:50056", BalancerAttributes: attributes.New("weight", 0)},
		},
	}
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
		}
		var err = r.cc.UpdateState(state)
		if err != nil {
			r.cc.ReportError(errors.Wrapf(err, "更新State失败："))
		}
		time.Sleep(time.Second)
		fmt.Printf("connected: %s\n", state.Addresses)
	}
}
