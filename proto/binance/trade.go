package binance

import (
	"encoding/json"
	"fmt"

	"github.com/sumer-meso/QuantMarketDM/proto"
)

type Trade struct {
	Event         string
	Time          int64
	Symbol        string
	TradeID       int64
	Price         float64
	Quantity      float64
	BuyerOrderId  int64
	SellerOrderId int64
	TradeTime     int64
	IsBuyerMaker  bool
	Placeholder   bool
	LocalTime     int64
	Source        string
}

func (t *Trade) RMQRoutingIdentifier() string {
	return fmt.Sprintf("binance.%s.trade.%s", t.Source, t.Symbol)
}

func (k *Trade) RMQDataIdentifier() string {
	return "binance.trade"
}

func (k *Trade) RMQDataStoreTable() string {
	return k.RMQRoutingIdentifier()
}

func (k *Trade) RMQDataStoreIndex() string {
	return "LocalTime:-1,TradeTime:-1"
}

func (t *Trade) RMQEncodeMessage() (*proto.MessageOverRabbitMQ, error) {
	if body, err := json.Marshal(t); err != nil {
		return &proto.MessageOverRabbitMQ{}, err
	} else {
		return &proto.MessageOverRabbitMQ{
			RoutingKey:     t.RMQRoutingIdentifier(),
			DataIdentifier: t.RMQDataIdentifier(),
			StoreTable:     t.RMQDataStoreTable(),
			StoreIndex:     t.RMQDataStoreIndex(),
			Body:           body,
		}, nil
	}
}

func (t *Trade) RMQDecodeMessage(m *proto.MessageOverRabbitMQ) error {
	if m.DataIdentifier != t.RMQDataIdentifier() {
		return proto.NotMatchError{Expected: t.RMQDataIdentifier(), Actual: m.DataIdentifier}
	}
	return json.Unmarshal(m.Body, t)
}

type WsTradeEvent struct {
	Time          int64  `json:"E"`
	Event         string `json:"e"`
	Symbol        string `json:"s"`
	TradeID       int64  `json:"t"`
	Price         string `json:"p"`
	Quantity      string `json:"q"`
	BuyerOrderId  int64  `json:"b"`
	SellerOrderId int64  `json:"a"`
	TradeTime     int64  `json:"T"`
	IsBuyerMaker  bool   `json:"m"`
	Placeholder   bool   `json:"M"` // add this field to avoid case insensitive unmarshaling
}
