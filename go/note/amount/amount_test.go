package amount

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

// 浮点数计算会有丢失精度问题，所以涉及到金额的计算一般需要转换成整数
// 计算。丢失精度可能出现在任何加减乘除运算中。
func TestLossOfPrecision(t *testing.T) {
	var num float64
	num = num + 0.1
	num = num + 0.1
	num = num + 0.1
	fmt.Println(num) // 0.30000000000000004

	var num2 float64
	num2 = num2 + 0.1
	fmt.Println(num2 * 3) // 0.30000000000000004

	var num3 float64 = 1.16
	fmt.Println(num3 * 100)  // 115.99999999999999
	fmt.Println(100.0 / 6.0) // 16.666666666666668
	fmt.Println(200.0 / 3.0) // 66.66666666666667

	var num4 float64 = 115.99999999999999
	fmt.Println(num4 / 100) // 1.16
}

// 所以涉及金额的浮点数计算，务必将该金额转换成最小单位的int类型数值（比如元需要分*100）
// 然后使用 math.Round 来获取最终值。
func TestPrecision(t *testing.T) {
	var numRound float64
	numRound = numRound + 0.2
	numRound = numRound + 0.1
	fmt.Println(math.Round(numRound*100) / 100)

	var numInt int = 1
	fmt.Println(float64(numInt) / 10) // 0.1

	var numTwo float64 = 1.2
	fmt.Printf("%.2f\n", numTwo) // 1.20

	var numStr = "1.2"
	f, _ := strconv.ParseFloat(numStr, 64)
	fmt.Printf("%.2f\n", f) // 1.20
}

// subDivided 减除法
func subDivided(amount, count int) []int {
	if count == 1 {
		return []int{amount}
	}

	ans := make([]int, 0, count)
	for count > 0 {
		aPart := amount / count
		amount = amount - aPart
		ans = append(ans, aPart)
		count--
	}

	return ans
}

func TestSubDivided(t *testing.T) {
	t.Log(subDivided(8, 3))  // []int{2, 3, 3}
	t.Log(subDivided(9, 3))  // []int{3, 3, 3}
	t.Log(subDivided(10, 3)) // []int{3, 3, 4}
	t.Log(subDivided(11, 3)) // []int{3, 4, 4}
}
