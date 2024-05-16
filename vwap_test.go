package vwap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVWAPWithNoTrades(t *testing.T) {
	t.Parallel()

	trades := []TradeData{
		{TokenName: "Token1", Volume: 0, Ratio: 0, Timestamp: 1621000000},
	}

	lastPrices = map[string]float64{
		"Token1": 1.5,
	}

	expectedVWAP := 1.5
	actualVWAP, err := calculateVWAP(trades)

	assert.Nil(t, err, "Unexpected error")
	assert.Equal(t, expectedVWAP, actualVWAP, "VWAP calculation with no trades is incorrect")
}

func TestVWAPWith10MinuteInterval(t *testing.T) {
	// Mock trade data with different timestamps
	trades := []TradeData{
		{TokenName: "Token1", Volume: 100, Ratio: 1.5, Timestamp: 1621000000},
		{TokenName: "Token1", Volume: 200, Ratio: 1.8, Timestamp: 1621000180},
		{TokenName: "Token1", Volume: 150, Ratio: 1.6, Timestamp: 1621000420},
		{TokenName: "Token1", Volume: 300, Ratio: 1.7, Timestamp: 1621000600},
		{TokenName: "Token1", Volume: 250, Ratio: 1.9, Timestamp: 1621000900},
	}

	// Calculate VWAP for each 10-minute interval
	var expectedVWAPs []float64
	var actualVWAPs []float64

	intervalStart := trades[0].Timestamp
	var intervalTrades []TradeData

	for _, trade := range trades {
		if trade.Timestamp < intervalStart+600 {
			intervalTrades = append(intervalTrades, trade)
		} else {
			expectedVWAP := calculateExpectedVWAP(intervalTrades)
			expectedVWAPs = append(expectedVWAPs, expectedVWAP)

			actualVWAP, err := calculateVWAP(intervalTrades)
			assert.Nil(t, err, "Unexpected error")

			actualVWAPs = append(actualVWAPs, actualVWAP)

			intervalStart = trade.Timestamp
			intervalTrades = []TradeData{trade}
		}
	}

	// Process the last interval
	expectedVWAP := calculateExpectedVWAP(intervalTrades)
	expectedVWAPs = append(expectedVWAPs, expectedVWAP)

	actualVWAP, err := calculateVWAP(intervalTrades)
	assert.Nil(t, err, "Unexpected error")

	actualVWAPs = append(actualVWAPs, actualVWAP)

	// Compare expected and actual VWAPs for each interval
	assert.Equal(t, len(expectedVWAPs), len(actualVWAPs), "Unexpected number of intervals")

	for i := 0; i < len(expectedVWAPs); i++ {
		assert.InDelta(t, expectedVWAPs[i], actualVWAPs[i], 1e-9, "Incorrect VWAP for interval %d", i+1)
	}
}

func calculateExpectedVWAP(trades []TradeData) float64 {
	var numerator, denominator float64

	for _, trade := range trades {
		numerator += trade.Volume * trade.Ratio
		denominator += trade.Volume
	}

	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}
