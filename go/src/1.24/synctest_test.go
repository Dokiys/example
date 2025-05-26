package __24

import (
	"testing"
	"testing/synctest"
	"time"
)

func TestSyncTest(t *testing.T) {
	synctest.Run(func() {
		before := time.Now()
		time.Sleep(time.Second)
		after := time.Now()
		if d := after.Sub(before); d != time.Second {
			t.Fatalf("took %v", d)
		}
	})
}
