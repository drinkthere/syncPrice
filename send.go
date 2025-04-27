package main

import (
	zmq "github.com/pebbe/zmq4"
	"google.golang.org/protobuf/proto"
	"os"
	"sort"
	"syncPrice/config"
	"syncPrice/protocol/pb"
	"syncPrice/utils/logger"
	"time"
)

func startSendMedianPrice() {
	mpChan := make(chan *string)
	startZmq(globalConfig.MedianPriceZMQ, mpChan)
	sendMedianPrice(mpChan)
}

func startZmq(zmqIPC string, msgCh chan *string) {
	go func() {
		defer func() {
			logger.Warn("[StartZmq] %s Pub Service Listening Exited.", zmqIPC)
		}()

		logger.Info("[StartZmq] %s Start Pub Service.", zmqIPC)

		var ctx *zmq.Context
		var pub *zmq.Socket
		var err error
		isPubStopped := true
		for {
			if isPubStopped {
				ctx, err = zmq.NewContext()
				if err != nil {
					logger.Error("[StartZmq] %s New Context Error: %s", zmqIPC, err.Error())
					time.Sleep(time.Second * 1)
					continue
				}

				pub, err = ctx.NewSocket(zmq.PUB)
				if err != nil {
					logger.Error("[StartZmq] %s New Socket Error: %s", zmqIPC, err.Error())
					ctx.Term()
					time.Sleep(time.Second * 1)
					continue
				}

				err = pub.Bind("ipc://" + zmqIPC)
				if err != nil {
					logger.Error("[StartZmq] Bind to Local ZMQ %s Error: %s", zmqIPC, err.Error())
					pub.Close() // 关闭套接字
					ctx.Term()  // 释放上下文
					time.Sleep(time.Second * 1)
					continue
				}

				// 修改 IPC 文件的权限为 0666（所有用户可读写）
				err = os.Chmod(zmqIPC, 0666)
				if err != nil {
					logger.Error("[StartZmq] Set %s Permission to 0666 Failed Error: %s", zmqIPC, err.Error())
					pub.Close() // 关闭套接字
					ctx.Term()  // 释放上下文
					time.Sleep(time.Second * 1)
					continue
				}
				isPubStopped = false
			}

			select {
			case msg := <-msgCh:
				_, err = pub.Send(*msg, 0)
				if err != nil {
					logger.Warn("[StartZmq] %s Error sending Median Price Data %s: %v", zmqIPC, *msg, err)
					isPubStopped = true
					pub.Close()
					ctx.Term()
					time.Sleep(time.Second * 5)
					continue
				}
			}
		}
	}()
}

func sendMedianPrice(msgCh chan *string) {
	go func() {
		ticker := time.NewTicker(time.Duration(globalConfig.SendIntervalMs) * time.Millisecond)
		for {
			_ = <-ticker.C

			bnUPerpWrapper := globalContext.TickerComposite.GetTickerWrapper(config.BinanceExchange, config.UPerpInstrument)
			bnSpotWrapper := globalContext.TickerComposite.GetTickerWrapper(config.BinanceExchange, config.SpotInstrument)
			bbUPerpWrapper := globalContext.TickerComposite.GetTickerWrapper(config.BybitExchange, config.UPerpInstrument)
			bbSpotWrapper := globalContext.TickerComposite.GetTickerWrapper(config.BybitExchange, config.SpotInstrument)
			okxUPerpWrapper := globalContext.TickerComposite.GetTickerWrapper(config.OkxExchange, config.UPerpInstrument)
			okxSpotWrapper := globalContext.TickerComposite.GetTickerWrapper(config.OkxExchange, config.SpotInstrument)

			var medianPriceList []*pb.MedianPrice
			coins := []string{"BTC", "ETH"}
			for _, coin := range coins {

				instID := coin + "USDT"
				bnUPerpTicker := bnUPerpWrapper.GetTicker(instID)
				bnSpotTicker := bnSpotWrapper.GetTicker(instID)
				bbUPerpTicker := bbUPerpWrapper.GetTicker(instID)
				bbSpotTicker := bbSpotWrapper.GetTicker(instID)

				instID = coin + "-USDT-SWAP"
				okxUPerpTicker := okxUPerpWrapper.GetTicker(instID)

				instID = coin + "-USDT"
				okxSpotTicker := okxSpotWrapper.GetTicker(instID)

				prices := []float64{bnUPerpTicker.BidPx, bnSpotTicker.BidPx, bbUPerpTicker.BidPx, bbSpotTicker.BidPx, okxUPerpTicker.BidPx, okxSpotTicker.BidPx}
				medianPrice := calculateMedian(prices)
				mp := pb.MedianPrice{
					Coin:  coin,
					Price: medianPrice,
				}
				medianPriceList = append(medianPriceList, &mp)
			}
			if len(medianPriceList) > 0 {
				mpList := &pb.MedianPriceList{
					Prices: medianPriceList,
				}
				data, err := proto.Marshal(mpList)
				if err != nil {
					logger.Warn("[SendMediaPrice] Error marshaling price list: %v", err)
					continue
				}
				dataStr := string(data)
				msgCh <- &dataStr
			}
		}
	}()
}
func calculateMedian(data []float64) float64 {
	if len(data) == 0 {
		logger.Fatal("[SendMedianPrice] ticker prices data is empty")
	}

	// 2. 对切片进行排序
	sort.Float64s(data)

	// 3. 计算中位数
	n := len(data)
	if n%2 == 1 {
		// 奇数长度，返回中间的值
		return data[n/2]
	}
	// 偶数长度，返回中间两个数的平均值
	return (data[n/2-1] + data[n/2]) / 2
}
