package t2

import (
	"testing"
	"time"
)

func TestA(t *testing.T) {
	t.Log(time.Now(), "a")
}

func TestZ(t *testing.T) {
	time.Sleep(1 * time.Second)
	t.Log(time.Now(), "z")
}

func TestB(t *testing.T) {
	t.Log(time.Now(), "b")
}
