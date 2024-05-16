package vwap

import (
	"fmt"
)

// Token name
type TokenIdentifier string

const (
	WUGNOT TokenIdentifier = "gno.land/r/demo/wugnot" // base token (ratio = 1)
	FOO    TokenIdentifier = "gno.land/r/demo/foo"
	BAR    TokenIdentifier = "gno.land/r/demo/bar"
	BAZ    TokenIdentifier = "gno.land/r/demo/baz"
	QUX    TokenIdentifier = "gno.land/r/demo/qux"
	GNS    TokenIdentifier = "gno.land/r/demo/gns"
)

// TradeData represents the data for a single trade.
type TradeData struct {
	TokenName string
	Volume    float64
	Ratio     float64
	Timestamp int
}

// lastPrices stores the last price of each token.
// This value will be used to show the last price if the token is not traded.
var lastPrices map[string]float64

func init() {
	lastPrices = make(map[string]float64)
}

func VWAP() (map[string]float64, error) {
	prices, err := fetchTokenPrices(priceEndpoint)
	if err != nil {
		return nil, err
	}

	volumeByToken := calculateVolume(prices)
	trades := extractTrades(prices, volumeByToken)
	vwapResults := make(map[string]float64)

	for tokenName, tradeData := range trades {
		vwap, err := calculateVWAP(tradeData)
		if err != nil {
			return nil, err
		}
		vwapResults[tokenName] = vwap
	}

	return vwapResults, nil
}

// calculateVWAP calculates the Volume Weighted Average Price (calculateVWAP) for the given set of trades.
// It returns the last price if there are no trades.
func calculateVWAP(trades []TradeData) (float64, error) {
	var numerator, denominator float64

	if len(trades) == 0 {
		return 0, fmt.Errorf("no trades found")
	}

	for _, trade := range trades {
		numerator += trade.Volume * trade.Ratio
		denominator += trade.Volume
	}

	// return last price if there is no trade
	if denominator == 0 {
		lastPrice, ok := lastPrices[trades[0].TokenName]
		if !ok {
			return 0, nil
		}
		return lastPrice, nil
	}

	vwap := numerator / denominator
	lastPrices[trades[0].TokenName] = vwap // save the last price

	store(trades[0].TokenName, vwap, trades[0].Timestamp)

	return vwap, nil
}
