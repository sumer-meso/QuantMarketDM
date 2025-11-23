package deribit

import (
	"github.com/sumer-meso/QuantMarketDM/proto"
	"github.com/sumer-meso/QuantMarketDM/proto/deribit/user"
)

var _ = []interface {
	proto.RMQIdentifier
	proto.RMQDataStorage
	proto.RMQSerializationOnWire
}{
	(*user.Portfolio)(nil),
	(*user.Position)(nil),
	(*user.Order)(nil),
	(*user.Trade)(nil),
}
