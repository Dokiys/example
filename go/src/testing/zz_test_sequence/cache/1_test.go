package t1

import (
	"testing"
	"time"
)

func Test1_1(t *testing.T) {
	t.Log("1_1")
}

func Test1_2(t *testing.T) {
	t.Log("1_2")
}

func TestTimeCache(t *testing.T) {
	nano := time.Now().Unix()
	if nano <= 1668152900 {
		t.Fatal("err")
	}

	// rand.Seed(int64(time.Now().Nanosecond()))
	// if rand.Intn(10) > 5 {
	//	t.Fatal("err")
	// }
}
