package binance

import "fmt"

type DataType string

const (
	TradeData     DataType = "trade"
	OrderbookData DataType = "orderbook"
	KlineData     DataType = "kline"
)

var AllDataTypes = []DataType{TradeData, OrderbookData, KlineData}

func (t DataType) UrlParam(symbol string, s Source) string {
	switch t {
	case TradeData:
		switch s {
		case Spot:
			return fmt.Sprintf("%s@trade", symbol)
		case Usdt:
			return fmt.Sprintf("%s@trade", symbol)
		case Coin:
			return fmt.Sprintf("%s@trade", symbol)
		}
	case OrderbookData:
		switch s {
		case Spot:
			return fmt.Sprintf("%s@depth@100ms", symbol)
		case Usdt:
			return fmt.Sprintf("%s@depth", symbol)
		case Coin:
			return fmt.Sprintf("%s@depth", symbol)
		}
	}
	return ""
}
