package client

import "syncPrice/config"

// ExchangeClient 定义交易所客户端接口
type ExchangeClient interface {
	Init(cfg *config.Config) bool
}

// NewExchangeClient 根据交易所名称创建客户端实例
func NewExchangeClient(exchange config.Exchange, cfg *config.Config) (ExchangeClient, bool) {
	switch exchange {
	case config.BinanceExchange:
		client := &BinanceClient{}
		return client, client.Init(cfg)
	case config.OkxExchange:
		client := &OkxClient{}
		return client, client.Init(cfg)
	case config.BybitExchange:
		client := &BybitClient{}
		return client, client.Init(cfg)
	default:
		return nil, false
	}
}
