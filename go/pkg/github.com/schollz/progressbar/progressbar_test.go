package progressbar

import (
	"testing"
	"time"

	"github.com/schollz/progressbar/v3"
)

func TestProgressbar(t *testing.T) {
	bar := progressbar.Default(100)
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(40 * time.Millisecond)
	}
}
