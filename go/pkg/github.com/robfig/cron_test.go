package robfig

import (
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

func TestCron(t *testing.T) {
	c := cron.New()
	// Accepts this spec: https://en.wikipedia.org/wiki/Cron
	_, err := c.AddFunc("* * * * *", func() {
		t.Log("cron do")
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCronParse(t *testing.T) {
	parser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse("2 * * * * 0")
	if err != nil {
		t.Fatal(err)
	}
	next := schedule.Next(time.Now())
	t.Log(next)
}

func TestCronParseNextAll(t *testing.T) {
	var start, end = time.Now(), time.Now().Add(time.Hour)
	var specs = []string{"2 */2 * * * *", "*/3 * * * *"}
	var internal = time.Minute

	var results = [][]string{}
	var parser = cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	var now = start
	for {
		if now.After(end) {
			break
		}

		result := []string{}
		for _, spec := range specs {
			schedule, err := parser.Parse(spec)
			if err != nil {
				t.Fatal(err)
			}
			next := schedule.Next(now)
			if next.After(end) {
				continue
			}
			result = append(result, next.Format("2006-01-02 15:04:05"))
		}
		if len(result) < 0 {
			continue
		}
		results = append(results, result)
		now = now.Add(internal)
	}

	for _, result := range results {
		t.Log(result)
	}
}
