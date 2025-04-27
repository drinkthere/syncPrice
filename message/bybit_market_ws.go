package message

import (
	"context"
	"fmt"
	"github.com/drinkthere/bybit"
	"syncPrice/config"
	mmContext "syncPrice/context"
	"syncPrice/utils/logger"
	"time"
)

func StartBybitUPerpMarketWs(
	globalConfig *config.Config,
	globalContext *mmContext.GlobalContext,
	tickerChan chan *bybit.V5WebsocketPublicTickerResponse) {

	bybitMarketWs := newBybitMarketWebsocket(tickerChan)
	bybitMarketWs.subscribeUPerpTicker(globalConfig, globalContext)
	logger.Info("[BybitMarketWs] Start Listen Bybit UPerp Tickers")
}

func StartBybitSpotMarketWs(
	globalConfig *config.Config,
	globalContext *mmContext.GlobalContext,
	tickerChan chan *bybit.V5WebsocketPublicTickerResponse) {

	bybitMarketWs := newBybitMarketWebsocket(tickerChan)
	bybitMarketWs.subscribeSpotTicker(globalConfig, globalContext)
	logger.Info("[BybitMarketWs] Start Listen Bybit Spot Tickers")
}

type bybitMarketWebSocket struct {
	tickerChan chan *bybit.V5WebsocketPublicTickerResponse
	isStopped  bool
}

func newBybitMarketWebsocket(tickerChan chan *bybit.V5WebsocketPublicTickerResponse) *bybitMarketWebSocket {

	return &bybitMarketWebSocket{
		tickerChan: tickerChan,
		isStopped:  true,
	}
}

func (ws *bybitMarketWebSocket) handleUPerpTickerEvent(event bybit.V5WebsocketPublicTickerResponse) error {
	ws.tickerChan <- &event
	return nil
}

func (ws *bybitMarketWebSocket) handleUPerpTickerError(isWebsocketClosed bool, err error) {
	if err != nil {
		logger.Error("[BybitMarketUPerpWs] Bybit UPerp Tickers Occur Error: %+v", err)
	}
	logger.Warn("[BybitMarketUPerpWs] Bybit UPerp Tickers Ws Will Reconnect Later")
	ws.isStopped = true
}

func (ws *bybitMarketWebSocket) subscribeUPerpTicker(cfg *config.Config, globalContext *mmContext.GlobalContext) {
	ctx := context.Background()
	go func() {
		defer func() {
			logger.Warn("[BybitMarketUPerpWs] Bybit UPerp Tickers Listening Exited.")
		}()

		const (
			baseReconnectDelay = 1 * time.Second
			maxReconnectDelay  = 30 * time.Second
		)

		reconnectDelay := baseReconnectDelay

		for {
			if !ws.isStopped {
				time.Sleep(time.Second * 1)
				continue
			}

			svc, err := ws.connectUPerpWs(globalContext, cfg.BybitWsBaseUrl)
			if err != nil {
				logger.Error("[BybitMarketUPerpWs] Failed to connect Bybit UPerp Ws: %v", err)
				time.Sleep(reconnectDelay)
				// Implement exponential backoff with a cap
				reconnectDelay = min(reconnectDelay*2, maxReconnectDelay)
				continue
			}

			reconnectDelay = baseReconnectDelay

			logger.Info("[BybitMarketUPerpWs] Subscribe Bybit UPerp Successfully")
			ws.isStopped = false

			err = svc.Start(ctx, ws.handleUPerpTickerError)
			if err != nil {
				logger.Error("[BybitMarketUPerpWs] Bybit UPerp Ws error: %v", err)
				ws.isStopped = true
			}
		}
	}()
}

func (ws *bybitMarketWebSocket) connectUPerpWs(globalContext *mmContext.GlobalContext, wsUrl string) (bybit.V5WebsocketPublicServiceI, error) {
	wsClient := bybit.NewWebsocketClient(wsUrl)
	svc, err := wsClient.V5().Public(bybit.CategoryV5Linear)
	if err != nil {
		return nil, fmt.Errorf("[BybitMarketUPerpWs] Start Bybit UPerp Ws Failed: %w", err)
	}

	instIDs := globalContext.InstrumentComposite.GetInstIDs(config.BybitExchange, config.UPerpInstrument)
	for _, instID := range instIDs {
		_, err = svc.SubscribeTicker(bybit.V5WebsocketPublicTickerParamKey{
			Symbol: bybit.SymbolV5(instID),
		}, ws.handleUPerpTickerEvent)

		if err != nil {
			return nil, fmt.Errorf("[BybitMarketUPerpWs] Bybit UPerp Subscribe %s Tickers Failed: %w", instID, err)
		}
		time.Sleep(time.Millisecond * 100)
	}
	return svc, nil
}

func (ws *bybitMarketWebSocket) handleSpotTickerEvent(event bybit.V5WebsocketPublicTickerResponse) error {
	ws.tickerChan <- &event
	return nil
}

func (ws *bybitMarketWebSocket) handleSpotTickerError(isWebsocketClosed bool, err error) {
	if err != nil {
		logger.Error("[BybitMarketSpotWs] Bybit Spot Tickers Occur Error: %+v", err)
	}
	logger.Warn("[BybitMarketSpotWs] Bybit Spot Tickers Ws Will Reconnect In Later")
	ws.isStopped = true
}

func (ws *bybitMarketWebSocket) subscribeSpotTicker(cfg *config.Config, globalContext *mmContext.GlobalContext) {
	ctx := context.Background()
	go func() {
		defer func() {
			logger.Warn("[BybitMarketSpotWs] Bybit Spot Tickers Listening Exited.")
		}()

		const (
			baseReconnectDelay = 1 * time.Second
			maxReconnectDelay  = 30 * time.Second
		)

		reconnectDelay := baseReconnectDelay

		for {
			if !ws.isStopped {
				time.Sleep(time.Second * 1)
				continue
			}

			svc, err := ws.connectSpotWs(globalContext, cfg.BybitWsBaseUrl)
			if err != nil {
				logger.Error("[BybitMarketSpotWs] Failed to connect Bybit Spot Ws: %v", err)
				time.Sleep(reconnectDelay)
				// Implement exponential backoff with a cap
				reconnectDelay = min(reconnectDelay*2, maxReconnectDelay)
				continue
			}

			reconnectDelay = baseReconnectDelay

			logger.Info("[BybitMarketSpotWs] Subscribe Bybit Spot Successfully")
			ws.isStopped = false

			err = svc.Start(ctx, ws.handleSpotTickerError)
			if err != nil {
				logger.Error("[BybitMarketSpotWs] Bybit Spot Ws error: %v", err)
				ws.isStopped = true
			}
		}
	}()
}

func (ws *bybitMarketWebSocket) connectSpotWs(globalContext *mmContext.GlobalContext, wsUrl string) (bybit.V5WebsocketPublicServiceI, error) {
	wsClient := bybit.NewWebsocketClient(wsUrl)
	svc, err := wsClient.V5().Public(bybit.CategoryV5Spot)
	if err != nil {
		return nil, fmt.Errorf("[BybitMarketSpotWs] Start Bybit Spot Ws Failed: %w", err)
	}

	instIDs := globalContext.InstrumentComposite.GetInstIDs(config.BybitExchange, config.SpotInstrument)
	for _, instID := range instIDs {
		_, err = svc.SubscribeTicker(bybit.V5WebsocketPublicTickerParamKey{
			Symbol: bybit.SymbolV5(instID),
		}, ws.handleSpotTickerEvent)

		if err != nil {
			return nil, fmt.Errorf("[BybitMarketSpotWs] Bybit Spot Subscribe %s Tickers Failed: %w", instID, err)
		}
		time.Sleep(time.Millisecond * 100)
	}
	return svc, nil
}
