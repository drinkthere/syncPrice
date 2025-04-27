package config

import (
	"encoding/json"
	"go.uber.org/zap/zapcore"
	"os"
)

type Config struct {
	// 日志配置
	CPUNum *int

	// 日志配置
	LogLevel zapcore.Level
	LogPath  string

	// 币安配置
	BinanceAPIKey    string
	BinanceSecretKey string
	BinanceLocalIP   string
	BinanceColo      bool

	// Okx配置
	OkxAPIKey    string
	OkxSecretKey string
	OkxPassword  string
	OkxLocalIP   string
	OkxColo      string // zoneB / zoneD

	// bybit配置
	BybitAPIKey    string
	BybitSecretKey string
	BybitLocalIP   string
	// 公网设置为空即可
	BybitWsBaseUrl string
	BybitColo      bool

	BinanceUPerpInstIDs []string
	BinanceSpotInstIDs  []string
	OkxUPerpInstIDs     []string
	OkxSpotInstIDs      []string
	BybitUPerpInstIDs   []string
	BybitSpotInstIDs    []string

	// 电报配置
	TgBotToken string
	TgChatID   int64

	// 推送中位数价格的频率，单位ms
	SendIntervalMs int64
	MedianPriceZMQ string

	MinAccuracy float64 // 价格最小精度

	// watch dog
	MaxErrorsPerMinute int64
}

func LoadConfig(filename string) *Config {
	config := new(Config)
	reader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	// 加载配置
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}
