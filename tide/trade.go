package tide

import "github.com/sumer-meso/QuantMarketDM/common"

type MarketTrade struct {
	Exchange string // "binance" / "deribit" / ...
	Symbol   string // Unified format，eg. "BTCUSDT" or "BTC-PERP"
	TradeID  int64
	TsUnixMs int64
	Price    common.Price
	Quantity common.Amount
	Side     common.Side
	IsMaker  bool
}
