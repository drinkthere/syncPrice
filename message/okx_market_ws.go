package message

import (
	"github.com/drinkthere/okx/events"
	"github.com/drinkthere/okx/events/public"
	wsRequestPublic "github.com/drinkthere/okx/requests/ws/public"
	"strings"
	"syncPrice/client"
	"syncPrice/config"
	"syncPrice/context"
	"syncPrice/utils/logger"
	"time"
)

func StartOkxMarketWs(globalConfig *config.Config, globalContext *context.GlobalContext,
	instType config.InstrumentType, tickerChan chan *public.Tickers) {

	logger.Info("[Okx%sTickerWs] Start Listen Tickers", instType)
	go func() {
		defer func() {
			if rc := recover(); rc != nil {
				logger.Error("[Okx%sTickerWs] Recovered from panic: %v", instType, rc)
			}

			logger.Warn("[Okx%sTickerWs] Ticker Listening Exited.", instType)
		}()
		for {
		ReConnect:
			errChan := make(chan *events.Error)
			subChan := make(chan *events.Subscribe)
			uSubChan := make(chan *events.Unsubscribe)
			loginCh := make(chan *events.Login)
			successCh := make(chan *events.Success)

			var okxClient = client.OkxClient{}
			okxClient.Init(globalConfig)

			okxClient.Client.Ws.SetChannels(errChan, subChan, uSubChan, loginCh, successCh)
			instIDs := globalContext.InstrumentComposite.GetInstIDs(config.OkxExchange, instType)
			for _, instID := range instIDs {
				err := okxClient.Client.Ws.Public.Tickers(wsRequestPublic.Tickers{
					InstID: instID,
				}, tickerChan)

				if err != nil {
					logger.Fatal("[Okx%sTickerWs] Fail To Listen Ticker For %s, %s", instType, instID, err.Error())
				} else {
					logger.Info("[Okx%sTickerWs] Ticker WebSocket Has Established For %s", instType, instID)
				}
				time.Sleep(100 * time.Millisecond)
			}

			for {
				select {
				case sub := <-subChan:
					channel, _ := sub.Arg.Get("channel")
					logger.Info("[Okx%sTickerWs] Subscribe %s", instType, channel)
				case err := <-errChan:
					if strings.Contains(err.Msg, "i/o timeout") {
						logger.Warn("[Okx%sTickerWs] Error occurred %s, Will reconnect after 1 second.", instType, err.Msg)
						time.Sleep(time.Second * 1)
						goto ReConnect
					}
					logger.Error("[Okx%sTickerWs] Occur Some Error %+v", instType, err)
					for _, datum := range err.Data {
						logger.Error("[Okx%sTickerWs] Error Data %+v", instType, datum)
					}

				case s := <-successCh:
					logger.Info("[Okx%sTickerWs] Receive Success: %+v", instType, s)
				case b := <-okxClient.Client.Ws.DoneChan:
					logger.Info("[Okx%sTickerWs] End %v", instType, b)
					// 暂停一秒再跳出，避免异常时频繁发起重连
					logger.Warn("[Okx%sTickerWs] Will Reconnect WebSocket After 1 Second", instType)
					time.Sleep(time.Second * 1)
					goto ReConnect
				}
			}
		}
	}()
}
