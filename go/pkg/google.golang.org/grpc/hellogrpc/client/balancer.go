package main

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

/*
Load Balancer
*/
const lbName = "mybalancer"

func newMyBalanceBuilder() balancer.Builder {
	return base.NewBalancerBuilder(lbName, &MyPickerBuilder{}, base.Config{HealthCheck: true})
}

// MyPickerBuilder Implement PickerBuilder
type MyPickerBuilder struct{}

func (m *MyPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	scs := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for sc := range info.ReadySCs {
		scs = append(scs, sc)
	}
	return &MyPicker{subConns: scs}
}

// MyPicker Implement Picker
type MyPicker struct {
	subConns []balancer.SubConn
	next     int
}

func (self *MyPicker) Pick(_ balancer.PickInfo) (balancer.PickResult, error) {
	sc := self.subConns[self.next]
	self.next = (self.next + 1) % len(self.subConns)
	return balancer.PickResult{SubConn: sc}, nil
}
