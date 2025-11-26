package localcache

import "testing"

func TestName(t *testing.T) {
	a := map[string]any{}
	a["1"] = 1
	a["2"] = "2"
	a["3"] = true
	t.Log(a)
}
