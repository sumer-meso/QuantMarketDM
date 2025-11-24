package user

import (
	"encoding/json"
	"fmt"

	"github.com/sumer-meso/QuantMarketDM/proto"
)

type Order struct {
	Quote                 *bool    `json:"quote,omitempty"`
	Triggered             *bool    `json:"triggered,omitempty"`
	Mobile                *bool    `json:"mobile,omitempty"`
	AppName               *string  `json:"app_name,omitempty"`
	Implv                 *float64 `json:"implv,omitempty"`
	RefreshAmount         *float64 `json:"refresh_amount,omitempty"`
	Usd                   *float64 `json:"usd,omitempty"`
	OtoOrderIds           []string `json:"oto_order_ids,omitempty"`
	API                   *bool    `json:"api,omitempty"`
	AveragePrice          *float64 `json:"average_price,omitempty"`
	Advanced              *string  `json:"advanced,omitempty"`
	OrderID               string   `json:"order_id"`
	PostOnly              *bool    `json:"post_only,omitempty"`
	FilledAmount          float64  `json:"filled_amount"`
	Trigger               *string  `json:"trigger,omitempty"`
	TriggerOrderID        *string  `json:"trigger_order_id,omitempty"`
	Direction             string   `json:"direction"`
	Contracts             *float64 `json:"contracts,omitempty"`
	IsSecondaryOto        *bool    `json:"is_secondary_oto,omitempty"`
	Replaced              *bool    `json:"replaced,omitempty"`
	MmpGroup              *string  `json:"mmp_group,omitempty"`
	Mmp                   *bool    `json:"mmp,omitempty"`
	LastUpdateTimestamp   int64    `json:"last_update_timestamp"`
	CreationTimestamp     int64    `json:"creation_timestamp"`
	CancelReason          *string  `json:"cancel_reason,omitempty"`
	MmpCancelled          *bool    `json:"mmp_cancelled,omitempty"`
	QuoteID               *string  `json:"quote_id,omitempty"`
	OrderState            string   `json:"order_state"`
	IsRebalance           *bool    `json:"is_rebalance,omitempty"`
	RejectPostOnly        *bool    `json:"reject_post_only,omitempty"`
	Label                 *string  `json:"label,omitempty"`
	IsLiquidation         *bool    `json:"is_liquidation,omitempty"`
	Price                 *float64 `json:"price,omitempty"` // 注意：可能是数字或字符串"market_price"
	Web                   *bool    `json:"web,omitempty"`
	TimeInForce           string   `json:"time_in_force"`
	TriggerReferencePrice *float64 `json:"trigger_reference_price,omitempty"`
	DisplayAmount         *float64 `json:"display_amount,omitempty"`
	OrderType             string   `json:"order_type"`
	IsPrimaryOtoco        *bool    `json:"is_primary_otoco,omitempty"`
	OriginalOrderType     *string  `json:"original_order_type,omitempty"`
	BlockTrade            *bool    `json:"block_trade,omitempty"`
	TriggerPrice          *float64 `json:"trigger_price,omitempty"`
	OcoRef                *string  `json:"oco_ref,omitempty"`
	TriggerOffset         *float64 `json:"trigger_offset,omitempty"`
	QuoteSetID            *string  `json:"quote_set_id,omitempty"`
	AutoReplaced          *bool    `json:"auto_replaced,omitempty"`
	ReduceOnly            *bool    `json:"reduce_only,omitempty"`
	Amount                float64  `json:"amount"`
	RiskReducing          *bool    `json:"risk_reducing,omitempty"`
	InstrumentName        string   `json:"instrument_name"`
	TriggerFillCondition  *string  `json:"trigger_fill_condition,omitempty"`
	PrimaryOrderID        *string  `json:"primary_order_id,omitempty"`
	LocalTime             string   `json:"localTime"`          // 本地时间戳
	Account               *string  `json:"account,omitempty"`  // 账户标识
	Kind                  *string  `json:"king,omitempty"`     // 标识
	Currency              *string  `json:"currency,omitempty"` // 货币标识
	Interval              *string  `json:"interval,omitempty"` // 时间间隔
}

func (o *Order) String() string {
	return fmt.Sprintf("Deribit Order (%s), dir:%s, amt:%.2f, fllAmt:%.2f, prc:%v, state:%s, lt:%s",
		o.InstrumentName, o.Direction, o.Amount, o.FilledAmount, o.Price, o.OrderState, o.LocalTime)
}

func (o *Order) RMQRoutingIdentifier() string {
	return fmt.Sprintf("deribit.%v.orders.%v.%v.%v", o.Account, o.Kind, o.Currency, o.Interval)
}

func (o *Order) RMQDataIdentifier() string {
	return "deribit.user.order"
}

func (o *Order) RMQDataStoreTable() string {
	return "deribit.user.orders"
}

func (o *Order) RMQDataStoreIndex() string {
	return "order_id:1,localTime:-1,true"
}

func (o *Order) RMQEncodeMessage() (*proto.MessageOverRabbitMQ, error) {
	if body, err := json.Marshal(o); err != nil {
		return &proto.MessageOverRabbitMQ{}, err
	} else {
		return &proto.MessageOverRabbitMQ{
			RoutingKey:     o.RMQRoutingIdentifier(),
			DataIdentifier: o.RMQDataIdentifier(),
			StoreTable:     o.RMQDataStoreTable(),
			StoreIndex:     o.RMQDataStoreIndex(),
			Body:           body,
		}, nil
	}
}

func (o *Order) RMQDecodeMessage(m *proto.MessageOverRabbitMQ) error {
	if m.DataIdentifier != o.RMQDataIdentifier() {
		return proto.NotMatchError{Expected: o.RMQDataIdentifier(), Actual: m.DataIdentifier}
	}
	return json.Unmarshal(m.Body, o)
}
