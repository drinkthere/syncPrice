package main

import (
	binanceSpot "github.com/dictxwang/go-binance"
	binanceFutures "github.com/dictxwang/go-binance/futures"
	"github.com/drinkthere/bybit"
	"github.com/drinkthere/okx/events/public"
	"syncPrice/config"
	"syncPrice/message"
)

func startTickerMessage() {
	startOkxUPerpTicker()
	startOkxSpotTicker()
	startBinanceUPerpTicker()
	startBinanceSpotTicker()
	startBybitUPerpTicker()
	startBybitSpotTicker()
}

func startOkxUPerpTicker() {
	tickerChan := make(chan *public.Tickers)
	message.StartOkxMarketWs(&globalConfig, &globalContext, config.UPerpInstrument, tickerChan)
	message.StartGatherOkxUPerpTicker(&globalConfig, &globalContext, tickerChan)
}

func startOkxSpotTicker() {
	tickerChan := make(chan *public.Tickers)
	message.StartOkxMarketWs(&globalConfig, &globalContext, config.SpotInstrument, tickerChan)
	message.StartGatherOkxSpotTicker(&globalConfig, &globalContext, tickerChan)
}

func startBinanceUPerpTicker() {
	tickerChan := make(chan *binanceFutures.WsBookTickerEvent)
	message.StartBinanceUPerpMarketWs(&globalConfig, &globalContext, tickerChan)
	message.StartGatherBinanceUPerpTicker(&globalConfig, &globalContext, tickerChan)
}

func startBinanceSpotTicker() {
	tickerChan := make(chan *binanceSpot.WsBookTickerEvent)
	message.StartBinanceSpotMarketWs(&globalConfig, &globalContext, tickerChan)
	message.StartGatherBinanceSpotTicker(&globalConfig, &globalContext, tickerChan)
}

func startBybitUPerpTicker() {
	tickerChan := make(chan *bybit.V5WebsocketPublicTickerResponse)
	message.StartBybitUPerpMarketWs(&globalConfig, &globalContext, tickerChan)
	message.StartGatherBybitUPerpTicker(&globalConfig, &globalContext, tickerChan)
}

func startBybitSpotTicker() {
	tickerChan := make(chan *bybit.V5WebsocketPublicTickerResponse)
	message.StartBybitSpotMarketWs(&globalConfig, &globalContext, tickerChan)
	message.StartGatherBybitSpotTicker(&globalConfig, &globalContext, tickerChan)
}
