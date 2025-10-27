package wiredata

import (
	"encoding/json"
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

func (ob *OrderBook) String() string {
	return fmt.Sprintf(
		"OrderBook (%s) (%s) Time(%d), LocalTime(%d):"+
			"%d -> %d, Pu: (%d)\n"+
			"Asks(%d):%+v ...\n"+
			"Bids(%d):%+v ...\n",
		ob.Symbol, ob.Event, ob.Time, ob.LocalTime,
		ob.FirstUpdateID, ob.LastUpdateID, ob.PreviousUpdateID,
		len(ob.Asks), ob.Asks[:min(len(ob.Asks), 50)],
		len(ob.Bids), ob.Bids[:min(len(ob.Bids), 50)],
	)
}

const obIndexSpecInRMQ = "{Time:-1,LocalTime:-1}"

func (ob *OrderBook) RMQRoutingIdentifier() string {
	return fmt.Sprintf(
		"binance.%s.orderbook.%s.%s",
		ob.Source, ob.Symbol, obIndexSpecInRMQ,
	)
}

func (ob *OrderBook) RMQDataIdentifier() string {
	return "binance.orderbook"
}

func (ob *OrderBook) RMQEncodeMessage() (MessageOverRabbitMQ, error) {
	if body, err := json.Marshal(ob); err != nil {
		return MessageOverRabbitMQ{}, err
	} else {
		return MessageOverRabbitMQ{
			RoutingKey:     ob.RMQRoutingIdentifier(),
			DataIdentifier: ob.RMQDataIdentifier(),
			Body:           body,
		}, nil
	}
}

func (ob *OrderBook) RMQDecodeMessage(m MessageOverRabbitMQ) error {
	if m.DataIdentifier != ob.RMQDataIdentifier() {
		return NotMatchError{Expected: ob.RMQDataIdentifier(), Actual: m.DataIdentifier}
	}
	return json.Unmarshal(m.Body, ob)
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
