package message

import (
	binanceSpot "github.com/dictxwang/go-binance"
	binanceFutures "github.com/dictxwang/go-binance/futures"
	"math/rand"
	"strconv"
	"syncPrice/config"
	"syncPrice/container"
	"syncPrice/context"
	"syncPrice/utils"
	"syncPrice/utils/logger"
	"time"
)

func StartGatherBinanceUPerpTicker(
	globalConfig *config.Config,
	globalContext *context.GlobalContext,
	tickerChan chan *binanceFutures.WsBookTickerEvent) {

	r := rand.New(rand.NewSource(2))
	go func() {
		defer func() {
			logger.Warn("[GatherBinanceMarketMsg] Binance UPerp Tickers Gather Exited.")
		}()
		instIDs := globalContext.InstrumentComposite.GetInstIDs(config.BinanceExchange, config.UPerpInstrument)
		wrapper := globalContext.TickerComposite.GetTickerWrapper(config.BinanceExchange, config.UPerpInstrument)
		if wrapper == nil {
			logger.Error("[GatherBinanceMarketMsg] Binance UPerp Tickers Wrapper is Nil")
			return
		}
		for t := range tickerChan {
			if !utils.InArray(t.Symbol, instIDs) {
				continue
			}
			tickerMsg := convertBinanceUPerpTickerEventToTickerMessage(t)
			wrapper.UpdateTicker(tickerMsg)

			if r.Int31n(10000) < 2 {
				logger.Info("[GatherBinanceMarketMsg] Receive Binance Futures Tickers %+v", t)
			}
		}
	}()

	logger.Info("[GatherBinanceMarketMsg] Start Gather Binance Futures Tickers")
}

func StartGatherBinanceSpotTicker(
	globalConfig *config.Config,
	globalContext *context.GlobalContext,
	tickerChan chan *binanceSpot.WsBookTickerEvent) {

	r := rand.New(rand.NewSource(2))
	go func() {
		defer func() {
			logger.Warn("[GatherBinanceMarketMsg] Binance Spot Tickers Gather Exited.")
		}()

		instIDs := globalContext.InstrumentComposite.GetInstIDs(config.BinanceExchange, config.SpotInstrument)
		wrapper := globalContext.TickerComposite.GetTickerWrapper(config.BinanceExchange, config.SpotInstrument)
		if wrapper == nil {
			logger.Error("[GatherBinanceMarketMsg] Binance Spot Tickers Wrapper is Nil")
			return
		}
		for t := range tickerChan {
			if !utils.InArray(t.Symbol, instIDs) {
				continue
			}
			tickerMsg := convertBinanceSpotTickerEventToTickerMessage(t)
			wrapper.UpdateTicker(tickerMsg)

			if r.Int31n(10000) < 2 {
				logger.Info("[GatherBinanceMarketMsg] Receive Binance Spot Tickers %+v", t)
			}
		}
	}()

	logger.Info("[GatherBinanceMarketMsg] Start Gather Binance Spot Tickers")
}

func convertBinanceSpotTickerEventToTickerMessage(ticker *binanceSpot.WsBookTickerEvent) container.TickerMessage {
	bestAskPrice, _ := strconv.ParseFloat(ticker.BestAskPrice, 64)
	bestAskQty, _ := strconv.ParseFloat(ticker.BestAskQty, 64)
	bestBidPrice, _ := strconv.ParseFloat(ticker.BestBidPrice, 64)
	bestBidQty, _ := strconv.ParseFloat(ticker.BestBidQty, 64)
	return container.TickerMessage{
		Exchange: config.BinanceExchange,
		InstType: config.SpotInstrument,
		Ticker: container.Ticker{
			InstID:       ticker.Symbol,
			AskPx:        bestAskPrice,
			AskSz:        bestAskQty,
			BidPx:        bestBidPrice,
			BidSz:        bestBidQty,
			UpdateID:     ticker.UpdateID,
			UpdateTimeMs: time.Now().UnixMilli(),
		},
	}
}

func convertBinanceUPerpTickerEventToTickerMessage(ticker *binanceFutures.WsBookTickerEvent) container.TickerMessage {
	bestAskPrice, _ := strconv.ParseFloat(ticker.BestAskPrice, 64)
	bestAskQty, _ := strconv.ParseFloat(ticker.BestAskQty, 64)
	bestBidPrice, _ := strconv.ParseFloat(ticker.BestBidPrice, 64)
	bestBidQty, _ := strconv.ParseFloat(ticker.BestBidQty, 64)
	return container.TickerMessage{
		Exchange: config.BinanceExchange,
		InstType: config.UPerpInstrument,
		Ticker: container.Ticker{
			InstID:       ticker.Symbol,
			AskPx:        bestAskPrice,
			AskSz:        bestAskQty,
			BidPx:        bestBidPrice,
			BidSz:        bestBidQty,
			UpdateID:     ticker.UpdateID,
			UpdateTimeMs: ticker.Time,
		},
	}
}
