package vwap

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
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
	client := &http.Client{}

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

func extractTrades(prices []TokenPrice, volumeByToken map[string]float64) map[string][]TradeData {
	trades := make(map[string][]TradeData)
	for _, price := range prices {
		usd, err := strconv.ParseFloat(price.USD, 64)
		if err != nil {
			fmt.Printf("failed to parse USD price for token %s: %v\n", price.Path, err)
			continue
		}

		volume, ok := volumeByToken[price.Path]
		if !ok {
			fmt.Printf("volume not found for token %s\n", price.Path)
			continue
		}

		trades[price.Path] = append(trades[price.Path], TradeData{
			TokenName: price.Path,
			Volume:    volume,
			Ratio:     usd,
			Timestamp: int(time.Now().Unix()),
		})
	}

	return trades
}
