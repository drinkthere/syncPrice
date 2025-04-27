package message

import (
	"github.com/drinkthere/bybit"
	"math/rand"
	"strconv"
	"syncPrice/config"
	"syncPrice/container"
	"syncPrice/context"
	"syncPrice/utils"
	"syncPrice/utils/logger"
)

func StartGatherBybitUPerpTicker(
	globalConfig *config.Config,
	globalContext *context.GlobalContext,
	tickerChan chan *bybit.V5WebsocketPublicTickerResponse) {

	r := rand.New(rand.NewSource(2))
	go func() {
		defer func() {
			logger.Error("[GatherBybitMarketMsg] Bybit UPerp Tickers Gather Exited.")
		}()

		instIDs := globalContext.InstrumentComposite.GetInstIDs(config.BybitExchange, config.UPerpInstrument)
		wrapper := globalContext.TickerComposite.GetTickerWrapper(config.BybitExchange, config.UPerpInstrument)
		if wrapper == nil {
			logger.Error("[GatherBybitMarketMsg] Bybit UPerp Tickers Wrapper is Nil")
			return
		}

		for t := range tickerChan {
			tickerInfo := t.Data.LinearInverse
			if !utils.InArray(string(tickerInfo.Symbol), instIDs) {
				continue
			}
			tickerMsg := convertBybitUPerpTickerEventToTickerMessage(tickerInfo, t.TimeStamp)
			wrapper.UpdateTicker(tickerMsg)

			if r.Int31n(10000) < 2 {
				logger.Info("[GatherBybitMarketMsg] Receive Bybit UPerp Tickers %+v", *t.Data.LinearInverse)
			}
		}
	}()

	logger.Info("[GatherBybitMarketMsg] Start Gather Bybit UPerp Tickers")
}

func StartGatherBybitSpotTicker(globalConfig *config.Config,
	globalContext *context.GlobalContext,
	tickerChan chan *bybit.V5WebsocketPublicTickerResponse) {

	r := rand.New(rand.NewSource(2))
	go func() {
		defer func() {
			logger.Error("[GatherBybitMarketMsg] Bybit Spot Tickers Gather Exited.")
		}()

		instIDs := globalContext.InstrumentComposite.GetInstIDs(config.BybitExchange, config.SpotInstrument)
		wrapper := globalContext.TickerComposite.GetTickerWrapper(config.BybitExchange, config.SpotInstrument)
		if wrapper == nil {
			logger.Error("[GatherBybitMarketMsg] Bybit Spot Tickers Wrapper is Nil")
			return
		}
		for t := range tickerChan {
			tickerInfo := t.Data.Spot
			if !utils.InArray(string(tickerInfo.Symbol), instIDs) {
				continue
			}
			tickerMsg := convertBybitSpotTickerEventToTickerMessage(globalContext, tickerInfo, t.TimeStamp)
			if tickerMsg != nil {
				wrapper.UpdateTicker(*tickerMsg)

				if r.Int31n(10000) < 2 {
					logger.Info("[GatherBybitMarketMsg] Receive Bybit Spot Tickers %+v", *t.Data.Spot)
				}
			}
		}
	}()

	logger.Info("[GatherBybitMarketMsg] Start Gather Bybit Spot Tickers")
}

func convertBybitSpotTickerEventToTickerMessage(globalContext *context.GlobalContext, ticker *bybit.V5WebsocketPublicTickerSpotResult, updateTs int64) *container.TickerMessage {
	instID := string(ticker.Symbol)
	value, ok := globalContext.InstrumentComposite.BybitSpotInstPrecisionMap.Load(instID)
	if ok {
		precision := value.(container.BybitSpotPrecision) // 转换为 []int

		lastPrice, _ := strconv.ParseFloat(ticker.LastPrice, 64)
		// 数据中没有方向，所以bid和ask同时变化了，否则只需要变动一侧就可以
		bidPx := lastPrice - precision.PriceTick
		askPx := lastPrice + precision.PriceTick

		// 因为bybit的spot ticker数据中，没有包含数量，这里用精度计算一个值
		sz := precision.MinSize

		return &container.TickerMessage{
			Exchange: config.BybitExchange,
			InstType: config.SpotInstrument,
			Ticker: container.Ticker{
				InstID:       instID,
				AskPx:        askPx,
				AskSz:        sz,
				BidPx:        bidPx,
				BidSz:        sz,
				UpdateTimeMs: updateTs,
			},
		}
	}
	return nil
}

func convertBybitUPerpTickerEventToTickerMessage(ticker *bybit.V5WebsocketPublicTickerLinearInverseResult, updateTs int64) container.TickerMessage {
	instID := string(ticker.Symbol)

	bidPx, _ := strconv.ParseFloat(ticker.Bid1Price, 64)
	bidSz, _ := strconv.ParseFloat(ticker.Bid1Size, 64)
	askPx, _ := strconv.ParseFloat(ticker.Ask1Price, 64)
	askSz, _ := strconv.ParseFloat(ticker.Ask1Size, 64)

	return container.TickerMessage{
		Exchange: config.BybitExchange,
		InstType: config.UPerpInstrument,
		Ticker: container.Ticker{
			InstID:       instID,
			AskPx:        askPx,
			AskSz:        askSz,
			BidPx:        bidPx,
			BidSz:        bidSz,
			UpdateTimeMs: updateTs,
		},
	}
}
