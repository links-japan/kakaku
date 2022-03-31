package store

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	AssetConstType    = "Const"
	AssetVariableType = "Variable"
)

type Asset struct {
	ID        int64           `json:"id"`
	Base      string          `json:"base" sql:"size:32"`
	Quote     string          `json:"quote" sql:"size:32"`
	Source    string          `json:"source" sql:"size:32"`
	Price     decimal.Decimal `json:"price" gorm:"not null;" sql:"type:decimal(8,0);"`
	Term      int64           `json:"term"`
	Type      string          `json:"type" sql:"size:32"`
	UpdatedAt time.Time       `json:"updated_at"`
	CreatedAt time.Time       `json:"created_at"`
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

func (a *AssetStore) Update(asset *Asset) error {
	return a.tx.Model(asset).
		Where("id = ? ", asset.ID).
		Updates(map[string]interface{}{"term": asset.Term, "price": asset.Price, "source": asset.Source}).Error
}
