package user

import (
	"encoding/json"
	"fmt"

	"github.com/sumer-meso/QuantMarketDM/proto"
)

type ComboTradeLeg struct {
	InstrumentName string  `json:"instrument_name"`
	Amount         float64 `json:"amount"`
	Price          float64 `json:"price"`
	Direction      string  `json:"direction"`
}

type ClientInfo struct {
	ClientID     int64  `json:"client_id"`
	ClientLinkID int64  `json:"client_link_id"`
	Name         string `json:"name"`
}

type TradeAllocation struct {
	Amount     float64     `json:"amount"`
	ClientInfo *ClientInfo `json:"client_info,omitempty"`
	UserID     *int64      `json:"user_id,omitempty"`
}

type Trade struct {
	TradeID          string            `json:"trade_id"`
	TickDirection    int               `json:"tick_direction"`
	FeeCurrency      string            `json:"fee_currency"`
	API              *bool             `json:"api,omitempty"`
	Advanced         *string           `json:"advanced,omitempty"`
	OrderID          string            `json:"order_id"`
	Liquidity        string            `json:"liquidity"`
	PostOnly         *bool             `json:"post_only,omitempty"`
	Direction        string            `json:"direction"`
	Contracts        *float64          `json:"contracts,omitempty"`
	Mmp              *bool             `json:"mmp,omitempty"`
	Fee              float64           `json:"fee"`
	QuoteID          *string           `json:"quote_id,omitempty"`
	IndexPrice       float64           `json:"index_price"`
	Label            *string           `json:"label,omitempty"`
	BlockTradeID     *string           `json:"block_trade_id,omitempty"`
	Price            float64           `json:"price"`
	ComboID          *string           `json:"combo_id,omitempty"`
	MatchingID       *string           `json:"matching_id,omitempty"`
	OrderType        string            `json:"order_type"`
	TradeAllocations []TradeAllocation `json:"trade_allocations,omitempty"`
	ProfitLoss       float64           `json:"profit_loss"`
	Timestamp        int64             `json:"timestamp"`
	IV               *float64          `json:"iv,omitempty"`
	State            string            `json:"state"`
	UnderlyingPrice  *float64          `json:"underlying_price,omitempty"`
	BlockRfqQuoteID  *int64            `json:"block_rfq_quote_id,omitempty"`
	QuoteSetID       *string           `json:"quote_set_id,omitempty"`
	MarkPrice        float64           `json:"mark_price"`
	BlockRfqID       *int64            `json:"block_rfq_id,omitempty"`
	ComboTradeID     *int64            `json:"combo_trade_id,omitempty"`
	ReduceOnly       *bool             `json:"reduce_only,omitempty"`
	Amount           float64           `json:"amount"`
	Liquidation      *string           `json:"liquidation,omitempty"`
	TradeSeq         int64             `json:"trade_seq"`
	RiskReducing     *bool             `json:"risk_reducing,omitempty"`
	InstrumentName   string            `json:"instrument_name"`
	Legs             []ComboTradeLeg   `json:"legs,omitempty"`
	LocalTime        string            `json:"localTime"`          // 本地时间戳
	Account          *string           `json:"account,omitempty"`  // 账户标识
	Kind             *string           `json:"king,omitempty"`     // 标识
	Currency         *string           `json:"currency,omitempty"` // 货币标识
	Interval         *string           `json:"interval,omitempty"` // 时间间隔
}

func (t *Trade) String() string {
	return fmt.Sprintf("Deribit Trade (%s), dir:%s, amt:%.2f, prc:%.2f, fee:%.4f, pl:%.4f, ts:%d, lt:%s",
		t.InstrumentName, t.Direction, t.Amount, t.Price, t.Fee, t.ProfitLoss, t.Timestamp, t.LocalTime)
}

func (t *Trade) RMQRoutingIdentifier() string {
	return fmt.Sprintf("deribit.%v.changes.%v.%v.%v", t.Account, t.Kind, t.Currency, t.Interval)
}

func (t *Trade) RMQDataIdentifier() string {
	return "deribit.user.trade"
}

func (t *Trade) RMQDataStoreTable() string {
	return "deribit.user.trades"
}

func (t *Trade) RMQDataStoreIndex() string {
	return "trade_id:-1,false"
}

func (t *Trade) RMQEncodeMessage() (*proto.MessageOverRabbitMQ, error) {
	if body, err := json.Marshal(t); err != nil {
		return &proto.MessageOverRabbitMQ{}, err
	} else {
		return &proto.MessageOverRabbitMQ{
			RoutingKey:     t.RMQRoutingIdentifier(),
			DataIdentifier: t.RMQDataIdentifier(),
			StoreTable:     t.RMQDataStoreTable(),
			StoreIndex:     t.RMQDataStoreIndex(),
			Body:           body,
		}, nil
	}
}

func (t *Trade) RMQDecodeMessage(m *proto.MessageOverRabbitMQ) error {
	if m.DataIdentifier != t.RMQDataIdentifier() {
		return proto.NotMatchError{Expected: t.RMQDataIdentifier(), Actual: m.DataIdentifier}
	}
	return json.Unmarshal(m.Body, t)
}
