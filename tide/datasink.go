package tide

import (
	"context"

	"github.com/sumer-meso/QuantMarketDM/proto"
	"github.com/sumer-meso/QuantMarketDM/proto/binance"
)

// DataSink receives decoded data from wire.
// Implement only the methods you need; others may be no-ops.
type DefaultDataSink interface {
	OnUnknown(ctx context.Context, u *proto.Unknown) error
}

// For binance data sink, implement BNDataSink instead.
type BNDataSink interface {
	OnBNTrade(ctx context.Context, t *binance.Trade) error
	OnBNKline(ctx context.Context, k *binance.Kline) error
	OnBNOrderBook(ctx context.Context, ob *binance.OrderBook) error
	OnBNTradeLite(ctx context.Context, ob *binance.TradeLite) error
	OnBNAccountUpdate(ctx context.Context, ob *binance.AccountUpdate) error
	OnBNOrderUpdate(ctx context.Context, ob *binance.OrderUpdate) error
}
