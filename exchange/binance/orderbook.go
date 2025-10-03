package binance

import (
	"fmt"
)

type OrderBookEntry struct {
	Price    float64
	Quantity float64
}

type MidPriceEntry struct {
	Amount   float64
	AskPrice float64
	BidPrice float64
	Value    float64
}

type SpreadEntry struct {
	MidPriceAmount float64
	SpreadAmount   float64
	AskPrice       float64
	BidPrice       float64
	Value          float64
}

type ImbalanceEntry struct {
	Layer   int
	BidSize float64
	AskSize float64
	Value   float64
}

type OrderBook struct {
	Symbol           string
	Event            string
	Time             int64
	LastUpdateID     int64
	FirstUpdateID    int64
	PreviousUpdateID int64
	Bids             []OrderBookEntry
	Asks             []OrderBookEntry
	LocalTime        int64
	Source           string
	MidPrice         MidPriceEntry
	Spread           SpreadEntry
	Imbalances       []ImbalanceEntry
}

func (o OrderBook) String() string {
	return fmt.Sprintf(
		"OrderBook (%s) (%s) Time(%d), LocalTime(%d):"+
			"%d -> %d, Pu: (%d)\n"+
			"Asks(%d):%+v ...\n"+
			"Bids(%d):%+v ...\n",
		o.Symbol, o.Event, o.Time, o.LocalTime,
		o.FirstUpdateID, o.LastUpdateID, o.PreviousUpdateID,
		len(o.Asks), o.Asks[:min(len(o.Asks), 50)],
		len(o.Bids), o.Bids[:min(len(o.Bids), 50)],
	)
}

func (o OrderBook) RoutingKey() string {
	return fmt.Sprintf(
		"binance.%s.orderbook.%s",
		o.Source, o.Symbol,
	)
}

type WsDepthEvent struct {
	Event            string     `json:"e"`
	Time             int64      `json:"E"`
	Symbol           string     `json:"s"`
	FirstUpdateID    int64      `json:"U"`
	LastUpdateID     int64      `json:"u"`
	PreviousUpdateID int64      `json:"pu"`
	Bids             [][]string `json:"b"`
	Asks             [][]string `json:"a"`
}
