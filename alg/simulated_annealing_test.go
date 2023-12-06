package alg

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

// 目标函数，可以根据具体问题进行定义
func costFunction(x float64) float64 {
	return math.Pow(x, 2)
}

// 模拟退火算法
func simulatedAnnealing() float64 {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// 初始温度和终止温度
	initialTemp := 100.0
	finalTemp := 0.1

	// 冷却率
	coolingRate := 0.95

	// 当前解和最优解
	currentSolution := rand.Float64() * 100 // 初始解
	bestSolution := currentSolution

	// 当前解的成本和最优解的成本
	currentCost := costFunction(currentSolution)
	bestCost := currentCost

	// 迭代循环，直到温度降到终止温度
	for temp := initialTemp; temp > finalTemp; temp *= coolingRate {
		// 生成新的候选解
		newSolution := currentSolution + (rand.Float64()*2 - 1)

		// 计算新解的成本
		newCost := costFunction(newSolution)

		// 计算成本差
		costDiff := newCost - currentCost

		// 如果新解成本更低，或者按照一定概率接受成本更高的解
		if costDiff < 0 || math.Exp(-costDiff/temp) > rand.Float64() {
			currentSolution = newSolution
			currentCost = newCost
		}

		// 更新最优解
		if currentCost < bestCost {
			bestSolution = currentSolution
			bestCost = currentCost
		}
	}

	return bestSolution
}

func TestSA(t *testing.T) {
	bestSolution := simulatedAnnealing()
	bestCost := costFunction(bestSolution)

	fmt.Printf("Best Solution: %.2f\n", bestSolution)
	fmt.Printf("Best Cost: %.2f\n", bestCost)
}
