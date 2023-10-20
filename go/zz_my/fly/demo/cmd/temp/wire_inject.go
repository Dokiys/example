//go:build wireinject

package temp

import (
	"github.com/Dokiys/go_test/go/zz_my/fly/demo/conf"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	conf.NewConfig,
	NewSubOne,
	NewSubTwo,
)
