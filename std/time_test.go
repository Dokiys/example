package std

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestTimeCheckNil 校验时间初始值
func TestTimeCheckNil(t *testing.T) {
	date, _ := time.ParseInLocation("2006-01-02", "", time.Local)
	t.Log(date)

	if date.IsZero() {
		t.Log(date.Format("2006-01-02 15:04:05"))
	}
}

// TestTimeAdd 添加时间
func TestTimeAdd(t *testing.T) {
	t.Log(time.Now().Format("2006-01-02"))
	t.Log(time.Now().AddDate(1, 0, 0).Format("2006-01-02"))
}

// TestTimeNano 输出当前时间
func TestTimeNano(t *testing.T) {
	t.Log(time.Now().UnixNano())
}

// TestTimeStrParse 根据字符串和格式转换成时间
func TestTimeStrParse(t *testing.T) {
	//从字符串转为时间戳，第一个参数是格式，第二个是要转换的时间字符串
	tm, err := time.Parse("01/02/2006", "02/08/2015")
	assert.NoError(t, err)
	tm2, err := time.Parse("2006-01-02 15:04:05", "0000-01-35 26:99:99")
	assert.NoError(t, err)
	tm3, err := time.Parse("2006-01-02 15:04:05", "")
	assert.NoError(t, err)

	fmt.Println(tm.Unix())
	fmt.Println(tm2)
	fmt.Println(tm3)
}

// TestTimeRound 四舍五入获取时间
func TestTimeRound(t *testing.T) {
	// Defining duration of Round method
	//d, _ := time.ParseDuration("3m73.371s")
	d, _ := time.Parse("2006-01-02T15:04:05.99Z", "2017-03-25T10:01:02.5234567Z")

	// Array of m
	R := []time.Duration{
		time.Microsecond,
		time.Second,
		3 * time.Second,
		9 * time.Minute,
	}

	// Using for loop and range to
	// iterate over an array
	for _, m := range R {

		// Prints rounded d of all
		// the items in an array
		fmt.Printf("Rounded(%s) is:%s\n",
			m, d.Round(m))
	}
}
