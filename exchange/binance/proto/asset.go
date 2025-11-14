package proto

import (
	"encoding/json"
	"fmt"
	"time"
)

type WsUsdtEventBase struct {
	WSEventBase
	TransTime     int64            `json:"T"`
	AccountUpdate *WsAccountUpdate `json:"a,omitempty"`
	OrderUpdate   *WsOrderUpdate   `json:"o,omitempty"`
}

func (wueb *WsUsdtEventBase) ToAccountUpdate(source string) AccountUpdate {
	if wueb.AccountUpdate == nil {
		return AccountUpdate{}
	}
	wau := wueb.AccountUpdate
	au := AccountUpdate{
		EventBase: EventBase(wueb.WSEventBase),
		TransTime: wueb.TransTime,
		Reason:    wau.Reason,
		Balances:  make([]AccountBalance, 0, len(wau.Balances)),
		Positions: make([]AccountPosition, 0, len(wau.Positions)),
		LocalBase: LocalBase{LocalTime: time.Now().UnixMilli(), Source: source},
	}
	for _, b := range wau.Balances {
		au.Balances = append(au.Balances, AccountBalance(b))
	}
	for _, p := range wau.Positions {
		au.Positions = append(au.Positions, AccountPosition(p))
	}

	return au
}

func (wueb *WsUsdtEventBase) ToOrderUpdate(source string) OrderUpdate {
	if wueb.OrderUpdate == nil {
		return OrderUpdate{}
	}
	wou := wueb.OrderUpdate
	return OrderUpdate{
		EventBase:       EventBase(wueb.WSEventBase),
		TransTime:       wueb.TransTime,
		Symbol:          wou.Symbol,
		ClientOrderID:   wou.ClientOrderID,
		Side:            wou.Side,
		Type:            wou.Type,
		TimeInForce:     wou.TimeInForce,
		OrigQty:         wou.OrigQty,
		Price:           wou.Price,
		AvgPrice:        wou.AvgPrice,
		StopPrice:       wou.StopPrice,
		ExecType:        wou.ExecType,
		OrderStatus:     wou.OrderStatus,
		OrderID:         wou.OrderID,
		LastFilledQty:   wou.LastFilledQty,
		CumFilledQty:    wou.CumFilledQty,
		LastFilledPrice: wou.LastFilledPrice,
		FeeAsset:        wou.FeeAsset,
		FeeAmount:       wou.FeeAmount,
		TradeTime:       wou.TradeTime,
		TradeID:         wou.TradeID,
		Maker:           wou.Maker,
		ReduceOnly:      wou.ReduceOnly,
		WorkingType:     wou.WorkingType,
		OrigType:        wou.OrigType,
		PositionSide:    wou.PositionSide,
		ClosePosition:   wou.ClosePosition,
		RealizedPnL:     wou.RealizedPnL,
		PriceProtect:    wou.PriceProtect,
		LocalBase:       LocalBase{LocalTime: time.Now().UnixMilli(), Source: source},
	}
}

// ========== ACCOUNT_UPDATE ==========

type AccountUpdate struct {
	EventBase
	TransTime int64
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

func (au *AccountUpdate) RMQRoutingIdentifier() string {
	return fmt.Sprintf(
		"binance.%s.accountupdate",
		au.LocalBase.Source,
	)
}

func (au *AccountUpdate) RMQDataIdentifier() string {
	return "binance.accountupdate"
}

func (au *AccountUpdate) RMQDataStoreTable() string {
	return au.RMQRoutingIdentifier()
}

func (au *AccountUpdate) RMQDataStoreIndex() string {
	return "LocalTime:-1"
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

func (au *AccountUpdate) RMQDecodeMessage(m MessageOverRabbitMQ) error {
	if m.DataIdentifier != au.RMQDataIdentifier() {
		return NotMatchError{Expected: au.RMQDataIdentifier(), Actual: m.DataIdentifier}
	}
	return json.Unmarshal(m.Body, au)
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
	TransTime       int64
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

func (ou *OrderUpdate) RMQRoutingIdentifier() string {
	return fmt.Sprintf(
		"binance.%s.orderupdate",
		ou.LocalBase.Source,
	)
}

func (ou *OrderUpdate) RMQDataIdentifier() string {
	return "binance.orderupdate"
}

func (ou *OrderUpdate) RMQDataStoreTable() string {
	return ou.RMQRoutingIdentifier()
}

func (ou *OrderUpdate) RMQDataStoreIndex() string {
	return "LocalTime:-1"
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

func (ou *OrderUpdate) RMQDecodeMessage(m MessageOverRabbitMQ) error {
	if m.DataIdentifier != ou.RMQDataIdentifier() {
		return NotMatchError{Expected: ou.RMQDataIdentifier(), Actual: m.DataIdentifier}
	}
	return json.Unmarshal(m.Body, ou)
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
	TradeID  int64
	Symbol   string
	Side     string
	Price    string
	Quantity string
	LocalBase
}

func (tl *TradeLite) RMQRoutingIdentifier() string {
	return fmt.Sprintf(
		"binance.%s.tradelite",
		tl.LocalBase.Source,
	)
}

func (tl *TradeLite) RMQDataIdentifier() string {
	return "binance.tradelite"
}

func (tl *TradeLite) RMQDataStoreTable() string {
	return tl.RMQRoutingIdentifier()
}

func (tl *TradeLite) RMQDataStoreIndex() string {
	return "LocalTime:-1;TradeID:-1"
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

func (tl *TradeLite) RMQDecodeMessage(m MessageOverRabbitMQ) error {
	if m.DataIdentifier != tl.RMQDataIdentifier() {
		return NotMatchError{Expected: tl.RMQDataIdentifier(), Actual: m.DataIdentifier}
	}
	return json.Unmarshal(m.Body, tl)
}

type WsTradeLite struct {
	WSEventBase
	TradeID  int64  `json:"T"`
	Symbol   string `json:"s"`
	Side     string `json:"S"`
	Price    string `json:"p"`
	Quantity string `json:"q"`
}

func (wtl *WsTradeLite) ToTradeLite(source string) TradeLite {
	return TradeLite{
		EventBase: EventBase(wtl.WSEventBase),
		TradeID:   wtl.TradeID,
		Symbol:    wtl.Symbol,
		Side:      wtl.Side,
		Price:     wtl.Price,
		Quantity:  wtl.Quantity,
		LocalBase: LocalBase{LocalTime: time.Now().UnixMilli(), Source: source},
	}
}
