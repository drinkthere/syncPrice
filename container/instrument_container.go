package container

import (
	"fmt"
	"github.com/drinkthere/bybit"
	"log"
	"strconv"
	"sync"
	"syncPrice/client"
	"syncPrice/config"
	"syncPrice/utils"
	"syncPrice/utils/logger"
)

type BybitSpotPrecision struct {
	PricePrecision int
	PriceTick      float64
	SizePrecision  int
	MinSize        float64
}
type InstrumentComposite struct {
	InstIDsMap                sync.Map  // 使用 sync.Map 来提高并发性能
	BybitSpotInstPrecisionMap *sync.Map // BYBIT 现货的price和size的精度，计算ticker价格时会用到
}

func NewInstrumentComposite(globalConfig *config.Config) *InstrumentComposite {
	composite := &InstrumentComposite{
		InstIDsMap: sync.Map{},
	}

	// 初始化 InstIDsMap
	bnUPerpKey := fmt.Sprintf("%s_%s", config.BinanceExchange, config.UPerpInstrument)
	composite.InstIDsMap.Store(bnUPerpKey, globalConfig.BinanceUPerpInstIDs)
	bnSpotKey := fmt.Sprintf("%s_%s", config.BinanceExchange, config.SpotInstrument)
	composite.InstIDsMap.Store(bnSpotKey, globalConfig.BinanceSpotInstIDs)

	okxUPerpKey := fmt.Sprintf("%s_%s", config.OkxExchange, config.UPerpInstrument)
	composite.InstIDsMap.Store(okxUPerpKey, globalConfig.OkxUPerpInstIDs)
	okxSpotKey := fmt.Sprintf("%s_%s", config.OkxExchange, config.SpotInstrument)
	composite.InstIDsMap.Store(okxSpotKey, globalConfig.OkxSpotInstIDs)

	bybitUPerpKey := fmt.Sprintf("%s_%s", config.BybitExchange, config.UPerpInstrument)
	composite.InstIDsMap.Store(bybitUPerpKey, globalConfig.BybitUPerpInstIDs)
	bybitSpotKey := fmt.Sprintf("%s_%s", config.BybitExchange, config.SpotInstrument)
	composite.InstIDsMap.Store(bybitSpotKey, globalConfig.BybitSpotInstIDs)

	// 初始化bybit 现货价格和数量的精度
	precisionMap := getBybitSpotInstPrecision(globalConfig)
	composite.BybitSpotInstPrecisionMap = precisionMap

	return composite
}

func (composite *InstrumentComposite) GetInstIDs(exchange config.Exchange, instType config.InstrumentType) []string {
	key := fmt.Sprintf("%s_%s", exchange, instType)
	if instIDs, ok := composite.InstIDsMap.Load(key); ok {
		return instIDs.([]string)
	}
	return nil
}

func getBybitSpotInstPrecision(globalConfig *config.Config) *sync.Map {
	precisionMap := new(sync.Map)
	exchangeClient, ok := client.NewExchangeClient(config.BybitExchange, globalConfig)
	if !ok {
		logger.Fatal("Error to Init Bybit Client")
	}

	if bybitClient, ok := exchangeClient.(*client.BybitClient); ok {
		resp, err := bybitClient.Client.V5().Market().GetInstrumentsInfo(bybit.V5GetInstrumentsInfoParam{
			Category: bybit.CategoryV5Spot,
		})
		if err != nil {
			logger.Fatal("Error to Get Spot Instruments")
		}

		for _, spot := range resp.Result.Spot.List {
			instID := string(spot.Symbol)
			if utils.InArray(instID, globalConfig.BybitSpotInstIDs) {
				pricePrecision := utils.GetDecimals(spot.PriceFilter.TickSize)
				priceTick, _ := strconv.ParseFloat(spot.PriceFilter.TickSize, 64)
				sizePrecision := utils.GetDecimals(spot.LotSizeFilter.MinOrderQty)
				minSize, _ := strconv.ParseFloat(spot.LotSizeFilter.MinOrderQty, 64)
				precisionMap.Store(instID, BybitSpotPrecision{
					PricePrecision: pricePrecision,
					PriceTick:      priceTick,
					SizePrecision:  sizePrecision,
					MinSize:        minSize,
				})
			}
		}
	} else {
		log.Fatal("Failed to cast client to BybitClient")
	}
	return precisionMap
}
