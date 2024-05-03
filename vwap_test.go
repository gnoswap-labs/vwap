package vwap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVWAP(t *testing.T) {
	tests := []struct {
		name         string
		trades       []Trade
		expectedVWAP float64
	}{
		{
			name: "Normal case",
			trades: []Trade{
				{Price: 100.0, Volume: 100},
				{Price: 101.0, Volume: 200},
				{Price: 102.0, Volume: 150},
				{Price: 103.0, Volume: 50},
			},
			expectedVWAP: 101.3,
		},
		{
			name:         "Empty trades",
			trades:       []Trade{},
			expectedVWAP: 0.0,
		},
		{
			name: "Single trade",
			trades: []Trade{
				{Price: 100.0, Volume: 100},
			},
			expectedVWAP: 100.0,
		},
		{
			name: "Zero volume trades",
			trades: []Trade{
				{Price: 100.0, Volume: 0},
				{Price: 101.0, Volume: 0},
			},
			expectedVWAP: 0.0,
		},
		{
			name: "Large volume",
			trades: []Trade{
				{Price: 100.0, Volume: 1000000},
				{Price: 101.0, Volume: 2000000},
			},
			expectedVWAP: 100.666667,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vwap := VWAP(tt.trades)
			assert.InDelta(t, tt.expectedVWAP, vwap, 0.000001, "Expected VWAP")
		})
	}
}
