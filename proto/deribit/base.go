package deribit

import (
	"github.com/sumer-meso/QuantMarketDM/proto"
	user "github.com/sumer-meso/QuantMarketDM/proto/deribit/dbuser"
)

var _ = []interface {
	proto.RMQIdentifier
	proto.RMQDataStorage
	proto.RMQSerializationOnWire
	proto.Routable
}{
	(*user.Portfolio)(nil),
	(*user.Position)(nil),
	(*user.Order)(nil),
	(*user.Trade)(nil),
	(*OrderBook)(nil),
	(*Trade)(nil),
}
