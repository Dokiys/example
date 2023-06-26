package src

import (
	"testing"
)

func TestDefer(t *testing.T) {
	fn := func() (str string) {
		defer func() {
			t.Log(str)
		}()
		return "将会输出这里的字符串"
	}
	fn()
}
