package vwap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVWAPStorage(t *testing.T) {
	vwapDataMap = make(map[string][]VWAPData)
	lastPrices = make(map[string]float64)

	trades := generateMockTrades()

	for _, trade := range trades {
		VWAP([]TradeData{trade})
	}

	expectedVWAPData := map[string][]VWAPData{
		"gno.land/r/demo/foo": {
			{TokenName: "gno.land/r/demo/foo", VWAP: 1.2, Timestamp: 1623200400},
			{TokenName: "gno.land/r/demo/foo", VWAP: 1.5, Timestamp: 1623201000},
		},
		"gno.land/r/demo/bar": {
			{TokenName: "gno.land/r/demo/bar", VWAP: 2.1, Timestamp: 1623200400},
			{TokenName: "gno.land/r/demo/bar", VWAP: 2.3, Timestamp: 1623201000},
		},
	}

	assert.Equal(t, expectedVWAPData, vwapDataMap)
}

func generateMockTrades() []TradeData {
	return []TradeData{
		{TokenName: "gno.land/r/demo/foo", Volume: 100, Ratio: 1.2, Timestamp: 1623200400},
		{TokenName: "gno.land/r/demo/bar", Volume: 200, Ratio: 2.1, Timestamp: 1623200400},
		{TokenName: "gno.land/r/demo/foo", Volume: 150, Ratio: 1.5, Timestamp: 1623201000},
		{TokenName: "gno.land/r/demo/bar", Volume: 250, Ratio: 2.3, Timestamp: 1623201000},
	}
}
