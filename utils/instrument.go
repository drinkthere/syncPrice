package utils

import (
	"strings"
	"syncPrice/config"
)

func ConvertToStdInstType(exchange config.Exchange, instType string) config.InstrumentType {
	if exchange == config.OkxExchange {
		switch instType {
		case "SWAP":
			return config.UPerpInstrument
		case "SPOT":
			return config.SpotInstrument
		}
	}
	return config.UnknownInstrument
}

func GetDecimals(numString string) int {
	// 去掉小数点后多余的零
	numString = strings.TrimRight(numString, "0")
	numString = strings.TrimSuffix(numString, ".")

	// 检查是否包含小数点
	if strings.Contains(numString, ".") {
		decimalIndex := strings.Index(numString, ".")
		decimals := len(numString) - decimalIndex - 1
		return decimals
	}

	// 如果没有小数点，返回 0
	return 0
}
