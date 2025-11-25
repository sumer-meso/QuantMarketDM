package user

import (
	"encoding/json"
	"fmt"

	"github.com/sumer-meso/QuantMarketDM/proto"
)

type Position struct {
	AveragePrice          float64  `json:"average_price"`
	AveragePriceUSD       *float64 `json:"average_price_usd,omitempty"`
	Delta                 float64  `json:"delta"`
	Direction             string   `json:"direction"`
	FloatingProfitLoss    float64  `json:"floating_profit_loss"`
	FloatingProfitLossUSD *float64 `json:"floating_profit_loss_usd,omitempty"`
	Gamma                 *float64 `json:"gamma,omitempty"`
	IndexPrice            float64  `json:"index_price"`
	InitialMargin         float64  `json:"initial_margin"`
	InstrumentName        string   `json:"instrument_name"`
	InterestValue         *float64 `json:"interest_value,omitempty"`
	Kind                  string   `json:"kind"`
	Leverage              *int     `json:"leverage,omitempty"`
	MaintenanceMargin     float64  `json:"maintenance_margin"`
	MarkPrice             float64  `json:"mark_price"`
	RealizedFunding       *float64 `json:"realized_funding,omitempty"`
	RealizedProfitLoss    float64  `json:"realized_profit_loss"`
	SettlementPrice       *float64 `json:"settlement_price,omitempty"`
	Size                  float64  `json:"size"`
	SizeCurrency          *float64 `json:"size_currency,omitempty"`
	Theta                 *float64 `json:"theta,omitempty"`
	TotalProfitLoss       float64  `json:"total_profit_loss"`
	Vega                  *float64 `json:"vega,omitempty"`
	Parastamp             int64    `json:"parastamp"`          // 到达时间戳/index
	LocalTime             string   `json:"localTime"`          // 本地时间戳
	Account               *string  `json:"account,omitempty"`  // 账户标识
	Currency              *string  `json:"currency,omitempty"` // 货币标识
}

func (p *Position) String() string {
	return fmt.Sprintf("Deribit Position (%s), dir:%s, size:%.2f, avgPrice:%.2f, mPrice:%.2f, dlt:%.2f, upl:%.2f, rpl:%.2f, tpl:%.2f, lt:%s",
		p.InstrumentName, p.Direction, p.Size, p.AveragePrice, p.MarkPrice,
		p.Delta, p.FloatingProfitLoss, p.RealizedProfitLoss, p.TotalProfitLoss, p.LocalTime)
}

func (p *Position) RMQRoutingIdentifier() string {
	return fmt.Sprintf("deribit.%s.positions.%s", proto.PtrStr(p.Account), proto.PtrStr(p.Currency))
}

func (p *Position) RMQDataIdentifier() string {
	return "deribit.user.position"
}

func (p *Position) RMQDataStoreTable() string {
	return "deribit.user.positions"
}

func (p *Position) RMQDataStoreIndex() string {
	return "instrument_name:1,localTime:-1,true"
}

func (p *Position) RMQEncodeMessage() (*proto.MessageOverRabbitMQ, error) {
	if body, err := json.Marshal(p); err != nil {
		return &proto.MessageOverRabbitMQ{}, err
	} else {
		return &proto.MessageOverRabbitMQ{
			RoutingKey:     p.RMQRoutingIdentifier(),
			DataIdentifier: p.RMQDataIdentifier(),
			StoreTable:     p.RMQDataStoreTable(),
			StoreIndex:     p.RMQDataStoreIndex(),
			Body:           body,
		}, nil
	}
}

func (p *Position) RMQDecodeMessage(m *proto.MessageOverRabbitMQ) error {
	if m.DataIdentifier != p.RMQDataIdentifier() {
		return proto.NotMatchError{Expected: p.RMQDataIdentifier(), Actual: m.DataIdentifier}
	}
	return json.Unmarshal(m.Body, p)
}
