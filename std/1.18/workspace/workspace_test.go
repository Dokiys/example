package workspace

import (
	"fmt"
	"testing"

	"golang.org/x/example/stringutil"
)

func TestWorkspace(t *testing.T) {
	fmt.Println(stringutil.Reverse("Hello"))
	fmt.Println(stringutil.ToUpper("Hello"))
}
