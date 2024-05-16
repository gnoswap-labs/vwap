package vwap

import (
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestStoreWithGORM(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&VWAPData{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	err = store(db, "FOO", 50000.0, 1000.0, time.Now())
	if err != nil {
		t.Errorf("error was not expected while storing data: %s", err)
	}

	var vwapData VWAPData
	db.First(&vwapData)

	if vwapData.TokenName != "FOO" || vwapData.VWAP != 50000.0 || vwapData.TotalVolume != 1000.0 {
		t.Errorf("stored data is incorrect")
	}
}

func TestVWAPWith10MinuteInterval(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&VWAPData{})
	assert.NoError(t, err)

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

			actualVWAP, err := calculateVWAP(db, intervalTrades)
			assert.Nil(t, err, "Unexpected error")

			actualVWAPs = append(actualVWAPs, actualVWAP)

			intervalStart = trade.Timestamp
			intervalTrades = []TradeData{trade}
		}
	}

	// Process the last interval
	expectedVWAP := calculateExpectedVWAP(intervalTrades)
	expectedVWAPs = append(expectedVWAPs, expectedVWAP)

	actualVWAP, err := calculateVWAP(db, intervalTrades)
	assert.Nil(t, err, "Unexpected error")

	actualVWAPs = append(actualVWAPs, actualVWAP)

	// Compare expected and actual VWAPs for each interval
	assert.Equal(t, len(expectedVWAPs), len(actualVWAPs), "Unexpected number of intervals")

	for i := 0; i < len(expectedVWAPs); i++ {
		assert.InDelta(t, expectedVWAPs[i], actualVWAPs[i], 1e-9, "Incorrect VWAP for interval %d", i+1)
	}

	var vwapDataList []VWAPData
	db.Find(&vwapDataList)

	assert.Equal(t, len(expectedVWAPs), len(vwapDataList), "Unexpected number of stored VWAP data")

	for i, vwapData := range vwapDataList {
		assert.Equal(t, "Token1", vwapData.TokenName)
		assert.InDelta(t, expectedVWAPs[i], vwapData.VWAP, 1e-9)
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
