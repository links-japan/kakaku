package client

import (
	"context"
	"github.com/shopspring/decimal"
)

const (
	BTC  = "BTC"
	JPYC = "JPYC"

	JPY = "JPY"
)

type Client interface {
	Price(ctx context.Context, base string, quote string) (decimal.Decimal, error)
	Source() string
}
