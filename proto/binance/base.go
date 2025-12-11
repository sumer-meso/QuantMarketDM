package binance

import (
	"github.com/sumer-meso/QuantMarketDM/proto"
)

type WSEventBase struct {
	Event string `json:"e"`
	Time  int64  `json:"E"`
}

type EventBase struct {
	Event string
	Time  int64
}

type LocalBase struct {
	LocalTime int64
	Source    string
}

var _ = []interface {
	proto.RMQIdentifier
	proto.RMQDataStorage
	proto.RMQSerializationOnWire
	proto.Routable
}{
	(*Trade)(nil),
	(*Kline)(nil),
	(*OrderBook)(nil),
	(*TradeLite)(nil),
	(*AccountUpdate)(nil),
	(*OrderUpdate)(nil),
}
