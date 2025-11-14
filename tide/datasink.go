package tide

import (
	"context"

	bWiredata "github.com/sumer-meso/QuantMarketDM/exchange/binance/proto"
)

// DataSink receives decoded data from wire.
// Implement only the methods you need; others may be no-ops.
type DefaultDataSink interface {
	OnUnknown(ctx context.Context, msgType string, body []byte) error
}

// For binance data sink, implement BNDataSink instead.
type BNDataSink interface {
	OnBNTrade(ctx context.Context, t bWiredata.Trade) error
	OnBNKline(ctx context.Context, k bWiredata.Kline) error
	OnBNOrderBook(ctx context.Context, ob bWiredata.OrderBook) error
}
