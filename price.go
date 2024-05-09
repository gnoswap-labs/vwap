package vwap

import (
	"encoding/json"
	"net/http"
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
