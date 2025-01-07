package chuiniu

import (
	"math/rand"
	"testing"

	"github.com/samber/lo"
)

func TestChuiniu1(t *testing.T) {
	var m = make([]int, 8)
	for i := 0; i < 100000; i++ {
		r := result(t)
		for j := 1; j < 6; j++ {
			m[r[j]] += 1
		}
	}

	var count int
	for _, v := range m {
		count += v
	}

	t.Log("对方可能的点数概率")
	for i := range m {
		var sum int
		for _, v2 := range m[i:] {
			sum += v2
		}
		t.Logf("%d+: \t%.f%%", i, float64(sum)/float64(count)*100)
	}
}

func result(t *testing.T) []int {
	var a = make([]int, 6)
	for i := 0; i < 5; i++ {
		a[rand.Intn(6)] += 1
	}
	if countZero := lo.Count(a, 0); countZero == 1 {
		return []int{0, 0, 0, 0, 0, 0}
	}
	if a[0] == 5 {
		return []int{7, 7, 7, 7, 7, 7}

	}

	var r = make([]int, 6)
	for i := 1; i < 6; i++ {
		if a[i] == 5 {
			r[i] = 7
		}

		r[i] = a[i] + a[0]
		if r[i] == 5 {
			r[i] = 5
		}
	}
	if r[0] == 5 {
		for i := range r {
			r[i] = 7
		}
	}

	return r
}
