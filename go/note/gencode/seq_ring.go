package gencode

import "time"

type ringResult struct {
	seq uint32
	now time.Time
}

type Ring struct {
	ch  chan ringResult
	Min uint32
	Max uint32
}

func newRing(min, max uint32) *Ring {
	r := &Ring{
		ch:  make(chan ringResult, 10),
		Min: min,
		Max: max,
	}
	r.init()
	return r
}

func (r *Ring) init() {
	go func() {
		const interval = 1 * time.Second
		var begin = time.Now()
		var ticker = time.NewTicker(interval)
		var i = r.Min
		for {
			select {
			case begin = <-ticker.C:
				i = r.Min // 1000
			case r.ch <- ringResult{i, begin}:
				i++
				if i <= r.Max /*9999*/ {
					continue
				}

				// 取完一轮后等待到下一个单位时间
				begin = <-ticker.C
				i = r.Min
			}
		}
	}()
}

func (r *Ring) Code() (time.Time, uint32) {
	next := <-r.ch
	return next.now, next.seq
}
