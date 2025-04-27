package client

import (
	"github.com/drinkthere/bybit"
	"net"
	"net/http"
	"syncPrice/config"
	"time"
)

type BybitClient struct {
	Client   *bybit.Client
	WsClient *bybit.WebSocketClient
}

func (cli *BybitClient) Init(cfg *config.Config) bool {
	// bybit的ticker没有内网，所以这里restBaseUrl都设置位空
	restBaseUrl := ""
	if cfg.BybitLocalIP != "" {
		// 创建一个 net.Dialer
		dialer := &net.Dialer{
			Timeout:   5 * time.Second,
			LocalAddr: &net.TCPAddr{IP: net.ParseIP(cfg.BybitLocalIP)},
		}

		// 创建一个 http.Transport，设置 DialContext
		transport := &http.Transport{
			DialContext: dialer.DialContext,
		}

		// 创建 http.Client，使用自定义的 Transport
		httpClient := &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second, // 设置请求超时时间
		}

		cli.Client = bybit.NewClient(restBaseUrl).WithHTTPClient(httpClient).WithAuth(cfg.BybitAPIKey, cfg.BybitSecretKey)
	} else {
		cli.Client = bybit.NewClient(restBaseUrl).WithAuth(cfg.BybitAPIKey, cfg.BybitSecretKey)
	}
	return true
}
