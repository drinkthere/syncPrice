package context

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"syncPrice/config"
	"syncPrice/container"
	"syncPrice/utils/logger"
)

type GlobalContext struct {
	TelegramBot         *tgbotapi.BotAPI // 电报机器人
	InstrumentComposite *container.InstrumentComposite
	TickerComposite     *container.TickerComposite
}

func (context *GlobalContext) Init(globalConfig *config.Config) {
	// 初始化电报机器人
	context.initTelegramBot(globalConfig)

	// 初始化交易对数据
	context.initInstrumentComposite(globalConfig)

	// 初始化ticker数据
	context.initTickerComposite()
}

func (context *GlobalContext) initTelegramBot(globalConfig *config.Config) {
	// 初始化 telegramBot
	bot, err := tgbotapi.NewBotAPI(globalConfig.TgBotToken)
	if err != nil {
		logger.Error("[Context] Init Telegram Bot Failed, Error Is %s", err.Error())
		os.Exit(1)
	}
	context.TelegramBot = bot
}

func (context *GlobalContext) initInstrumentComposite(globalConfig *config.Config) {
	instrumentComposite := container.NewInstrumentComposite(globalConfig)
	context.InstrumentComposite = instrumentComposite
}

func (context *GlobalContext) initTickerComposite() {
	context.TickerComposite = container.NewTickerComposite()
}
