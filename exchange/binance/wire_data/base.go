package wiredata

type RMQRoutingIdentifier interface{ RMQRoutingIdentifier() string }
type RMQDataIdentifier interface{ RMQDataIdentifier() string }

var _ = []interface {
	RMQRoutingIdentifier
	RMQDataIdentifier
}{
	(*Trade)(nil),
	(*Kline)(nil),
	(*OrderBook)(nil),
}
