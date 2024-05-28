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
	return base.NewBalancerBuilder(lbName, &MyPickerBuilder{}, base.Config{HealthCheck: false})
}

// MyPickerBuilder Implement PickerBuilder
type MyPickerBuilder struct{}

func (m *MyPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	scs := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for conn, sc := range info.ReadySCs {
		weight := 1
		if sc.Address.BalancerAttributes != nil {
			val := sc.Address.BalancerAttributes.Value("weight")
			if val != nil {
				weight = val.(int)
			}
		}
		for i := 0; i < weight; i++ {
			scs = append(scs, conn)
		}
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
