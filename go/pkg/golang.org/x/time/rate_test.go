package time

import (
	"context"
	"sync"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestRate(t *testing.T) {
	duration, count := 1*time.Millisecond, 5
	limiter := rate.NewLimiter(rate.Every(duration), count)
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	i := 1
	for {
		select {
		case <-ctx.Done():
			goto RETURN
		default:
			t.Log("Start a new return(", i, ")...")
			wg.Add(1)
			for j := 1; j <= i; j++ {
				go func(id int) {
					wg.Wait()
					if limiter.Allow() {
						t.Log(time.Now(), "Goroutine(", id, ") request.")
					} else {
						cancel()
					}
				}(j)
			}

			wg.Done()
			time.Sleep(duration + 1000*time.Millisecond)
			i += 1
		}
	}

RETURN:
	t.Log("finished!")
}
