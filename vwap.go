package vwap

import "fmt"

const priceEndpoint = "http://dev.api.gnoswap.io/v1/tokens/prices"

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

// VWAP calculates the Volume Weighted Average Price (VWAP) for the given set of trades.
// It returns the last price if there are no trades.
func VWAP(trades []TradeData) (float64, error) {
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
			return 0, fmt.Errorf("no last price available for token %s", trades[0].TokenName)
		}
		return lastPrice, nil
	}

	vwap := numerator / denominator
	lastPrices[trades[0].TokenName] = vwap // save the last price

	store(trades[0].TokenName, vwap, trades[0].Timestamp)

	return vwap, nil
}

func FetchAndCalculateVWAP() (map[string]float64, error) {
	prices, err := fetchTokenPrices(priceEndpoint)
	if err != nil {
		return nil, err
	}

	trades := extractTrades(prices)
	vwapResults := make(map[string]float64)

	for tokenName, tradeData := range trades {
		vwap, err := VWAP(tradeData)
		if err != nil {
			return nil, err
		}
		vwapResults[tokenName] = vwap
	}

	return vwapResults, nil
}
