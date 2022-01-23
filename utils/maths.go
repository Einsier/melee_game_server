package utils

import (
	"math"
)

/**
*@Author Sly
*@Date 2022/1/15
*@Version 1.0
*@Description:
 */

func MinInt32(i, j int32) int32 {
	if i < j {
		return i
	}
	return j
}

func FloatEqual(i, j float64) bool {
	return math.Abs(i-j) < 10e-8
}
