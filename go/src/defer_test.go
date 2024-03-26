package src

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
		defer t.Log(err) // nil
		defer func() {
			t.Log(err) // 123
		}()
		if err := fmt.Errorf("123"); err != nil {
			return err
		}
		return nil
	}
	_ = fn()
}

func TestDeferClosure(t *testing.T) {
	assert.Equal(t, closureErr().Error(), "test")
}

func closureErr() (err error) {
	fn := func(err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	defer func() { fn(err) }()
	return fmt.Errorf("test")
}
