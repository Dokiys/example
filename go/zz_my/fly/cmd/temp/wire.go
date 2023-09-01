//go:build wireinject

package temp

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewSubOne,
	NewSubOneCmd,

	NewSubTwo,
	NewSubTwoCmd,

	NewTempCmd,
	wire.Struct(new(Temp), "*"),
)
