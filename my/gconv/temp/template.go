package temp

import (
	"fmt"

	"go_test/my/genconv/data"
)

func ConvA(a *data.A) (a2 *data.A2) {
	atob := func() {
		a2.Name = fmt.Sprintf("gc_%s", a.Name)
	}
	panic(atob)
}
