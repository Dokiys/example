//go:build wireinject

package temp

import (
	"github.com/Dokiys/go_test/go/zz_my/fly"
	"github.com/google/wire"
)

func NewSubOneFly(loader fly.ConfigLoader) *SubOne {
	panic(wire.Build(ProviderSet))
}

func NewSubTwoFly(loader fly.ConfigLoader) *SubTwo {
	panic(wire.Build(ProviderSet))
}
