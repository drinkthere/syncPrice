package message

import (
	"github.com/drinkthere/okx/events/public"
	"github.com/drinkthere/okx/models/market"
	"math/rand"
	"syncPrice/config"
	"syncPrice/container"
	"syncPrice/context"
	"syncPrice/utils"
	"syncPrice/utils/logger"
	"time"
)

func convertToOkxTickerMessage(ticker *market.Ticker) container.TickerMessage {
	return container.TickerMessage{
		Exchange: config.OkxExchange,
		InstType: utils.ConvertToStdInstType(config.OkxExchange, string(ticker.InstType)),
		Ticker: container.Ticker{
			InstID:       ticker.InstID,
			AskPx:        float64(ticker.AskPx),
			AskSz:        float64(ticker.AskSz),
			BidPx:        float64(ticker.BidPx),
			BidSz:        float64(ticker.BidSz),
			UpdateTimeMs: time.Time(ticker.TS).UnixMilli(),
			UpdateID:     0,
		},
	}
}

func StartGatherOkxUPerpTicker(
	globalConfig *config.Config,
	globalContext *context.GlobalContext,
	tickChan chan *public.Tickers) {

	r := rand.New(rand.NewSource(2))
	go func() {
		defer func() {
			if rc := recover(); rc != nil {
				logger.Error("[GatherOkxTicker] Recovered Okx %s ticker from panic: %v", config.UPerpInstrument, rc)
			}

			logger.Warn("[GatherOkxTicker] Okx %s Ticker Gather Exited.", config.UPerpInstrument)
		}()

		instIDs := globalContext.InstrumentComposite.GetInstIDs(config.OkxExchange, config.UPerpInstrument)
		wrapper := globalContext.TickerComposite.GetTickerWrapper(config.OkxExchange, config.UPerpInstrument)
		if wrapper == nil {
			logger.Error("[GatherOkxTicker] Okx %s Tickers Wrapper is Nil", config.UPerpInstrument)
			return
		}

		for {
			s := <-tickChan
			for _, t := range s.Tickers {
				if !utils.InArray(t.InstID, instIDs) {
					continue
				}
				tickerMsg := convertToOkxTickerMessage(t)
				wrapper.UpdateTicker(tickerMsg)
			}
			if r.Int31n(10000) < 5 && len(s.Tickers) > 0 {
				logger.Info("[GatherOkxTicker] Receive Okx %s Ticker %+v", s.Tickers[0], config.UPerpInstrument)
			}
		}
	}()

	logger.Info("[GatherOkxTicker] Start Gather Okx %s Ticker", config.UPerpInstrument)
}

func StartGatherOkxSpotTicker(
	globalConfig *config.Config,
	globalContext *context.GlobalContext,
	tickChan chan *public.Tickers) {

	r := rand.New(rand.NewSource(2))
	go func() {
		defer func() {
			if rc := recover(); rc != nil {
				logger.Error("[GatherOkxTicker] Recovered Okx %s ticker from panic: %v", config.SpotInstrument, rc)
			}

			logger.Warn("[GatherOkxTicker] Okx %s Ticker Gather Exited.", config.SpotInstrument)
		}()

		instIDs := globalContext.InstrumentComposite.GetInstIDs(config.OkxExchange, config.SpotInstrument)
		wrapper := globalContext.TickerComposite.GetTickerWrapper(config.OkxExchange, config.SpotInstrument)
		if wrapper == nil {
			logger.Error("[GatherOkxTicker] Okx %s Tickers Wrapper is Nil", config.SpotInstrument)
			return
		}
		for {
			s := <-tickChan
			for _, t := range s.Tickers {
				if !utils.InArray(t.InstID, instIDs) {
					continue
				}

				tickerMsg := convertToOkxTickerMessage(t)
				wrapper.UpdateTicker(tickerMsg)
			}
			if r.Int31n(10000) < 5 && len(s.Tickers) > 0 {
				logger.Info("[GatherOkxTicker] Receive Okx %s Ticker %+v", s.Tickers[0], config.SpotInstrument)
			}
		}
	}()

	logger.Info("[GatherOkxTicker] Start Gather Okx %s Ticker", config.SpotInstrument)
}
