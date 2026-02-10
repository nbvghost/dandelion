package model

import (
	"log"

	"github.com/shopspring/decimal"
)

func init() {
	decimal.MarshalJSONWithoutQuotes = true
	decimal.DivisionPrecision = 6
	log.Println("DivisionPrecision:", decimal.DivisionPrecision)
	log.Println("MarshalJSONWithoutQuotes:", decimal.MarshalJSONWithoutQuotes)
}
