package utils

import (
	"math"
	"math/rand"
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

func FloatEqual(i, j float32) bool {
	return math.Abs(float64(i-j)) < 10e-8
}

//RandomFloat64 随机生成一个min到max的float64的数
func RandomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func RandomInt32(min, max int32) int32 {
	return rand.Int31n(max-min) + min
}

//TransNaN 有时运算结果太小产生的NaN变成0
func TransNaN(f *float64) {
	if math.IsNaN(*f) {
		*f = 0
	}
}
