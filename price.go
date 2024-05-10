package vwap

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"strconv"
)

const priceEndpoint = "http://dev.api.gnoswap.io/v1/tokens/prices"

type TokenPrice struct {
	Path              string       `json:"path"`
	USD               string       `json:"usd"`
	PricesBefore      PricesBefore `json:"pricesBefore"`
	MarketCap         string       `json:"marketCap"`
	LockedTokensUSD   string       `json:"lockedTokensUsd"`
	VolumeUSD24h      string       `json:"volumeUsd24h"`
	FeeUSD24h         string       `json:"feeUsd24h"`
	MostLiquidityPool string       `json:"mostLiquidityPool"`
	Last7d            []Last7d     `json:"last7d"`
}

type PricesBefore struct {
	LatestPrice string `json:"latestPrice"`
	Price1h     string `json:"price1h"`
	PriceToday  string `json:"priceToday"`
	Price1d     string `json:"price1d"`
	Price7d     string `json:"price7d"`
	Price30d    string `json:"price30d"`
	Price60d    string `json:"price60d"`
	Price90d    string `json:"price90d"`
}

type Last7d struct {
	Date  string `json:"date"`
	Price string `json:"price"`
}

type APIResponse struct {
	Error json.RawMessage `json:"error"`
	Data  []TokenPrice    `json:"data"`
}

func fetchTokenPrices(endpoint string) ([]TokenPrice, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResponse APIResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// calculateVolume calculates the total volume in the USD for each token.
func calculateVolume(prices []TokenPrice) map[string]float64 {
	volumeByToken := make(map[string]float64)

	for _, price := range prices {
		volume, err := strconv.ParseFloat(price.VolumeUSD24h, 64)
		if err != nil {
			fmt.Printf("failed to parse volume for token %s: %v\n", price.Path, err)
		}

		volumeByToken[price.Path] = volume
	}

	return volumeByToken
}

// calculateTokenUSDPrices calculates the USD price based on a base token (wugnot) price and token ratios.
func calculateTokenUSDPrices(tokenData *TokenData, baseTokenPrice float64) map[string]float64 {
	tokenPrices := make(map[string]float64)
	baseRatio := new(big.Float)

	// fund the base token ratio
	for _, token := range tokenData.TokenRatio {
		if token.TokenName == string(wugnot) {
			ratio, _ := new(big.Float).SetString(token.Ratio)
			baseRatio.Quo(ratio, big.NewFloat(math.Pow(2, 96)))
			break
		}
	}

	// calculate token prices based on the base token price and ratios.
	for _, token := range tokenData.TokenRatio {
		if token.TokenName != string(wugnot) {
			ratio, _ := new(big.Float).SetString(token.Ratio)
			tokenRatio := new(big.Float).Quo(ratio, big.NewFloat(math.Pow(2, 96)))
			tokenPrice := new(big.Float).Quo(baseRatio, tokenRatio)

			price, _ := tokenPrice.Float64()
			tokenPrices[token.TokenName] = price * baseTokenPrice
		}
	}

	return tokenPrices
}
