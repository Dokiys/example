package t1

import (
	"testing"
	"time"
)

func Test1_1(t *testing.T) {
	t.Log(time.Now(), "1_1")
}

func Test1_2(t *testing.T) {
	time.Sleep(1 * time.Second)
	t.Log(time.Now(), "1_2")
}
