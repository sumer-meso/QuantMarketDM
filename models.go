package model

import (
	"fmt"
)

// Config file model section.
type RabbitMQ struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type IndicatorSettings struct {
	Spread     string `json:"spread"`
	Imbalances string `json:"imbalances"`
}

type ConfigFile struct {
	RabbitmqTest RabbitMQ                                `json:"rabbitmqTest"`
	RabbitmqProd RabbitMQ                                `json:"rabbitmqProd"`
	Sources      map[string]map[string]IndicatorSettings `json:"sources"`
}

///////////////////////////////////

type Trade struct {
	Event         string
	Time          int64
	Symbol        string
	TradeID       int64
	Price         float64
	Quantity      float64
	BuyerOrderId  int64
	SellerOrderId int64
	TradeTime     int64
	IsBuyerMaker  bool
	Placeholder   bool
	LocalTime     int64
	Source        string
}

type OrderBookEntry struct {
	Price    float64
	Quantity float64
}

type MidPriceEntry struct {
	Amount   float64
	AskPrice float64
	BidPrice float64
	Value    float64
}

type SpreadEntry struct {
	MidPriceAmount float64
	SpreadAmount   float64
	AskPrice       float64
	BidPrice       float64
	Value          float64
}

type ImbalanceEntry struct {
	Layer   int
	BidSize float64
	AskSize float64
	Value   float64
}

type OrderBook struct {
	Symbol           string
	Event            string
	Time             int64
	LastUpdateID     int64
	FirstUpdateID    int64
	PreviousUpdateID int64
	Bids             []OrderBookEntry
	Asks             []OrderBookEntry
	LocalTime        int64
	Source           string
	MidPrice         MidPriceEntry
	Spread           SpreadEntry
	Imbalances       []ImbalanceEntry
}

func (o OrderBook) String() string {
	return fmt.Sprintf(
		"OrderBook (%s) (%s) Time(%d), LocalTime(%d):"+
			"%d -> %d, Pu: (%d)\n"+
			"Asks(%d):%+v ...\n"+
			"Bids(%d):%+v ...\n",
		o.Symbol, o.Event, o.Time, o.LocalTime,
		o.FirstUpdateID, o.LastUpdateID, o.PreviousUpdateID,
		len(o.Asks), o.Asks[:min(len(o.Asks), 50)],
		len(o.Bids), o.Bids[:min(len(o.Bids), 50)],
	)
}

type TrueRangeRatio struct {
	StartTime    int64
	EndTime      int64
	High         float64
	Low          float64
	OBMpAmount   float64
	OBMidPrice   float64
	OBUpdateTime int64
	Value        float64
}

// This is Kline model which is used to send out.
type Kline struct {
	Symbol          string
	TrueRangeRatios map[string]TrueRangeRatio
	Volumns         map[string]float64
	LocalTime       int64
	Source          string
}

func (k Kline) String() string {
	return fmt.Sprintf("Kline (%s), t:%v, v:%v, s:%s, l:%d",
		k.Symbol, k.TrueRangeRatios, k.Volumns, k.Source, k.LocalTime)
}

// WebSocket Event Structs
type WsTradeEvent struct {
	Time          int64  `json:"E"`
	Event         string `json:"e"`
	Symbol        string `json:"s"`
	TradeID       int64  `json:"t"`
	Price         string `json:"p"`
	Quantity      string `json:"q"`
	BuyerOrderId  int64  `json:"b"`
	SellerOrderId int64  `json:"a"`
	TradeTime     int64  `json:"T"`
	IsBuyerMaker  bool   `json:"m"`
	Placeholder   bool   `json:"M"` // add this field to avoid case insensitive unmarshaling
}

// WsKlineEvent define websocket kline event
type WsKlineEvent struct {
	Event  string  `json:"e"`
	Time   int64   `json:"E"`
	Symbol string  `json:"s"`
	Kline  WsKline `json:"k"`
}

// WsKline define websocket kline
type WsKline struct {
	StartTime            int64  `json:"t"`
	EndTime              int64  `json:"T"`
	Symbol               string `json:"s"`
	Interval             string `json:"i"`
	FirstTradeID         int64  `json:"f"`
	LastTradeID          int64  `json:"L"`
	Open                 string `json:"o"`
	Close                string `json:"c"`
	High                 string `json:"h"`
	Low                  string `json:"l"`
	Volume               string `json:"v"`
	TradeNum             int64  `json:"n"`
	IsFinal              bool   `json:"x"`
	QuoteVolume          string `json:"q"`
	ActiveBuyVolume      string `json:"V"`
	ActiveBuyQuoteVolume string `json:"Q"`
}

type WsDepthEvent struct {
	Event            string     `json:"e"`
	Time             int64      `json:"E"`
	Symbol           string     `json:"s"`
	FirstUpdateID    int64      `json:"U"`
	LastUpdateID     int64      `json:"u"`
	PreviousUpdateID int64      `json:"pu"`
	Bids             [][]string `json:"b"`
	Asks             [][]string `json:"a"`
}

type IndicatorSet struct {
	DataSrc string
	Symbol  string

	BidAmount float64
	AskAmount float64
	MidPrice  float64
	Spread    float64

	BidTop5     float64
	AskTop5     float64
	Imbalance5  float64
	Imbalance16 float64
	Imbalance32 float64

	TrueRangeRatio1s  float64
	TrueRangeRatio5s  float64
	TrueRangeRatio60s float64
	TrueRangeRatio5m  float64
	Volume1s          float64
	Volume5s          float64

	OrderBookUpdateTime int64
	KlineUpdateTime     int64
}
