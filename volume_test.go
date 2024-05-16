package vwap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateVolume(t *testing.T) {
	// Test case 1: Valid token prices
	prices := []TokenPrice{
		{Path: "TOKEN1", VolumeUSD24h: "1000.50"},
		{Path: "TOKEN2", VolumeUSD24h: "500.25"},
	}
	volumeByToken := calculateVolume(prices)
	assert.Equal(t, 1000.50, volumeByToken["TOKEN1"])
	assert.Equal(t, 500.25, volumeByToken["TOKEN2"])

	// Test case 2: Invalid volume format
	prices = []TokenPrice{
		{Path: "TOKEN3", VolumeUSD24h: "invalid"},
	}
	volumeByToken = calculateVolume(prices)
	assert.Equal(t, 0.0, volumeByToken["TOKEN3"])
}
