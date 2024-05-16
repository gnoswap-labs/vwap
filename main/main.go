package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gnoswap-labs/vwap"
)

// testing purposes

func main() {
	start := time.Now().Unix()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		<-ticker.C
		vwapResults, err := vwap.VWAP()
		if err != nil {
			log.Printf("Failed to calculate VWAP: %v", err)
			continue
		}

		fmt.Println("VWAP Results:")
		for token, vwap := range vwapResults {
			fmt.Printf("%s: %.8f\n", token, vwap)
		}

		fmt.Println(time.Now().Unix() - start)
		fmt.Println("--------------------")
	}
}
