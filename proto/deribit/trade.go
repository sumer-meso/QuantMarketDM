package deribit

import (
	"encoding/json"
	"fmt"

	"github.com/sumer-meso/QuantMarketDM/proto"
)

type Trade struct {
	Amount             float64  `json:"amount"`
	BlockRfqID         *int64   `json:"block_rfq_id,omitempty"`
	BlockTradeID       *string  `json:"block_trade_id,omitempty"`
	BlockTradeLegCount *int     `json:"block_trade_leg_count,omitempty"`
	ComboID            *string  `json:"combo_id,omitempty"`
	ComboTradeID       *string  `json:"combo_trade_id,omitempty"`
	Contracts          *float64 `json:"contracts,omitempty"`
	Direction          string   `json:"direction"`
	IndexPrice         float64  `json:"index_price"`
	InstrumentName     string   `json:"instrument_name"`
	IV                 *float64 `json:"iv,omitempty"`
	Liquidation        *string  `json:"liquidation,omitempty"`
	MarkPrice          float64  `json:"mark_price"`
	Price              float64  `json:"price"`
	TickDirection      int      `json:"tick_direction"`
	Timestamp          int64    `json:"timestamp"`
	TradeID            string   `json:"trade_id"`
	TradeSeq           int64    `json:"trade_seq"`
	LocalTime          string   `json:"localTime"`          // 本地时间戳
	Kind               *string  `json:"kind,omitempty"`     // 标识
	Currency           *string  `json:"currency,omitempty"` // 货币标识
	Interval           *string  `json:"interval,omitempty"` // 时间间隔
}

func (t *Trade) String() string {
	return fmt.Sprintf("Deribit Trade (%s), dir:%s, amt:%.2f, prc:%.2f, ts:%d, lt:%s",
		t.InstrumentName, t.Direction, t.Amount, t.Price, t.Timestamp, t.LocalTime)
}

func (t *Trade) RMQRoutingIdentifier() string {
	return fmt.Sprintf("deribit.trades.%v.%v.%v", t.Kind, t.Currency, t.Interval)
}

func (t *Trade) RMQDataIdentifier() string {
	return "deribit.trade"
}

func (t *Trade) RMQDataStoreTable() string {
	return fmt.Sprintf("deribit.trades.%v.%v", t.Kind, t.Currency)
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
