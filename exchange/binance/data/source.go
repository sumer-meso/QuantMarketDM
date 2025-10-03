package data

type Source string

const (
	Spot Source = "spot"
	Usdt Source = "usdt"
	Coin Source = "coin"
)

var AllDataSrc = []Source{Spot, Usdt, Coin}
