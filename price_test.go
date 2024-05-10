package vwap

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchTokenPrices(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenPricesResponse := APIResponse{
			Data: []TokenPrice{
				{
					Path: "gno.land/r/demo/bar",
					USD:  "30.087345",
					PricesBefore: PricesBefore{
						LatestPrice: "30.087345",
						Price1h:     "30.087345",
						PriceToday:  "30.087345",
						Price1d:     "30.087345",
						Price7d:     "0.000000",
						Price30d:    "0.000000",
						Price60d:    "0.000000",
						Price90d:    "0.000000",
					},
					MarketCap:         "15043681760.213797",
					LockedTokensUSD:   "351004463.170473",
					VolumeUSD24h:      "0.000000",
					FeeUSD24h:         "0",
					MostLiquidityPool: "gno.land/r/demo/bar:gno.land/r/demo/baz:100",
					Last7d: []Last7d{
						{
							Date:  "2024-05-09T08:00:00Z",
							Price: "30.087345468020313",
						},
						{
							Date:  "2024-05-09T07:00:00Z",
							Price: "30.087345468020313",
						},
					},
				},
			},
		}

		err := json.NewEncoder(w).Encode(tokenPricesResponse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}))
	defer server.Close()

	apiEndpoint := server.URL

	prices, err := fetchTokenPrices(apiEndpoint)
	assert.NoError(t, err, "Failed to fetch token prices")

	expectedPrices := []TokenPrice{
		{
			Path: "gno.land/r/demo/bar",
			USD:  "30.087345",
			PricesBefore: PricesBefore{
				LatestPrice: "30.087345",
				Price1h:     "30.087345",
				PriceToday:  "30.087345",
				Price1d:     "30.087345",
				Price7d:     "0.000000",
				Price30d:    "0.000000",
				Price60d:    "0.000000",
				Price90d:    "0.000000",
			},
			MarketCap:         "15043681760.213797",
			LockedTokensUSD:   "351004463.170473",
			VolumeUSD24h:      "0.000000",
			FeeUSD24h:         "0",
			MostLiquidityPool: "gno.land/r/demo/bar:gno.land/r/demo/baz:100",
			Last7d: []Last7d{
				{
					Date:  "2024-05-09T08:00:00Z",
					Price: "30.087345468020313",
				},
				{
					Date:  "2024-05-09T07:00:00Z",
					Price: "30.087345468020313",
				},
			},
		},
	}

	assert.Equal(t, expectedPrices, prices, "Fetched token prices do not match expected prices")
}

func TestFetchTokenPricesLive(t *testing.T) {
	t.Parallel()

	prices, err := fetchTokenPrices(priceEndpoint)
	if err != nil {
		t.Fatalf("Failed to fetch token prices: %v", err)
	}

	if len(prices) != 6 {
		t.Errorf("Expected 6 token prices, got %d", len(prices))
	}

	for _, price := range prices {
		if price.Path == "" {
			t.Errorf("Token path is empty")
		}
		if price.USD == "" {
			t.Errorf("Token USD price is empty")
		}
		if price.MarketCap == "" {
			t.Errorf("Token market cap is empty")
		}
		if price.VolumeUSD24h == "" {
			t.Errorf("Token 24h volume USD is empty")
		}
	}
}

func TestCalculateVolume(t *testing.T) {
	t.Parallel()
	mockTokenPrices := []TokenPrice{
		{
			Path:         "token1",
			VolumeUSD24h: "1000.50",
		},
		{
			Path:         "token2",
			VolumeUSD24h: "2500.75",
		},
		{
			Path:         "token3",
			VolumeUSD24h: "500.25",
		},
	}

	volumes := calculateVolume(mockTokenPrices)

	assert.Equal(t, 3, len(volumes))
	assert.InDelta(t, 1000.50, volumes["token1"], 0.001)
	assert.InDelta(t, 2500.75, volumes["token2"], 0.001)
	assert.InDelta(t, 500.25, volumes["token3"], 0.001)
}

func TestCalculateTokenPrices(t *testing.T) {
	t.Parallel()
	mockTokenData := &TokenData{
		Status: struct {
			Height    int `json:"height"`
			Timestamp int `json:"timestamp"`
		}{
			Height:    61946,
			Timestamp: 1715064552,
		},
		TokenRatio: []struct {
			TokenName string `json:"token"`
			Ratio     string `json:"ratio"`
		}{
			{
				TokenName: "gno.land/r/demo/wugnot",
				Ratio:     "79228162514264337593543950336",
			},
			{
				TokenName: "gno.land/r/demo/qux",
				Ratio:     "60536002769587966558221762891",
			},
			{
				TokenName: "gno.land/r/demo/foo",
				Ratio:     "121074204438706762251182081654",
			},
			{
				TokenName: "gno.land/r/demo/gns",
				Ratio:     "236174327852992866806677676716",
			},
		},
	}

	baseTokenPrice := 1.0
	tokenPrices := calculateTokenUSDPrices(mockTokenData, baseTokenPrice)

	assert.InDelta(t, 1.3087775685458196, tokenPrices["gno.land/r/demo/qux"], 1e-5)
	assert.InDelta(t, 0.6543768995349725, tokenPrices["gno.land/r/demo/foo"], 1e-5)
	assert.InDelta(t, 0.3354647528142011, tokenPrices["gno.land/r/demo/gns"], 1e-5)
}
