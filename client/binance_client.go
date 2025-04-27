package client

import (
	"github.com/dictxwang/go-binance"
	"github.com/dictxwang/go-binance/futures"
	"syncPrice/config"
)

type BinanceClient struct {
	futuresClient *futures.Client
	spotClient    *binance.Client
}

func (cli *BinanceClient) Init(cfg *config.Config) bool {
	if cfg.BinanceLocalIP == "" {
		cli.futuresClient = futures.NewClient(cfg.BinanceAPIKey, cfg.BinanceSecretKey)
		cli.spotClient = binance.NewClient(cfg.BinanceAPIKey, cfg.BinanceSecretKey)
	} else {
		cli.futuresClient = futures.NewClientWithIP(cfg.BinanceAPIKey, cfg.BinanceSecretKey, cfg.BinanceLocalIP)
		cli.spotClient = binance.NewClientWithIP(cfg.BinanceAPIKey, cfg.BinanceSecretKey, cfg.BinanceLocalIP)
	}
	return true
}
