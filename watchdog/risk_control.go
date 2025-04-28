package watchdog

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os/exec"
	"strconv"
	"strings"
	"syncPrice/config"
	"syncPrice/context"
	"syncPrice/utils/logger"
	"time"
)

func StartRiskControl(globalConfig *config.Config, globalContext *context.GlobalContext) {
	statErrorsInLog(globalConfig, globalContext, time.Minute*1)
	logger.Info("[Watchdog] Start to Stat Errors Num In Log")

	checkPriceUpdateTime(globalConfig, globalContext, time.Second*10)
	logger.Info("[Watchdog] Start to Check Ticker Update in Time")
}

// 检查错误日志是不是超过阈值
func statErrorsInLog(globalConfig *config.Config, globalContext *context.GlobalContext, interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)

			// 获取服务器设置的时区
			now := time.Now()
			local, err := time.LoadLocation("Local")
			if err != nil {
				logger.Error("[ErrorWatchdog] Load timezone failed, message is %s", err.Error())
				continue
			}

			lastMinute := now.In(local).Format("2006-01-02T15:04")
			count := getErrorCountFromLog(globalConfig.LogPath, lastMinute)

			if count > globalConfig.MaxErrorsPerMinute {
				msg := "ERROR logs nums over max num, stop syncing median price."
				logger.Error("[ErrorWatchdog] %s", msg)
				globalContext.TelegramBot.Send(tgbotapi.NewMessage(globalConfig.TgChatID, msg))
			}
			logger.Info("[ErrorWatchdog] Finish Stat Errors Num In Log: Num=%d", count)
		}
	}()
}

func getErrorCountFromLog(logFile string, lastMinute string) int64 {
	shell := fmt.Sprintf("/usr/bin/grep '%s' %s | grep ERROR | grep -v \"10 seconds\" | wc -l", lastMinute, logFile)

	// 通过 grep 获取 Error 日志出现次数
	cmd := exec.Command("/bin/sh", "-c", shell)
	countRaw, err := cmd.Output()
	if err != nil {
		logger.Error("[ErrorWatchdog] Failed to stat ERROR logs, message is %s", err.Error())
		return 0
	}

	// 默认执行结果里面有"\n"，需要去掉
	countStr := strings.Trim(string(countRaw), "\n")
	countStr = strings.Trim(countStr, " ")
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		logger.Error("[ErrorWatchdog] CheckErrors failed to convert count to int, count is %s, message is %s", countStr, err.Error())
		return 0
	}

	return count
}

// 检查价格是否没有更新
func checkPriceUpdateTime(
	globalConfig *config.Config,
	globalContext *context.GlobalContext,
	interval time.Duration) {

	initInterval := interval
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Recovered from panic in price update check: %v", r)
			}
		}()
		for {
			time.Sleep(initInterval)
			exchanges := []config.Exchange{config.BinanceExchange, config.OkxExchange, config.BybitExchange}
			instTypes := []config.InstrumentType{config.UPerpInstrument, config.SpotInstrument}

			var expiredWs []string
			for _, exchange := range exchanges {
				for _, instType := range instTypes {
					instIDs := globalContext.InstrumentComposite.GetInstIDs(exchange, instType)
					tickerWrapper := globalContext.TickerComposite.GetTickerWrapper(exchange, instType)
					tickerExpired := false
					for _, instID := range instIDs {
						var timeoutThresholdSec int64 = 50
						ticker := tickerWrapper.GetTicker(instID)
						if ticker == nil || ticker.IsExpired(timeoutThresholdSec*1000*1000) {
							logger.Error("[MedianPrice] %s %s %s Ticker %+v Did Not Update For %d seconds",
								exchange, instType, instID, ticker, timeoutThresholdSec)
							tickerExpired = true
						}
					}

					if tickerExpired {
						expiredWs = append(expiredWs, fmt.Sprintf("%s_%s", exchange, instType))
					}
				}
			}

			const maxInterval = 10 * time.Minute
			if len(expiredWs) > 0 {
				msg := strings.Join(expiredWs, "&") + " ws ticker update delay"
				logger.Error("[MedianPrice] " + msg)
				// globalContext.TelegramBot.Send(tgbotapi.NewMessage(globalConfig.TgChatID, msg))
				if initInterval < maxInterval {
					initInterval = initInterval * 2
				}
			} else {
				initInterval = interval
			}
		}
	}()
}
