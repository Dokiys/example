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
	fmt.Println(num3 * 100) // 115.99999999999999
	fmt.Println(1.0 / 6.0)  // 0.16666666666666666
	fmt.Println(2.0 / 3.0)  // 0.6666666666666666

	var num4 float64 = 115.99999999999999
	fmt.Println(num4 / 100) // 1.16
}

// 所以最终计算结果为int类型的时候，务必需要使用 math.Round 来获取最终值。
func TestPrecision(t *testing.T) {
	var numRound float64 = 1.16
	fmt.Println(math.Round(numRound * 100)) // 116

	var numInt int = 1
	fmt.Println(float64(numInt) / 10) // 0.1

	var numTwo float64 = 1.2
	fmt.Printf("%.2f\n", numTwo) // 1.20

	var numStr = "1.2"
	f, _ := strconv.ParseFloat(numStr, 64)
	fmt.Printf("%.2f\n", f) // 1.20
}

// SubDivided 减除法
func SubDivided(amount, count int) []int {
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
	t.Log(SubDivided(8, 3))  // []int{2, 3, 3}
	t.Log(SubDivided(9, 3))  // []int{3, 3, 3}
	t.Log(SubDivided(10, 3)) // []int{3, 3, 4}
	t.Log(SubDivided(11, 3)) // []int{3, 4, 4}
}
