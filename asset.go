package main

import "github.com/shopspring/decimal"

type Asset struct {
	//gorm.Model
	Symbol   string
	PriceJPY decimal.Decimal
}
