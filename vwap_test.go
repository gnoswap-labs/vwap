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
	actualVWAP := VWAP(trades)

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

			actualVWAP := VWAP(intervalTrades)
			actualVWAPs = append(actualVWAPs, actualVWAP)

			intervalStart = trade.Timestamp
			intervalTrades = []TradeData{trade}
		}
	}

	// Process the last interval
	expectedVWAP := calculateExpectedVWAP(intervalTrades)
	expectedVWAPs = append(expectedVWAPs, expectedVWAP)

	actualVWAP := VWAP(intervalTrades)
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

func TestUpdateTrades(t *testing.T) {
	t.Parallel()

	// Mock the RPC response
	mockResponse := `{
		"stat": {
			"height": 1000,
			"timestamp": 1621000000
		},
		"response": [
			{
				"token": "gno.land/r/demo/wugnot",
				"ratio": "79228162514264337593543950336"
			},
			{
				"token": "gno.land/r/demo/foo",
				"ratio": "121074204438706762251182081654"
			},
			{
				"token": "gno.land/r/demo/qux",
				"ratio": "60536002769587966558221762891"
			},
			{
				"token": "gno.land/r/demo/gns",
				"ratio": "236174327852992866806677676716"
			},
			{
				"token": "gno.land/r/demo/bar",
				"ratio": "0"
			},
			{
				"token": "gno.land/r/demo/baz",
				"ratio": "0"
			}
		]
	}`

	trades, err := updateTrades(mockResponse)

	expectedTimestamp := 1621000000

	assert.NoError(t, err, "updateTrades returned an error")
	assert.Len(t, trades, 6, "Incorrect number of trades returned")

	expectedTrades := []TradeData{
		{TokenName: string(wugnot), Volume: 100, Ratio: 1, Timestamp: expectedTimestamp},
		{TokenName: string(foo), Volume: 100, Ratio: 1.5281713042, Timestamp: expectedTimestamp},
		{TokenName: string(qux), Volume: 100, Ratio: 0.7640717751, Timestamp: expectedTimestamp},
		{TokenName: string(gns), Volume: 100, Ratio: 2.980939105, Timestamp: expectedTimestamp},
		{TokenName: string(bar), Volume: 100, Ratio: 0, Timestamp: expectedTimestamp},
		{TokenName: string(baz), Volume: 100, Ratio: 0, Timestamp: expectedTimestamp},
	}

	for i, trade := range trades {
		assert.Equal(t, expectedTrades[i].TokenName, trade.TokenName, "Incorrect token name")
		assert.Equal(t, expectedTrades[i].Volume, trade.Volume, "Incorrect quantity")
		assert.InDelta(t, expectedTrades[i].Ratio, trade.Ratio, 1e-9, "Incorrect ratio")
		assert.Equal(t, expectedTrades[i].Timestamp, trade.Timestamp, "Incorrect timestamp")
	}
}
