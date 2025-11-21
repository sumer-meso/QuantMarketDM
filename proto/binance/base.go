package binance

import (
	"fmt"

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

type NotMatchError struct {
	Expected string // e.g., "orderbook"
	Actual   string // e.g., "trade"
}

func (e NotMatchError) Error() string {
	return fmt.Sprintf("data identifier not match, expected: %s, actual: %s", e.Expected, e.Actual)
}

var _ = []interface {
	proto.RMQIdentifier
	proto.RMQDataStorage
	proto.RMQSerializationOnWire
}{
	(*Trade)(nil),
	(*Kline)(nil),
	(*OrderBook)(nil),
	(*TradeLite)(nil),
	(*AccountUpdate)(nil),
	(*OrderUpdate)(nil),
}
