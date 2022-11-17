package t1

import (
	"testing"
	"time"
)

func Test2_1(t *testing.T) {
	t.Log(time.Now(), "2_1")
}

func Test2_2(t *testing.T) {
	time.Sleep(1 * time.Second)
	t.Log(time.Now(), "2_2")
}
