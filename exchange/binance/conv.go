package binance

import (
	"github.com/sumer-meso/QuantMarketDM/common"
	"github.com/sumer-meso/QuantMarketDM/tide"
)

func (t Trade) ToMarketTrade() (tide.MarketTrade, error) {
	side := common.SideBuy
	if t.BuyerOrderId > 0 {
		side = common.SideSell
	}

	return tide.MarketTrade{
		Exchange: "binance",
		Symbol:   t.Symbol,
		TradeID:  t.TradeID,
		TsUnixMs: t.LocalTime,
		Price:    t.Price,
		Quantity: t.Quantity,
		Side:     side,
		IsMaker:  t.IsBuyerMaker,
	}, nil
}
