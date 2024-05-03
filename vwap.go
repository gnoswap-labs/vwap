package vwap

type Trade struct {
	Price  float64
	Volume int
}

func VWAP(trades []Trade) float64 {
	var (
		totalVolume int
		totalValue  float64
	)

	for _, trade := range trades {
		totalVolume += trade.Volume
		totalValue += trade.Price * float64(trade.Volume)
	}

	if totalVolume == 0 {
		return 0
	}

	return totalValue / float64(totalVolume)
}
