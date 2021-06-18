package client

import (
	"context"
	"encoding/json"
	"github.com/shopspring/decimal"
	"net/http"
)

type CoinbaseClient struct {
	client *http.Client
}

func NewCoinBaseClient() *CoinbaseClient {
	return &CoinbaseClient{
		client: &http.Client{},
	}
}


type priceResponse struct {
	Data struct {
		Base     string
		Currency string
		Amount   decimal.Decimal
	}
}

func (co *CoinbaseClient) Price(ctx context.Context, base string, quote string) (decimal.Decimal, error) {
	pair := base + "-" + quote
	uri := "https://api.coinbase.com/v2/prices/" + pair + "/spot"
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return decimal.Zero, err
	}

	resp, err := co.client.Do(req)
	if err != nil {
		return decimal.Zero, err
	}
	defer resp.Body.Close()

	var re priceResponse
	if err = json.NewDecoder(resp.Body).Decode(&re); err != nil {
		return decimal.Zero, err
	}

	return  re.Data.Amount, nil
}

func (co *CoinbaseClient) Name() string {
	return "CoinbaseClient"
}
