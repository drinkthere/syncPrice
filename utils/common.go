package utils

import (
	"math"
	"time"
)

func MaxFloat64(list []float64) (max float64) {
	max = list[0]
	for _, v := range list {
		if v > max {
			max = v
		}
	}
	return
}

func MinFloat64(list []float64) (min float64) {
	min = list[0]
	for _, v := range list {
		if v < min {
			min = v
		}
	}
	return
}

func MinDecimal(decimalPlaces int) float64 {
	if decimalPlaces < 0 {
		return 0 // 负数位数返回 0
	}
	return 1 / math.Pow(10, float64(decimalPlaces))
}

func Round(value float64, decimals int) float64 {
	shift := math.Pow(10, float64(decimals))
	rounded := math.Round(value*shift) / shift
	return rounded
}

func InArray(target string, strArray []string) bool {
	for _, element := range strArray {
		if target == element {
			return true
		}
	}
	return false
}

func GetTimestampInMS() int64 {
	return time.Now().UnixNano() / 1e6
}

func IsMsExpired(updateTimeMs int64, milli int64) bool {
	return time.Now().UnixMilli()-updateTimeMs > milli
}
