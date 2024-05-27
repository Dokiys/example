package main

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/resolver"
)

/*
Server discovery
*/
const scheme = "myresolve"

// MyResolveBuilder Implement ResolveBuilder
type MyResolveBuilder struct{}

func (m *MyResolveBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// 发现服务
	// state = resolver.State{
	//	Addresses: []resolver.Address{{Addr: addrM[target.Endpoint]}},
	// }
	var state resolver.State
	if target.Endpoint() == "mytarget" {
		state = resolver.State{
			// Addresses: []resolver.Address{{Addr: "localhost:50055"}},
			Addresses: []resolver.Address{{Addr: "localhost:50055"}, {Addr: "localhost:50056"}},
		}
	}
	err := cc.UpdateState(state)
	if err != nil {
		cc.ReportError(errors.Wrapf(err, "更新State失败："))
	}
	return &MyResolver{}, nil
}
func (m *MyResolveBuilder) Scheme() string { return scheme }

// MyResolver Implement Resolver
type MyResolver struct{}

func (m *MyResolver) ResolveNow(resolver.ResolveNowOptions) {}
func (m *MyResolver) Close()                                {}
