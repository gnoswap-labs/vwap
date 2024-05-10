package main

import (
	"log"
	"sort"
	"time"

	"github.com/gnoswap-labs/vwap"
)

// testing purposes

func main() {
	interval := time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	calculateAndPrintVWAP()

	for range ticker.C {
		calculateAndPrintVWAP()
	}
}

func calculateAndPrintVWAP() {
	vwapResults, err := vwap.FetchAndCalculateVWAP()
	if err != nil {
		log.Printf("Error fetching and calculating VWAP: %v\n", err)
		return
	}

	var tokens []string
	for token := range vwapResults {
		tokens = append(tokens, token)
	}
	sort.Strings(tokens)

	log.Println("VWAP results updated at", time.Now().Format("15:04:05"))

	for _, token := range tokens {
		log.Printf("Token: %s, VWAP: %.4f\n", token, vwapResults[token])
	}
}
