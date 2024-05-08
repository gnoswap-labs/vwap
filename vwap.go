package vwap

import (
	"fmt"
	"math/big"
)

// Token name
type TokenIdentifier string

const (
	wugnot TokenIdentifier = "gno.land/r/demo/wugnot" // base token (ratio = 1)
	foo    TokenIdentifier = "gno.land/r/demo/foo"
	bar    TokenIdentifier = "gno.land/r/demo/bar"
	baz    TokenIdentifier = "gno.land/r/demo/baz"
	qux    TokenIdentifier = "gno.land/r/demo/qux"
	gns    TokenIdentifier = "gno.land/r/demo/gns"
)

// TradeData represents the data for a single trade.
type TradeData struct {
	TokenName string
	Quantity  float64
	Ratio     float64
	Timestamp int
}

// lastPrices stores the last price of each token.
// This value will be used to show the last price if the token is not traded.
var lastPrices map[string]float64

// VWAP calculates the Volume Weighted Average Price (VWAP) for the given set of trades.
// It returns the last price if there are no trades.
func VWAP(trades []TradeData) float64 {
	var numerator, denominator float64

	for _, trade := range trades {
		numerator += trade.Quantity * trade.Ratio
		denominator += trade.Quantity
	}

	// return last price if there is no trade
	if denominator == 0 {
		return lastPrices[trades[0].TokenName]
	}

	vwap := numerator / denominator
	lastPrices[trades[0].TokenName] = vwap // save the last price

	return vwap
}

// updateTrades retrieves and updates the trade data from RPC API.
// It returns the updated trades and any error that occurred during the process.
func updateTrades(jsonStr string) ([]TradeData, error) {
	data, err := unmarshalResponseData(jsonStr)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response data: %v", err)
	}

	tradableTokens := []TokenIdentifier{wugnot, foo, bar, baz, qux, gns}

	var trades []TradeData
	for _, r := range data.Response {
		ratio, ok := new(big.Int).SetString(r.Ratio, 10)
		if !ok {
			return nil, fmt.Errorf("error converting ratio to big.Int")
		}

		dvisor := new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
		result := new(big.Float).Quo(new(big.Float).SetInt(ratio), new(big.Float).SetInt(dvisor))
		floatRatio, _ := result.Float64()

		// TODO; remove testing logic
		// testing purpose
		tokenName := TokenIdentifier(r.Token)
		if contains(tradableTokens, tokenName) {
			trade := TradeData{
				TokenName: string(tokenName),
				Quantity:  100,
				Ratio:     floatRatio,
				Timestamp: data.Stat.Timestamp,
			}
			trades = append(trades, trade)
		}
	}

	return trades, nil
}

func contains(tokens []TokenIdentifier, name TokenIdentifier) bool {
	for _, t := range tokens {
		if t == name {
			return true
		}
	}

	return false
}
