package main

import (
	"fmt"
	"sort"
	"time"
)

type Transaction struct {
	id         string
	token0Path string
	token1Path string
	amount0    int
	amount1    int
	time       time.Time
}

func main() {
	layout := "2006-01-02 15:04:05"
	transactions := []Transaction{
		{"ccb4668d", "gno.land/r/demo/gns", "gno.land/r/demo/wugnot", 100000, -685659, parseTime("2024-05-16 05:21:17", layout)},
		{"58f51962", "gno.land/r/demo/gns", "gno.land/r/demo/wugnot", 10000, -68658, parseTime("2024-05-16 05:20:34", layout)},
		{"b4c2a9c0", "gno.land/r/demo/wugnot", "gno.land/r/demo/gns", 27000000, -18399281, parseTime("2024-05-16 05:15:00", layout)},
		{"b8a0ad7d", "gno.land/r/demo/wugnot", "gno.land/r/demo/gns", 2000000, -3771726, parseTime("2024-05-16 05:14:17", layout)},
		{"65d7ad35", "gno.land/r/demo/bar", "gno.land/r/demo/gns", 1000000, -131195131, parseTime("2024-05-16 05:05:14", layout)},
		{"c06cdf98", "gno.land/r/demo/gns", "gno.land/r/demo/bar", 10000, -19961, parseTime("2024-05-16 05:04:51", layout)},
		{"792098bf", "gno.land/r/demo/baz", "gno.land/r/demo/gns", 1000000, -6437928, parseTime("2024-05-16 04:42:41", layout)},
		{"6d07c81c", "gno.land/r/demo/foo", "gno.land/r/demo/gns", 50000, -96865, parseTime("2024-05-16 02:01:28", layout)},
		{"a16085c3", "gno.land/r/demo/wugnot", "gno.land/r/demo/gns", 245, -242450006, parseTime("2024-05-14 14:29:22", layout)},
		{"389b0fa9", "gno.land/r/demo/wugnot", "gno.land/r/demo/gns", 2, -9985, parseTime("2024-05-14 14:28:47", layout)},
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].time.Before(transactions[j].time)
	})

	priceMap := make(map[string]float64)
	priceMap["gno.land/r/demo/wugnot"] = 1.0
	priceMap["gno.land/r/demo/gns"] = 0.0
	priceMap["gno.land/r/demo/bar"] = 0.0
	priceMap["gno.land/r/demo/baz"] = 0.0
	priceMap["gno.land/r/demo/foo"] = 0.0

	priceHistory := calculatePriceHistory(transactions, priceMap)

	for _, entry := range priceHistory {
		fmt.Printf("Time: %s\n", entry.time.Format(layout))
		for token, price := range entry.prices {
			fmt.Printf("%s: $%.4f\n", token, price)
		}
		fmt.Println("-----------")
	}
}

func parseTime(timeStr, layout string) time.Time {
	t, _ := time.Parse(layout, timeStr)
	return t
}

func calculatePriceHistory(transactions []Transaction, initialPrices map[string]float64) []PriceEntry {
	var priceHistory []PriceEntry
	currentPrices := make(map[string]float64)
	for token, price := range initialPrices {
		currentPrices[token] = price
	}

	start := transactions[0].time

	for current := start; current.Before(transactions[len(transactions)-1].time); current = current.Add(10 * time.Minute) {
		for _, tx := range transactions {
			if tx.time.Before(current.Add(10 * time.Minute)) {
				updatePrices(tx, currentPrices)
			}
		}
		priceHistory = append(priceHistory, PriceEntry{time: current, prices: copyMap(currentPrices)})
	}

	return priceHistory
}

func updatePrices(tx Transaction, prices map[string]float64) {
	var exchangeRate float64
	if tx.token0Path == "gno.land/r/demo/wugnot" || tx.token1Path == "gno.land/r/demo/wugnot" {
		if tx.token0Path == "gno.land/r/demo/wugnot" {
			exchangeRate = float64(tx.amount0) / float64(-tx.amount1)
			prices[tx.token1Path] = prices[tx.token0Path] / exchangeRate
		} else {
			exchangeRate = float64(-tx.amount1) / float64(tx.amount0)
			prices[tx.token0Path] = prices[tx.token1Path] * exchangeRate
		}
	} else if tx.token0Path == "gno.land/r/demo/gns" || tx.token1Path == "gno.land/r/demo/gns" {
		if tx.token0Path == "gno.land/r/demo/gns" {
			exchangeRate = float64(tx.amount0) / float64(-tx.amount1)
			if price, exists := prices["gno.land/r/demo/gns"]; exists && price != 0 {
				prices[tx.token1Path] = price / exchangeRate
			}
		} else {
			exchangeRate = float64(-tx.amount1) / float64(tx.amount0)
			if price, exists := prices["gno.land/r/demo/gns"]; exists && price != 0 {
				prices[tx.token0Path] = price * exchangeRate
			}
		}
	}
}

func copyMap(original map[string]float64) map[string]float64 {
	newMap := make(map[string]float64)
	for key, value := range original {
		newMap[key] = value
	}
	return newMap
}

type PriceEntry struct {
	time   time.Time
	prices map[string]float64
}
