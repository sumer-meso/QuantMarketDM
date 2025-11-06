package wiredata

import (
	"encoding/json"
	"fmt"
)

type WsUsdtEventBase struct {
	WSEventBase
	AccountUpdate *WsAccountUpdate `json:"a,omitempty"`
	OrderUpdate   *WsOrderUpdate   `json:"o,omitempty"`
}

// ========== ACCOUNT_UPDATE ==========

type AccountUpdate struct {
	EventBase
	Reason    string
	Balances  []AccountBalance
	Positions []AccountPosition
	LocalBase
}

type AccountBalance struct {
	Asset         string
	WalletBalance string
	CrossWallet   string
}

type AccountPosition struct {
	Symbol            string
	PositionAmount    string
	EntryPrice        string
	CumRealizedProfit string
	UnrealizedProfit  string
	MarginType        string
	IsolatedWallet    string
	PositionSide      string
}

const auIndexSpecInRMQ = "{LocalTime:-1}"

func (au *AccountUpdate) RMQRoutingIdentifier() string {
	return fmt.Sprintf(
		"binance.accountupdate.%s",
		auIndexSpecInRMQ,
	)
}

func (au *AccountUpdate) RMQDataIdentifier() string {
	return "binance.accountupdate"
}

func (au *AccountUpdate) RMQEncodeMessage() (MessageOverRabbitMQ, error) {
	if body, err := json.Marshal(au); err != nil {
		return MessageOverRabbitMQ{}, err
	} else {
		return MessageOverRabbitMQ{
			RoutingKey:     au.RMQRoutingIdentifier(),
			DataIdentifier: au.RMQDataIdentifier(),
			Body:           body,
		}, nil
	}
}

type WsAccountUpdate struct {
	Reason    string              `json:"m"`
	Balances  []WsAccountBalance  `json:"B"`
	Positions []WsAccountPosition `json:"P"`
}

type WsAccountBalance struct {
	Asset         string `json:"a"`
	WalletBalance string `json:"wb"`
	CrossWallet   string `json:"cw"`
}

type WsAccountPosition struct {
	Symbol            string `json:"s"`
	PositionAmount    string `json:"pa"`
	EntryPrice        string `json:"ep"`
	CumRealizedProfit string `json:"cr"`
	UnrealizedProfit  string `json:"up"`
	MarginType        string `json:"mt"`
	IsolatedWallet    string `json:"iw"`
	PositionSide      string `json:"ps"`
}

// ========== ORDER_TRADE_UPDATE ==========

type OrderUpdate struct {
	EventBase
	Symbol          string
	ClientOrderID   string
	Side            string
	Type            string
	TimeInForce     string
	OrigQty         string
	Price           string
	AvgPrice        string
	StopPrice       string
	ExecType        string
	OrderStatus     string
	OrderID         int64
	LastFilledQty   string
	CumFilledQty    string
	LastFilledPrice string
	FeeAsset        string
	FeeAmount       string
	TradeTime       int64
	TradeID         int64
	Maker           bool
	ReduceOnly      bool
	WorkingType     string
	OrigType        string
	PositionSide    string
	ClosePosition   bool
	RealizedPnL     string
	PriceProtect    bool
	LocalBase
}

const ouIndexSpecInRMQ = "{LocalTime:-1}"

func (ou *OrderUpdate) RMQRoutingIdentifier() string {
	return fmt.Sprintf(
		"binance.orderupdate.%s",
		ouIndexSpecInRMQ,
	)
}

func (au *OrderUpdate) RMQDataIdentifier() string {
	return "binance.orderupdate"
}

func (ou *OrderUpdate) RMQEncodeMessage() (MessageOverRabbitMQ, error) {
	if body, err := json.Marshal(ou); err != nil {
		return MessageOverRabbitMQ{}, err
	} else {
		return MessageOverRabbitMQ{
			RoutingKey:     ou.RMQRoutingIdentifier(),
			DataIdentifier: ou.RMQDataIdentifier(),
			Body:           body,
		}, nil
	}
}

type WsOrderUpdate struct {
	Symbol          string `json:"s"`
	ClientOrderID   string `json:"c"`
	Side            string `json:"S"`
	Type            string `json:"o"`
	TimeInForce     string `json:"f"`
	OrigQty         string `json:"q"`
	Price           string `json:"p"`
	AvgPrice        string `json:"ap"`
	StopPrice       string `json:"sp"`
	ExecType        string `json:"x"`
	OrderStatus     string `json:"X"`
	OrderID         int64  `json:"i"`
	LastFilledQty   string `json:"l"`
	CumFilledQty    string `json:"z"`
	LastFilledPrice string `json:"L"`
	FeeAsset        string `json:"N"`
	FeeAmount       string `json:"n"`
	TradeTime       int64  `json:"T"`
	TradeID         int64  `json:"t"`
	Maker           bool   `json:"m"`
	ReduceOnly      bool   `json:"R"`
	WorkingType     string `json:"wt"`
	OrigType        string `json:"ot"`
	PositionSide    string `json:"ps"`
	ClosePosition   bool   `json:"cp"`
	RealizedPnL     string `json:"rp"`
	PriceProtect    bool   `json:"pP"`
}

// ========== TRADE_LITE ==========

type TradeLite struct {
	EventBase
	Symbol   string
	Side     string
	Price    string
	Quantity string
	LocalBase
}

const tlIndexSpecInRMQ = "{LocalTime:-1,TransTime:-1}"

func (tl *TradeLite) RMQRoutingIdentifier() string {
	return fmt.Sprintf(
		"binance.tradelite.%s",
		tlIndexSpecInRMQ,
	)
}

func (tl *TradeLite) RMQDataIdentifier() string {
	return "binance.tradelite"
}

func (tl *TradeLite) RMQEncodeMessage() (MessageOverRabbitMQ, error) {
	if body, err := json.Marshal(tl); err != nil {
		return MessageOverRabbitMQ{}, err
	} else {
		return MessageOverRabbitMQ{
			RoutingKey:     tl.RMQRoutingIdentifier(),
			DataIdentifier: tl.RMQDataIdentifier(),
			Body:           body,
		}, nil
	}
}

type WsTradeLite struct {
	WSEventBase
	Symbol   string `json:"s"`
	Side     string `json:"S"`
	Price    string `json:"p"`
	Quantity string `json:"q"`
}
