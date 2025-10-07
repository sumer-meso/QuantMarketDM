package data

import "fmt"

type Type string

const (
	Trade     Type = "trade"
	Orderbook Type = "orderbook"
	Kline     Type = "kline"
)

var AllDataTypes = []Type{Trade, Orderbook, Kline}

func (t Type) UrlParam(symbol string, s Source) string {
	switch t {
	case Trade:
		switch s {
		case Spot:
			return fmt.Sprintf("%s@trade", symbol)
		case Usdt:
			return fmt.Sprintf("%s@trade", symbol)
		case Coin:
			return fmt.Sprintf("%s@trade", symbol)
		}
	case Orderbook:
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
