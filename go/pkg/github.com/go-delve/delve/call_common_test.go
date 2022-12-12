package delve

import "testing"

func TestCallCommon(t *testing.T) {
	t.Run("dlv debug", func(t *testing.T) {
		Common("123")
	})
}
