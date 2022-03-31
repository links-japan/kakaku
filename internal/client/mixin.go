package client

import (
	"context"
	"strings"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/shopspring/decimal"
)

var Symbol2AssetID = map[string]string{
	"BTC":  "c6d0c728-2624-429b-8e0d-d9d19b6592fa",
	"CNB":  "965e5c6e-434c-3fa9-b780-c50f43cd955c",
	"ETH":  "43d61dcd-e413-450d-80b8-101d5e903357",
	"JPYC": "0ff3f325-4f34-334d-b6c0-a3bd8850fc06",
}

type MixinClient struct {
	client *mixin.Client
}

func NewMixinClient(keystore *mixin.Keystore) *MixinClient {
	cli, err := mixin.NewFromKeystore(keystore)
	if err != nil {
		panic(err)
	}
	return &MixinClient{cli}
}

func (m *MixinClient) Price(ctx context.Context, base string, quote string) (decimal.Decimal, error) {

	assetID := Symbol2AssetID[base]

	asset, err := m.client.ReadAsset(ctx, assetID)
	if err != nil {
		return decimal.Zero, err
	}

	rates, err := m.client.ReadExchangeRates(ctx)
	if err != nil {
		return decimal.Zero, err
	}

	var rate decimal.Decimal
	for _, r := range rates {
		if strings.EqualFold(r.Code, quote) {
			rate = r.Rate
		}
	}

	return asset.PriceUSD.Mul(rate), nil
}

func (m *MixinClient) Source() string {
	return "Mixin"
}
