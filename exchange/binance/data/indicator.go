package data

type IndicatorSet struct {
	DataSrc Source
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
