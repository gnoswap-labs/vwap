package vwap

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

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
	client := &http.Client{} // Timeout is managed by context

	// Create a context with a 30-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var apiResponse APIResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			log.Printf("Error decoding response: %v\n", err)
			return nil, err
		}
		return apiResponse.Data, nil
	}

	log.Printf("Received non-OK response: %d\n", resp.StatusCode)
	return nil, fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
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
		if token.TokenName == string(WUGNOT) {
			ratio, _ := new(big.Float).SetString(token.Ratio)
			baseRatio.Quo(ratio, big.NewFloat(math.Pow(2, 96)))
			break
		}
	}

	// calculate token prices based on the base token price and ratios.
	for _, token := range tokenData.TokenRatio {
		if token.TokenName != string(WUGNOT) {
			ratio, _ := new(big.Float).SetString(token.Ratio)
			tokenRatio := new(big.Float).Quo(ratio, big.NewFloat(math.Pow(2, 96)))
			tokenPrice := new(big.Float).Quo(baseRatio, tokenRatio)

			price, _ := tokenPrice.Float64()
			tokenPrices[token.TokenName] = price * baseTokenPrice
		}
	}

	return tokenPrices
}

func extractTrades(prices []TokenPrice) map[string][]TradeData {
	trades := make(map[string][]TradeData)
	for _, price := range prices {
		// calculatedVolume, err := strconv.ParseFloat(price.VolumeUSD24h, 64)
		// if err != nil {
		// 	fmt.Printf("Failed to parse volume for token %s: %v\n", price.Path, err)
		// 	continue
		// }
		usd, err := strconv.ParseFloat(price.USD, 64)
		if err != nil {
			fmt.Printf("failed to parse USD price for token %s: %v\n", price.Path, err)
			continue
		}

		trades[price.Path] = append(trades[price.Path], TradeData{
			TokenName: price.Path,
			// Volume: calculatedVolume,

			// hard coded because current calculated volume always be 0.
			// therefore, we can't calculate the VWAP correctly.
			// TODO: remove this hard coded value and use `calculatedVolume` instead.
			Volume:    100,
			Ratio:     usd,
			Timestamp: int(time.Now().Unix()),
		})
	}

	return trades
}
