package deribit

import (
	"encoding/json"
	"fmt"

	"github.com/sumer-meso/QuantMarketDM/proto"
)

type OrderBook struct {
	Instrument      string      `json:"instrument_name"`  // 合约名称，如 "BTC-PERPETUAL"
	Timestamp       int64       `json:"timestamp"`        // 订单簿时间戳（毫秒）
	ChangeID        int64       `json:"change_id"`        // 当前变更ID，用于跟踪订单簿更新
	PrevChangeID    int64       `json:"prev_change_id"`   // 上一个变更ID
	State           string      `json:"state"`            // 合约状态："open", "closed", "settled"
	Bids            [][]float64 `json:"bids"`             // 买单数组，每个元素为[价格, 数量]
	Asks            [][]float64 `json:"asks"`             // 卖单数组，每个元素为[价格, 数量]
	BestBidPrice    float64     `json:"best_bid_price"`   // 最优买价
	BestBidAmount   float64     `json:"best_bid_amount"`  // 最优买价对应的数量
	BestAskPrice    float64     `json:"best_ask_price"`   // 最优卖价
	BestAskAmount   float64     `json:"best_ask_amount"`  // 最优卖价对应的数量
	MarkPrice       float64     `json:"mark_price"`       // 标记价格，用于保证金计算和强平
	LastPrice       float64     `json:"last_price"`       // 最新成交价格
	IndexPrice      float64     `json:"index_price"`      // 指数价格，标的资产的参考价格
	OpenInterest    float64     `json:"open_interest"`    // 未平仓合约数量
	MaxPrice        float64     `json:"max_price"`        // 当日最高价格限制
	MinPrice        float64     `json:"min_price"`        // 当日最低价格限制
	SettlementPrice float64     `json:"settlement_price"` // 结算价格
	//EstimatedPrice  interface{} `json:"estimated_delivery_price"` // 预估交割价格
	Stats struct {
		Volume      float64 `json:"volume"`       // 24小时交易量
		PriceChange float64 `json:"price_change"` // 24小时价格变化量
		Low         float64 `json:"low"`          // 24小时最低价
		High        float64 `json:"high"`         // 24小时最高价
	} `json:"stats"` // 24小时统计数据
	DeliveryPrice   float64 `json:"delivery_price"`   // 交割价格
	MarkIV          float64 `json:"mark_iv"`          // 标记隐含波动率（期权合约）
	UnderlyingIndex string  `json:"underlying_index"` // 标的指数名称
	UnderlyingPrice float64 `json:"underlying_price"` // 标的资产价格
	InterestRate    float64 `json:"interest_rate"`    // 无风险利率（期权定价）
	Greeks          struct {
		Delta float64 `json:"delta"` // Delta值，价格敏感度
		Gamma float64 `json:"gamma"` // Gamma值，Delta的变化率
		Rho   float64 `json:"rho"`   // Rho值，利率敏感度
		Theta float64 `json:"theta"` // Theta值，时间衰减
		Vega  float64 `json:"vega"`  // Vega值，波动率敏感度
	} `json:"greeks"` // 希腊字母（期权风险指标）
	BidIV            float64 `json:"bid_iv"`             // 买价隐含波动率
	AskIV            float64 `json:"ask_iv"`             // 卖价隐含波动率
	LastPrice24h     float64 `json:"last_price_24h"`     // 24小时前的价格
	UsdIndexPrice    float64 `json:"usd_index_price"`    // USD指数价格
	BaseCurrency     string  `json:"base_currency"`      // 基础货币，如 "BTC"
	QuoteCurrency    string  `json:"quote_currency"`     // 计价货币，如 "USD"
	InstrumentType   string  `json:"instrument_type"`    // 合约类型："future", "option", "spot"
	MinTradeAmount   float64 `json:"min_trade_amount"`   // 最小交易数量
	SettlementPeriod string  `json:"settlement_period"`  // 结算周期
	LocalTime        string  `json:"localTime"`          // 本地时间戳
	Currency         *string `json:"currency,omitempty"` // 货币标识
}

func (ob *OrderBook) String() string {
	return "Deribit OrderBook for " + ob.Instrument
}

func (ob *OrderBook) RMQRoutingIdentifier() string {
	return fmt.Sprintf("deribit.orderbook.%s", proto.Ptr2Str(ob.Currency))
}

func (ob *OrderBook) RMQDataIdentifier() string {
	return "deribit.orderbook"
}

func (ob *OrderBook) RMQDataStoreTable() string {
	return ob.RMQRoutingIdentifier()
}

func (ob *OrderBook) RMQDataStoreIndex() string {
	return "instrument:1,timestamp:-1,true"
}

func (ob *OrderBook) RMQEncodeMessage() (*proto.MessageOverRabbitMQ, error) {
	if body, err := json.Marshal(ob); err != nil {
		return &proto.MessageOverRabbitMQ{}, err
	} else {
		return &proto.MessageOverRabbitMQ{
			RoutingKey:     ob.RMQRoutingIdentifier(),
			DataIdentifier: ob.RMQDataIdentifier(),
			StoreTable:     ob.RMQDataStoreTable(),
			StoreIndex:     ob.RMQDataStoreIndex(),
			Body:           body,
		}, nil
	}
}

func (ob *OrderBook) RMQDecodeMessage(m *proto.MessageOverRabbitMQ) error {
	if m.DataIdentifier != ob.RMQDataIdentifier() {
		return proto.NotMatchError{Expected: ob.RMQDataIdentifier(), Actual: m.DataIdentifier}
	}
	return json.Unmarshal(m.Body, ob)
}
