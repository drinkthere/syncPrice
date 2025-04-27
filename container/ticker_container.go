package container

import (
	"fmt"
	"sync"
	"syncPrice/config"
	"time"
)

type TickerMessage struct {
	Exchange config.Exchange
	InstType config.InstrumentType
	Ticker
}

type Ticker struct {
	InstID          string
	BidPx           float64 // 买1价
	BidSz           float64 // 买1量
	AskPx           float64 // 卖1价
	AskSz           float64 // 卖1量
	UpdateTimeMicro int64   // 更新时间（微秒）
	UpdateTimeMs    int64   // 更新时间（毫秒）
	UpdateID        int64   // 更新ID
}

func (t *Ticker) init(instID string) {
	t.InstID = instID
	t.BidPx, t.BidSz, t.AskPx, t.AskSz, t.UpdateTimeMicro, t.UpdateID, t.UpdateTimeMs = 0.0, 0.0, 0.0, 0.0, 0, 0, 0
}

func (t *Ticker) IsExpired(micro int64) bool {
	return time.Now().UnixMicro()-t.UpdateTimeMicro > micro
}

func (t *Ticker) updateTicker(message TickerMessage) bool {
	if (message.Exchange == config.BinanceExchange && message.InstType == config.SpotInstrument && t.UpdateID >= message.UpdateID) ||
		(t.UpdateTimeMs >= message.UpdateTimeMs) {
		return false // message is expired
	}

	if message.AskPx > 0.0 && message.AskSz > 0.0 {
		t.AskPx = message.AskPx
		t.AskSz = message.AskSz
	}
	if message.BidPx > 0.0 && message.BidSz > 0.0 {
		t.BidPx = message.BidPx
		t.BidSz = message.BidSz
	}

	t.UpdateID = message.UpdateID
	t.UpdateTimeMs = message.UpdateTimeMs
	t.UpdateTimeMicro = time.Now().UnixMicro()
	return true
}

type TickerWrapper struct {
	Exchange  config.Exchange
	InstType  config.InstrumentType
	tickerMap sync.Map // 使用 sync.Map
}

func newTickerWrapper(exchange config.Exchange, instType config.InstrumentType) *TickerWrapper {
	return &TickerWrapper{
		Exchange:  exchange,
		InstType:  instType,
		tickerMap: sync.Map{},
	}
}

func (wrapper *TickerWrapper) GetTicker(instID string) *Ticker {
	if ticker, ok := wrapper.tickerMap.Load(instID); ok {
		return ticker.(*Ticker) // 类型断言
	}
	return nil
}

func (wrapper *TickerWrapper) UpdateTicker(message TickerMessage) bool {
	if wrapper.Exchange != message.Exchange || wrapper.InstType != message.InstType {
		return false
	}

	ticker, _ := wrapper.tickerMap.Load(message.InstID)
	var newTicker *Ticker

	if ticker == nil {
		newTicker = &Ticker{}
		newTicker.init(message.InstID)
		newTicker.updateTicker(message)
		wrapper.tickerMap.Store(message.InstID, newTicker)
	} else {
		newTicker = ticker.(*Ticker) // 类型断言
		newTicker.updateTicker(message)
		wrapper.tickerMap.Store(message.InstID, newTicker)
	}
	return true
}

type TickerComposite struct {
	tickerWrapperMap *sync.Map // 使用 sync.Map
}

func NewTickerComposite() *TickerComposite {
	composite := &TickerComposite{
		tickerWrapperMap: new(sync.Map),
	}

	// 初始化 InstIDsMap
	bnUPerpKey := fmt.Sprintf("%s_%s", config.BinanceExchange, config.UPerpInstrument)
	composite.tickerWrapperMap.Store(bnUPerpKey, newTickerWrapper(config.BinanceExchange, config.UPerpInstrument))
	bnSpotKey := fmt.Sprintf("%s_%s", config.BinanceExchange, config.SpotInstrument)
	composite.tickerWrapperMap.Store(bnSpotKey, newTickerWrapper(config.BinanceExchange, config.SpotInstrument))

	okxUPerpKey := fmt.Sprintf("%s_%s", config.OkxExchange, config.UPerpInstrument)
	composite.tickerWrapperMap.Store(okxUPerpKey, newTickerWrapper(config.OkxExchange, config.UPerpInstrument))
	okxSpotKey := fmt.Sprintf("%s_%s", config.OkxExchange, config.SpotInstrument)
	composite.tickerWrapperMap.Store(okxSpotKey, newTickerWrapper(config.OkxExchange, config.SpotInstrument))

	bbUPerpKey := fmt.Sprintf("%s_%s", config.BybitExchange, config.UPerpInstrument)
	composite.tickerWrapperMap.Store(bbUPerpKey, newTickerWrapper(config.BybitExchange, config.UPerpInstrument))
	bbSpotKey := fmt.Sprintf("%s_%s", config.BybitExchange, config.SpotInstrument)
	composite.tickerWrapperMap.Store(bbSpotKey, newTickerWrapper(config.BybitExchange, config.SpotInstrument))

	return composite
}

func (c *TickerComposite) GetTickerWrapper(exchange config.Exchange, instType config.InstrumentType) *TickerWrapper {
	key := fmt.Sprintf("%s_%s", exchange, instType)
	if wrapper, ok := c.tickerWrapperMap.Load(key); ok {
		return wrapper.(*TickerWrapper)
	}
	return nil
}
