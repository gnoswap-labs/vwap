package vwap

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

// Token name
type TokenIdentifier string

const (
	WUGNOT TokenIdentifier = "gno.land/r/demo/wugnot" // base token (ratio = 1)
	FOO    TokenIdentifier = "gno.land/r/demo/foo"
	BAR    TokenIdentifier = "gno.land/r/demo/bar"
	BAZ    TokenIdentifier = "gno.land/r/demo/baz"
	QUX    TokenIdentifier = "gno.land/r/demo/qux"
	GNS    TokenIdentifier = "gno.land/r/demo/gns"
)

// TradeData represents the data for a single trade.
type TradeData struct {
	TokenName string
	Volume    float64
	Ratio     float64
	Timestamp int
}

// lastPrices stores the last price of each token.
// This value will be used to show the last price if the token is not traded.
var (
	lastPrices      map[string]float64
	lastPricesMutex sync.Mutex
)

func init() {
	lastPrices = make(map[string]float64)
}

func VWAP(db *gorm.DB) (map[string]float64, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	prices, err := fetchTokenPrices(priceEndpoint)
	if err != nil {
		return nil, err
	}

	volumeByToken := calculateVolume(prices)
	trades := extractTrades(prices, volumeByToken)
	vwapResults := make(map[string]float64)

	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
	)

	for tokenName, tradeData := range trades {
		wg.Add(1)
		go func(tokenName string, tradeData []TradeData) {
			defer wg.Done()
			res, err := calculateVWAP(db, tradeData)
			if err != nil {
				log.Printf("failed to calculate VWAP for token %s: %v\n", tokenName, err)
				return
			}
			mutex.Lock()
			vwapResults[tokenName] = res
			mutex.Unlock()
		}(tokenName, tradeData)
	}

	wg.Wait()

	return vwapResults, nil
}

// calculateVWAP calculates the Volume Weighted Average Price (calculateVWAP) for the given set of trades.
// It returns the last price if there are no trades.
func calculateVWAP(db *gorm.DB, trades []TradeData) (float64, error) {
	var numerator, denominator float64

	if len(trades) == 0 {
		return 0, fmt.Errorf("no trades found")
	}

	for _, trade := range trades {
		numerator += trade.Volume * trade.Ratio
		denominator += trade.Volume
	}

	// return last price if there is no trade
	if denominator == 0 {
		lastPricesMutex.Lock()
		lastPrice, ok := lastPrices[trades[0].TokenName]
		lastPricesMutex.Unlock()
		if !ok {
			return 0, nil
		}
		return lastPrice, nil
	}

	vwap := numerator / denominator
	lastPricesMutex.Lock()
	lastPrices[trades[0].TokenName] = vwap // save the last price
	lastPricesMutex.Unlock()

	totalVolume := denominator
	calculatedAt := time.Now()

	err := store(db, trades[0].TokenName, vwap, totalVolume, calculatedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to store data: %v", err)
	}

	return vwap, nil
}
