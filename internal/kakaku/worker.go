package kakaku

import (
	"context"
	"fmt"
	"github.com/links-japan/kakaku/internal/store"
	"github.com/sirupsen/logrus"
)

func UpdateAssetPrice(oracle *Oracle, assets *store.AssetStore, base, quote string) error {
	var asset store.Asset
	err := assets.Find(&asset, base, quote)
	if err != nil {
		return err
	}
	prevTerm := asset.Term
	term := prevTerm + 1
	logrus.WithField("term", term).Debug("update asset price")

	nullPrice, source := oracle.Price(context.TODO(), base, quote)
	if !nullPrice.Valid {
		return fmt.Errorf("failed term: %v, base: %v, quote %v\n", term, base, quote)
	}

	asset.Price = nullPrice.Decimal
	asset.Term = term
	asset.Source = source

	return assets.Update(&asset, base, quote, prevTerm)
}
