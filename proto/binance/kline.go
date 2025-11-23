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

// This is Kline model which is used to send out.
type Kline struct {
	Symbol          string
	TrueRangeRatios map[string]TrueRangeRatio
	Volumes         map[string]float64
	LocalTime       int64
	Source          string
}

func (k *Kline) String() string {
	return fmt.Sprintf("Kline (%s), t:%v, v:%v, s:%s, l:%d",
		k.Symbol, k.TrueRangeRatios, k.Volumes, k.Source, k.LocalTime)
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
		return NotMatchError{Expected: k.RMQDataIdentifier(), Actual: m.DataIdentifier}
	}
	return json.Unmarshal(m.Body, k)
}

// WsKlineEvent define websocket kline event
type WsKlineEvent struct {
	Event  string  `json:"e"`
	Time   int64   `json:"E"`
	Symbol string  `json:"s"`
	Kline  WsKline `json:"k"`
}

// WsKline define websocket kline
type WsKline struct {
	StartTime            int64  `json:"t"`
	EndTime              int64  `json:"T"`
	Symbol               string `json:"s"`
	Interval             string `json:"i"`
	FirstTradeID         int64  `json:"f"`
	LastTradeID          int64  `json:"L"`
	Open                 string `json:"o"`
	Close                string `json:"c"`
	High                 string `json:"h"`
	Low                  string `json:"l"`
	Volume               string `json:"v"`
	TradeNum             int64  `json:"n"`
	IsFinal              bool   `json:"x"`
	QuoteVolume          string `json:"q"`
	ActiveBuyVolume      string `json:"V"`
	ActiveBuyQuoteVolume string `json:"Q"`
}
