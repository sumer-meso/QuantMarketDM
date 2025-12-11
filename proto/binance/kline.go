package binance

import (
	"encoding/json"
	"fmt"

	"github.com/sumer-meso/QuantMarketDM/proto"
)

type TrueRangeRatio struct {
	StartTime    int64
	EndTime      int64
	High         float64
	Low          float64
	OBMpAmount   float64
	OBMidPrice   float64
	OBUpdateTime int64
	Value        float64
}

type KlineLite struct {
	StartTime int64
	EndTime   int64
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Volume    float64
}

type KlineDetails struct {
	Kline           KlineLite
	TrueRangeRatios TrueRangeRatio
	Volume          float64
}

// This is Kline model which is used to send out.
type Kline struct {
	Symbol    string
	Range     map[string]KlineDetails
	LocalTime int64
	Source    string
}

func (k *Kline) String() string {
	return fmt.Sprintf("Kline (%s), r:%v, s:%s, l:%d",
		k.Symbol, k.Range, k.Source, k.LocalTime)
}

func (k *Kline) RMQRoutingIdentifier() string {
	return fmt.Sprintf(
		"binance.%s.kline.%s",
		k.Source, k.Symbol,
	)
}

func (k *Kline) RMQDataIdentifier() string {
	return "binance.kline"
}

func (k *Kline) RMQDataStoreTable() string {
	return k.RMQRoutingIdentifier()
}

func (k *Kline) RMQDataStoreIndex() string {
	return "LocalTime:-1"
}

func (k *Kline) RMQEncodeMessage() (*proto.MessageOverRabbitMQ, error) {
	if body, err := json.Marshal(k); err != nil {
		return &proto.MessageOverRabbitMQ{}, err
	} else {
		return &proto.MessageOverRabbitMQ{
			RoutingKey:     k.RMQRoutingIdentifier(),
			DataIdentifier: k.RMQDataIdentifier(),
			StoreTable:     k.RMQDataStoreTable(),
			StoreIndex:     k.RMQDataStoreIndex(),
			Body:           body,
		}, nil
	}
}

func (k *Kline) RMQDecodeMessage(m *proto.MessageOverRabbitMQ) error {
	if m.DataIdentifier != k.RMQDataIdentifier() {
		return proto.NotMatchError{Expected: k.RMQDataIdentifier(), Actual: m.DataIdentifier}
	}
	return json.Unmarshal(m.Body, k)
}

type KlineHandler interface {
	TideHandleBNKline(*Kline) proto.TideRoutable
}

func (k *Kline) TideDispatch(target any) proto.TideRoutable {
	if h, ok := target.(KlineHandler); ok {
		return h.TideHandleBNKline(k)
	}
	return nil
}
