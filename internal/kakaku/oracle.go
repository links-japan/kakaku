package kakaku

import (
	"context"
	"sync"
	"time"

	"github.com/links-japan/kakaku/internal/client"
	"github.com/links-japan/kakaku/internal/config"
	"github.com/links-japan/kakaku/internal/store"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type Oracle struct {
	clients []client.Client
	assets  *store.AssetStore
	cfg     *config.Oracle
}

func NewOracle(clients []client.Client, assets *store.AssetStore, cfg *config.Oracle) *Oracle {
	return &Oracle{
		clients: clients,
		assets:  assets,
		cfg:     cfg,
	}
}

func (o *Oracle) Price(ctx context.Context, base string, quote string) (decimal.NullDecimal, string) {
	result := decimal.NullDecimal{}
	approveCnt := 0
	start := time.Now()
	mu := sync.Mutex{}
	source := ""

	ctx, cancel := context.WithTimeout(ctx, o.cfg.RequestTimeout)
	defer cancel()

	for _, cli := range o.clients {
		go func(cli client.Client) {
			price, err := cli.Price(ctx, base, quote)
			logrus.WithField("name", cli.Source()).WithField("price", price).Debug("client price")

			if err != nil {
				logrus.Error("client price err", cli.Source(), err)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			if !result.Valid {
				result.Decimal = price
				result.Valid = true
				source = cli.Source()
				approveCnt += 1
				return
			}

			value := result.Decimal
			if value == decimal.Zero {
				result.Valid = false
				return
			}
			delta := value.Sub(price).Abs().Div(value)
			if delta.LessThan(o.cfg.PriceDelta) {
				approveCnt += 1
			}
			return
		}(cli)
	}

	for time.Since(start) < o.cfg.RequestTimeout {
		if approveCnt >= o.cfg.ApproveThreshold {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if approveCnt < o.cfg.ApproveThreshold {
		result.Valid = false
	}
	return result, source
}
