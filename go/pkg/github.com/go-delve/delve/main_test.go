package delve

import (
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	println("start", time.Now().Nanosecond())
	m.Run()
	println("end", time.Now().Nanosecond())
}
