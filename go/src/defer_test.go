package src

import (
	"fmt"
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

func TestDeferReturn(t *testing.T) {
	fn := func() (err error) {
		defer func() {
			t.Log(err) // 123
		}()
		if err := fmt.Errorf("123"); err != nil {
			return err
		}
		return nil
	}
	fn()
}
