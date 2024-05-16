package vwap

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// temporary filtering func. should move to a router later.

type SwapToken struct {
	Symbol string `json:"symbol"`
}

type Swap struct {
	Time         string    `json:"time"`
	TokenA       SwapToken `json:"tokenA"`
	TokenAAmount string    `json:"tokenAAmount"`
	TokenB       SwapToken `json:"tokenB"`
	TokenBAmount string    `json:"tokenBAmount"`
	TotalUsd     string    `json:"totalUsd"`
}

type ActivitySwapResponse struct {
	Data []Swap `json:"data"`
}

const (
	ActivitySwapEndpoint = "http://dev.api.gnoswap.io/v1/activity?type=%s"
	QueryTypeSwap        = "SWAP"
)

func FetchActivitySwap(endpoint, queryType string) ([]Swap, error) {
	client := &http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rpc := fmt.Sprintf(endpoint, queryType)
	req, err := http.NewRequestWithContext(ctx, "GET", rpc, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %v", resp.Status)
	}

	var apiResponse ActivitySwapResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	return apiResponse.Data, nil
}

var filterToken = map[string]bool{
	"USDC": true,
	"GNOT": true,
	"GNS":  true,
}

func FilterSwaps(swaps []Swap) []Swap {
	filteredSwaps := make([]Swap, 0)
	for _, swap := range swaps {
		_, okA := filterToken[swap.TokenA.Symbol]
		_, okB := filterToken[swap.TokenB.Symbol]
		if okA || okB {
			filteredSwaps = append(filteredSwaps, swap)
		}
	}
	return filteredSwaps
}
