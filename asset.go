package main

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Asset struct {
	gorm.Model
	Symbol   string          `gorm:"unique"`
	PriceJPY decimal.Decimal `gorm:"not null;" sql:"type:decimal(8,0);"`
}

func FirstOrCreate(symbol string, tx *gorm.DB) (*Asset, error) {
	var asset Asset
	err := tx.Where(Asset{Symbol: symbol}).FirstOrCreate(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}
