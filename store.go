package vwap

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

/* Schema:
* CREATE TABLE vwap_data (
*     id SERIAL PRIMARY KEY,
*     token_name VARCHAR(50) NOT NULL,
*     calculated_at TIMESTAMP NOT NULL DEFAULT NOW(),
*     vwap DECIMAL(20, 10) NOT NULL,
*     total_volume DECIMAL(20, 10) NOT NULL
* );
 */

type VWAPData struct {
	gorm.Model
	TokenName    string
	VWAP         float64
	TotalVolume  float64
	CalculatedAt time.Time
}

func store(db *gorm.DB, tokenName string, vwap, totalVolume float64, calculatedAt time.Time) error {
	vwapData := VWAPData{
		TokenName:    tokenName,
		VWAP:         vwap,
		TotalVolume:  totalVolume,
		CalculatedAt: calculatedAt,
	}

	result := db.Create(&vwapData)
	if result.Error != nil {
		return fmt.Errorf("failed to insert data: %v", result.Error)
	}

	return nil
}

// testing purpose

func PopulateVWAPData(db *gorm.DB, count int) error {
	for i := 0; i < count; i++ {
		vwapData := VWAPData{
			TokenName:    faker.Currency(),
			VWAP:         rand.Float64(),
			TotalVolume:  rand.Float64(),
			CalculatedAt: time.Now().Add(time.Duration(rand.Intn(1000)) * time.Minute),
		}

		result := db.Create(&vwapData)
		if result.Error != nil {
			return fmt.Errorf("failed to insert data: %v", result.Error)
		}
	}

	return nil
}
