package tide

import (
	"context"

	"github.com/sumer-meso/QuantMarketDM/proto"
	"github.com/sumer-meso/QuantMarketDM/proto/binance"
	"github.com/sumer-meso/QuantMarketDM/proto/deribit/user"
)

// DataSink receives decoded data from wire.
// Implement only the methods you need; others may be no-ops.
type DefaultDataSink interface {
	OnUnknown(ctx context.Context, u *proto.Unknown) error
}

// For binance data sink, implement BNDataSink from when consuming these data types.
type BNDataSink interface {
	OnBNTrade(ctx context.Context, t *binance.Trade) error
	OnBNKline(ctx context.Context, k *binance.Kline) error
	OnBNOrderBook(ctx context.Context, ob *binance.OrderBook) error
	OnBNTradeLite(ctx context.Context, ob *binance.TradeLite) error
	OnBNAccountUpdate(ctx context.Context, ob *binance.AccountUpdate) error
	OnBNOrderUpdate(ctx context.Context, ob *binance.OrderUpdate) error
}

// For deribit data sink, implement DBDataSink from when consuming these data types.
type DBDataSink interface {
	OnDBPortfolio(ctx context.Context, t *user.Portfolio) error
	OnDBPosition(ctx context.Context, t *user.Position) error
	OnDBOrder(ctx context.Context, t *user.Order) error
}
