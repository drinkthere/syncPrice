package message

import (
	binanceSpot "github.com/dictxwang/go-binance"
	binanceFutures "github.com/dictxwang/go-binance/futures"
	"syncPrice/config"
	"syncPrice/context"
	"syncPrice/utils/logger"
	"time"
)

func StartBinanceUPerpMarketWs(
	globalConfig *config.Config,
	globalContext *context.GlobalContext,
	tickerChan chan *binanceFutures.WsBookTickerEvent) {

	binanceMarketWs := newBinanceUPerpMarketWebSocket(tickerChan)
	binanceMarketWs.subscribeBookTickers(globalConfig, globalContext)
	logger.Info("[BinanceMarketWs] Start Listen Binance UPerp Tickers")
}

type binanceUPerpMarketWebSocket struct {
	tickerChan chan *binanceFutures.WsBookTickerEvent
	isStopped  bool
	stopChan   chan struct{}
}

func newBinanceUPerpMarketWebSocket(tickerChan chan *binanceFutures.WsBookTickerEvent) *binanceUPerpMarketWebSocket {
	return &binanceUPerpMarketWebSocket{
		tickerChan: tickerChan,
		isStopped:  true,
		stopChan:   make(chan struct{}),
	}
}

func (ws *binanceUPerpMarketWebSocket) handleTickerEvent(event *binanceFutures.WsBookTickerEvent) {
	ws.tickerChan <- event
}

func (ws *binanceUPerpMarketWebSocket) handleError(err error) {
	// 出错断开连接，再重连
	logger.Error("[BinanceMarketWs] Binance UPerp Tickers Occur Error And Reconnect Ws: %s", err.Error())
	ws.stopChan <- struct{}{}
	ws.isStopped = true
}

func (ws *binanceUPerpMarketWebSocket) subscribeBookTickers(globalConfig *config.Config, globalContext *context.GlobalContext) {

	go func() {
		defer func() {
			logger.Warn("[BinanceMarketWs] Binance UPerp Tickers Listening Exited.")
		}()

		for {
			if !ws.isStopped {
				time.Sleep(time.Second * 1)
				continue
			}
			ip := globalConfig.BinanceLocalIP
			instIDs := globalContext.InstrumentComposite.GetInstIDs(config.BinanceExchange, config.UPerpInstrument)
			var stopChan chan struct{}
			var err error
			if ip == "" {
				_, stopChan, err = binanceFutures.WsCombinedBookTickerServe(instIDs, ws.handleTickerEvent, ws.handleError)
			} else {
				_, stopChan, err = binanceFutures.WsCombinedBookTickerServeWithIP(ip, instIDs, ws.handleTickerEvent, ws.handleError)
			}

			if err != nil {
				logger.Error("[BinanceMarketWs] Subscribe Binance UPerp Tickers Error: %s", err.Error())
				time.Sleep(time.Second * 1)
				continue
			}
			logger.Info("[BinanceMarketWs] Subscribe Binance UPerp Tickers: %d", len(instIDs))
			// 重置channel和时间
			ws.stopChan = stopChan
			ws.isStopped = false
		}
	}()
}

func StartBinanceSpotMarketWs(
	globalConfig *config.Config,
	globalContext *context.GlobalContext,
	tickerChan chan *binanceSpot.WsBookTickerEvent) {

	binanceMarketWs := newBinanceSpotMarketWebSocket(tickerChan)
	binanceMarketWs.subscribeBookTickers(globalConfig, globalContext)
	logger.Info("[BinanceMarketWs] Start Listen Binance Spot Tickers")
}

type binanceSpotMarketWebSocket struct {
	tickerChan chan *binanceSpot.WsBookTickerEvent
	isStopped  bool
	stopChan   chan struct{}
}

func newBinanceSpotMarketWebSocket(tickerChan chan *binanceSpot.WsBookTickerEvent) *binanceSpotMarketWebSocket {
	return &binanceSpotMarketWebSocket{
		tickerChan: tickerChan,
		isStopped:  true,
		stopChan:   make(chan struct{}),
	}
}

func (ws *binanceSpotMarketWebSocket) handleTickerEvent(event *binanceSpot.WsBookTickerEvent) {
	ws.tickerChan <- event
}

func (ws *binanceSpotMarketWebSocket) handleError(err error) {
	// 出错断开连接，再重连
	logger.Error("[BinanceMarketWs] Binance Spot Tickers Occur Error And Reconnect Ws: %s", err.Error())
	ws.stopChan <- struct{}{}
	ws.isStopped = true
}

func (ws *binanceSpotMarketWebSocket) subscribeBookTickers(globalConfig *config.Config, globalContext *context.GlobalContext) {

	go func() {
		defer func() {
			logger.Warn("[BinanceMarketWs] Binance Spot Tickers Listening Exited.")
		}()
		for {
			if !ws.isStopped {
				time.Sleep(time.Second * 1)
				continue
			}
			ip := globalConfig.BinanceLocalIP
			instIDs := globalContext.InstrumentComposite.GetInstIDs(config.BinanceExchange, config.SpotInstrument)
			var stopChan chan struct{}
			var err error
			if ip == "" {
				_, stopChan, err = binanceSpot.WsCombinedBookTickerServe(instIDs, ws.handleTickerEvent, ws.handleError)
			} else {
				_, stopChan, err = binanceSpot.WsCombinedBookTickerServeWithIP(ip, instIDs, ws.handleTickerEvent, ws.handleError)
			}
			if err != nil {
				logger.Error("[BinanceMarketWs] Subscribe Binance Spot Tickers Error: %s", err.Error())
				time.Sleep(time.Second * 1)
				continue
			}
			logger.Info("[BinanceMarketWs] Subscribe Binance Spot Tickers Successfully: %d", len(instIDs))
			// 重置channel和时间
			ws.stopChan = stopChan
			ws.isStopped = false
		}
	}()
}
