package config

type (
	RiskType        int
	Exchange        string
	InstrumentType  string
	OrderSide       string
	OrderStatus     string
	TickerSource    string
	DeltaChangeType string
)

const (
	// NoRisk 可以挂单，其他RiskType暂停挂单。1表示出Error误消息太多，2表示处于结算时间，3表示系统在退出中国， 4表示暂停等待价格更新， 5表示账户的FeeRate有问题 6表示收到了外部暂停下单的信号
	NoRisk                         = RiskType(iota)
	FatalErrorRisk                 = NoRisk + 1
	SettlingRisk                   = FatalErrorRisk + 1
	ExitRisk                       = SettlingRisk + 1
	Price10sNotUpdateRisk          = ExitRisk + 1
	FeeRateWrongRisk               = Price10sNotUpdateRisk + 1
	SignalRisk                     = FeeRateWrongRisk + 1
	MaxUnRealisedProfitRisk        = SignalRisk + 1
	OrderNumExceedMaxThresholdRisk = MaxUnRealisedProfitRisk + 1
	RestNWsPositionDiffRisk        = OrderNumExceedMaxThresholdRisk + 1
	OrderWsNotUpdateRisk           = RestNWsPositionDiffRisk + 1

	BinanceExchange = Exchange("Binance")
	OkxExchange     = Exchange("Okx")
	BybitExchange   = Exchange("Bybit")

	UnknownInstrument = InstrumentType("UNKNOWN")
	SpotInstrument    = InstrumentType("SPOT")
	UPerpInstrument   = InstrumentType("UPERP")
	MarginInstrument  = InstrumentType("MARGIN")

	BuyOrderSide  = OrderSide("buy")
	SellOrderSide = OrderSide("sell")

	OrderCreate          = OrderStatus("create")         // 发起创建请求，但不确定订单是否创建成功
	OrderLive            = OrderStatus("live")           // 新订单创建成功
	OrderCanceling       = OrderStatus("OrderCanceling") // 发起取消请求，但不确定订单是否取消成功
	OrderPartiallyFilled = OrderStatus("partially_filled")
	OrderFilled          = OrderStatus("filled")
	OrderCancel          = OrderStatus("canceled")

	RapidMarketTickerSource = TickerSource("RapidMarket")
	WebSocketTickerSource   = TickerSource("WebSocket")

	NoDeltaChangeType   = DeltaChangeType("NoChange")
	BuyDeltaChangeType  = DeltaChangeType("BuyChange")
	SellDeltaChangeType = DeltaChangeType("SellChange")
)
