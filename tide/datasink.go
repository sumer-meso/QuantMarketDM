package tide

import (
	"context"

	"github.com/sumer-meso/QuantMarketDM/proto/binance"
)

// DataSink receives decoded data from wire.
// Implement only the methods you need; others may be no-ops.
type DefaultDataSink interface {
	OnUnknown(ctx context.Context, msgType string, body []byte) error
}

// For binance data sink, implement BNDataSink instead.
type BNDataSink interface {
	OnBNTrade(ctx context.Context, t binance.Trade) error
	OnBNKline(ctx context.Context, k binance.Kline) error
	OnBNOrderBook(ctx context.Context, ob binance.OrderBook) error
}
