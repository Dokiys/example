package gencode

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkName(b *testing.B) {
	m := make(map[string]struct{})
	g := NewGenerator("test_hash_salt", "prefix-", &localWorkIdPicker{})
	finish := time.After(2 * time.Second)
	for {
		select {
		case <-finish:
			goto FINISH
		default:
			code := g.Code()
			if _, ok := m[code]; ok {
				panic("conflict")
			} else {
				m[code] = struct{}{}
			}
		}
	}
FINISH:
	assert.Equal(b, 18001, len(m))
}
