package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.NewSource(time.Now().UnixNano())
}

func Abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

// 限制3位小数
func FloatLimit(input float64) float64 {
	if input > 10 { // 10m/s2过大，fix with random
		input = rand.Float64() * 10
	}

	s := fmt.Sprintf("%03f", input)
	q, _ := strconv.ParseFloat(s, 64)
	return q
}
