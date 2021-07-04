package store

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

const (
	AssetConstType    = "Const"
	AssetVariableType = "Variable"
)

type Asset struct {
	ID        int64
	Base      string          `sql:"size:32"`
	Quote     string          `sql:"size:32"`
	Source    string          `sql:"size:32"`
	Price     decimal.Decimal `gorm:"not null;" sql:"type:decimal(8,0);"`
	Term      int64
	Type      string `sql:"size:32"`
	UpdatedAt time.Time
	CreatedAt time.Time
}

type AssetStore struct {
	tx *gorm.DB
}

func NewAssetStore() *AssetStore {
	return &AssetStore{tx: db}
}

func (a *AssetStore) FirstOrCreate(asset *Asset) error {
	return a.tx.FirstOrCreate(asset).Error
}

func (a *AssetStore) Find(asset *Asset, base, quote string) error {
	return a.tx.Where("base = ? AND quote = ?", base, quote).First(asset).Error
}

func (a *AssetStore) ListAll() ([]*Asset, error) {
	var assets []*Asset
	if err := a.tx.Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

func (a *AssetStore) ListVariable() ([]*Asset, error) {
	var assets []*Asset
	if err := a.tx.Where("type = ?", AssetVariableType).Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

func (a *AssetStore) Update(asset *Asset, base, quote string, prevTerm int64) error {
	return a.tx.Model(asset).
		Where("base = ? AND quote = ? AND term = ?", base, quote, prevTerm).
		Updates(map[string]interface{}{"term": asset.Term, "price": asset.Price, "source": asset.Source}).Error
}
