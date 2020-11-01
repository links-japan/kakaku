package kakaku

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
	err := tx.Where(Asset{Symbol: symbol}).First(&asset).Error
	switch {
	case err == gorm.ErrRecordNotFound:
		asset.Symbol = symbol
		asset.PriceJPY = decimal.Zero
		if err := tx.Create(&asset).Error; err != nil {
			return nil, err
		}
		return &asset, nil
	case err != nil:
		return nil, err
	}
	return &asset, nil
}
