package main

import (
	"fmt"
	"log"

	"github.com/gnoswap-labs/vwap"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("vwap.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&vwap.VWAPData{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	err = vwap.PopulateVWAPData(db, 10)
	if err != nil {
		log.Fatalf("failed to populate vwap data: %v", err)
	}

	var VWAPDataList []vwap.VWAPData
	result := db.Find(&VWAPDataList)
	if result.Error != nil {
		log.Fatalf("failed to retrieve vwap data: %v", result.Error)
	}

	for _, VWAPData := range VWAPDataList {
		fmt.Printf("TokenName: %s, VWAP: %.2f, TotalVolume: %.2f, CalculatedAt: %s\n",
			VWAPData.TokenName, VWAPData.VWAP, VWAPData.TotalVolume, VWAPData.CalculatedAt)
	}
}
