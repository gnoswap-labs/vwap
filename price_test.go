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
		tokenPricesResponse := PricesResponse{
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
