package main

import (
	"fmt"
)

var Speed []int
var Efficiency []int
var MaxRes int
var N int

func maxPerformance(n int, speed []int, efficiency []int, k int) int {
	Speed = speed
	Efficiency = efficiency
	N = n
	dfshelper(make([]int, 0), 0, k)
	return MaxRes
}

func dfshelper(tmparr []int, i, k int) {
	if k <= 0 {
		tmpsum := calc(tmparr)
		MaxRes = max(MaxRes, tmpsum)
		return
	}
	if i >= N {
		return
	}
	dfshelper(append(tmparr, i), i+1, k-1)
	dfshelper(tmparr, i+1, k)
}

func calc(tmparr []int) int {
	mine := Efficiency[tmparr[0]]
	sum := Speed[tmparr[0]]
	for i, v := range tmparr {
		if i == 0 {
			continue
		}
		sum += Speed[v]
		mine = min(mine, Efficiency[v])
	}
	fmt.Println(tmparr, mine, sum)

	return mine * sum
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main12() {
	speed := []int{2, 8, 2}
	efficiency := []int{2, 7, 1}
	res := maxPerformance(3, speed, efficiency, 2)
	fmt.Println(res)

}
