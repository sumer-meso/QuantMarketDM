package user

import (
	"encoding/json"
	"fmt"

	"github.com/sumer-meso/QuantMarketDM/proto"
)

type Portfolio struct {
	OptionsPL                    float64            `json:"options_pl"`
	ProjectedDeltaTotal          float64            `json:"projected_delta_total"`
	OptionsThetaMap              map[string]float64 `json:"options_theta_map"`
	TotalMarginBalanceUSD        *float64           `json:"total_margin_balance_usd,omitempty"`
	TotalDeltaTotalUSD           *float64           `json:"total_delta_total_usd,omitempty"`
	AvailableWithdrawalFunds     float64            `json:"available_withdrawal_funds"`
	EstimatedLiquidationRatioMap map[string]float64 `json:"estimated_liquidation_ratio_map"`
	OptionsSessionRPL            float64            `json:"options_session_rpl"`
	FuturesSessionRPL            float64            `json:"futures_session_rpl"`
	TotalPL                      float64            `json:"total_pl"`
	AdditionalReserve            float64            `json:"additional_reserve"`
	OptionsSessionUPL            float64            `json:"options_session_upl"`
	CrossCollateralEnabled       bool               `json:"cross_collateral_enabled"`
	DeltaTotalMap                map[string]float64 `json:"delta_total_map"`
	OptionsValue                 float64            `json:"options_value"`
	OptionsVegaMap               map[string]float64 `json:"options_vega_map"`
	MaintenanceMargin            float64            `json:"maintenance_margin"`
	FuturesSessionUPL            float64            `json:"futures_session_upl"`
	PortfolioMarginingEnabled    bool               `json:"portfolio_margining_enabled"`
	FuturesPL                    float64            `json:"futures_pl"`
	OptionsGammaMap              map[string]float64 `json:"options_gamma_map"`
	Currency                     string             `json:"currency"`
	OptionsDelta                 float64            `json:"options_delta"`
	InitialMargin                float64            `json:"initial_margin"`
	ProjectedMaintenanceMargin   float64            `json:"projected_maintenance_margin"`
	AvailableFunds               float64            `json:"available_funds"`
	Equity                       float64            `json:"equity"`
	MarginModel                  string             `json:"margin_model"`
	Balance                      float64            `json:"balance"`
	SessionUPL                   float64            `json:"session_upl"`
	MarginBalance                float64            `json:"margin_balance"`
	OptionsTheta                 float64            `json:"options_theta"`
	TotalInitialMarginUSD        *float64           `json:"total_initial_margin_usd,omitempty"`
	EstimatedLiquidationRatio    *float64           `json:"estimated_liquidation_ratio,omitempty"` // DEPRECATED
	SessionRPL                   float64            `json:"session_rpl"`
	FeeBalance                   float64            `json:"fee_balance"`
	TotalMaintenanceMarginUSD    *float64           `json:"total_maintenance_margin_usd,omitempty"`
	OptionsVega                  float64            `json:"options_vega"`
	ProjectedInitialMargin       float64            `json:"projected_initial_margin"`
	OptionsGamma                 float64            `json:"options_gamma"`
	TotalEquityUSD               *float64           `json:"total_equity_usd,omitempty"`
	DeltaTotal                   float64            `json:"delta_total"`
	LocalTime                    string             `json:"localTime"`         // 本地时间戳
	Account                      *string            `json:"account,omitempty"` // 账户标识
}

func (p *Portfolio) String() string {
	return fmt.Sprintf("Deribit Portfolio (%s), eq:%.2f, bal:%.2f, marg:%.2f, upl:%.2f, rpl:%.2f, tpl:%.2f, lt:%s",
		p.Currency, p.Equity, p.Balance, p.MarginBalance, p.SessionUPL, p.SessionRPL, p.TotalPL, p.LocalTime)
}

func (p *Portfolio) RMQRoutingIdentifier() string {
	return fmt.Sprintf("deribit.%v.portfolio.%s", p.Account, p.Currency)
}

func (p *Portfolio) RMQDataIdentifier() string {
	return "deribit.user.portfolio"
}

func (p *Portfolio) RMQDataStoreTable() string {
	return "deribit.user.portfolio"
}

func (p *Portfolio) RMQDataStoreIndex() string {
	return "localTime:-1"
}

func (p *Portfolio) RMQEncodeMessage() (*proto.MessageOverRabbitMQ, error) {
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

func (p *Portfolio) RMQDecodeMessage(m *proto.MessageOverRabbitMQ) error {
	if m.DataIdentifier != p.RMQDataIdentifier() {
		return proto.NotMatchError{Expected: p.RMQDataIdentifier(), Actual: m.DataIdentifier}
	}
	return json.Unmarshal(m.Body, p)
}
