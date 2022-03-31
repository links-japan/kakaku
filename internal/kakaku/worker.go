package kakaku

import (
	"context"
	"fmt"

	"github.com/links-japan/kakaku/internal/store"
	"github.com/links-japan/log"
)

func UpdateAssetPrice(oracle *Oracle, assets *store.AssetStore, base, quote string) error {
	var asset store.Asset
	err := assets.Find(&asset, base, quote)
	if err != nil {
		return err
	}
	term := asset.Term + 1
	log.WithField("term", term).Debug("update asset price")

	nullPrice, source := oracle.Price(context.TODO(), base, quote)
	if !nullPrice.Valid {
		return fmt.Errorf("failed term: %v, base: %v, quote %v\n", term, base, quote)
	}

	asset.Price = nullPrice.Decimal
	asset.Term = term
	asset.Source = source

	return assets.Update(&asset)
}
