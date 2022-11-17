package sync

import (
	"sync/atomic"
	"testing"
)

func TestAtomic(t *testing.T) {
	var i int32 = 1
	atomic.CompareAndSwapInt32(&i, 2, i+1)
	t.Logf("i: %d", i)
}
