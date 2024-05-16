package vwap

import (
	"fmt"
	"strconv"
)

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
