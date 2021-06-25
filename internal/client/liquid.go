package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	LiquidApiBase    = "https://api.liquid.com" // LiquidClient API endpoint
	LiquidApiVersion = "2"
)

var PairProductID = map[string]int{
	"BTC-JPY": 5,
}

type product struct {
	ID                 string          `json:"id"`
	Symbol             string          `json:"symbol"`
	ProductType        string          `json:"product_type"`
	Code               string          `json:"code"`
	Name               string          `json:"name"`
	MarketAsk          decimal.Decimal `json:"market_ask"`
	MarketBid          decimal.Decimal `json:"market_bid"`
	Currency           string          `json:"currency"`
	CurrencyPairCode   string          `json:"currency_pair_code"`
	LastTradedPrice    decimal.Decimal `json:"last_traded_price"`
	QuotedCurrency     string          `json:"quoted_currency"`
	BaseCurrency       string          `json:"base_currency"`
	LastEventTimestamp string          `json:"last_event_timestamp"`
	ExchangeRate       int             `json:"exchange_rate"`
}

type LiquidClient struct {
	httpClient *http.Client
}

// NewLiquidClient return a new LiquidClient HTTP client
func NewLiquidClient() (c *LiquidClient) {
	return &LiquidClient{&http.Client{}}
}

func (l *LiquidClient) Price(ctx context.Context, base string, quote string) (decimal.Decimal, error) {
	pair := base + "-" + quote
	p, err := l.getTicker(ctx, PairProductID[pair])
	if err != nil {
		return decimal.Zero, err
	}
	return p.LastTradedPrice, nil
}

func (l *LiquidClient) getTicker(ctx context.Context, productId int) (*product, error) {
	r, err := l.do(ctx, "GET", "/products/"+strconv.Itoa(productId), "")
	if err != nil {
		return nil, err
	}

	var p product
	if err = json.Unmarshal(r, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

// do prepare and process HTTP request to LiquidClient API
func (l *LiquidClient) do(ctx context.Context, method string, resource string, payload string) (response []byte, err error) {
	var rawurl string
	if strings.HasPrefix(resource, "http") {
		rawurl = resource
	} else {
		rawurl = fmt.Sprintf("%s%s", LiquidApiBase, resource)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawurl, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return
	}

	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Quoine-API-Version", LiquidApiVersion)

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("LiquidClient response error, status code %d, status %s", resp.StatusCode, resp.Status)
	}
	return response, err
}

func (l *LiquidClient) Source() string {
	return "Liquid"
}
