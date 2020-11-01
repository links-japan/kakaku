package kakaku

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	LIQUID_API_BASE    = "https://api.liquid.com" // Liquid API endpoint
	LIQUID_API_VERSION = "2"
	BTC_JPY_PAIR       = 5
)

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

type liquid struct {
	httpClient  *http.Client
	httpTimeout time.Duration
}

// newClient return a new Liquid HTTP client
func newClient() (c *liquid) {
	return &liquid{&http.Client{}, 30 * time.Second}
}

func (l *liquid) getTicker(productId int) (*product, error) {
	r, err := l.do("GET", "/products/"+strconv.Itoa(productId), "")
	if err != nil {
		return nil, err
	}

	var p product
	if err = json.Unmarshal(r, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

// do prepare and process HTTP request to Liquid API
func (l *liquid) do(method string, resource string, payload string) (response []byte, err error) {
	connectTimer := time.NewTimer(l.httpTimeout)

	var rawurl string
	if strings.HasPrefix(resource, "http") {
		rawurl = resource
	} else {
		rawurl = fmt.Sprintf("%s%s", LIQUID_API_BASE, resource)
	}

	req, err := http.NewRequest(method, rawurl, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return
	}

	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Quoine-API-Version", LIQUID_API_VERSION)

	resp, err := l.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("liquid response error, status code %d, status %s", resp.StatusCode, resp.Status)
	}
	return response, err
}

// doTimeoutRequest do a HTTP request with timeout
func (l *liquid) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		resp, err := l.httpClient.Do(req)
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, fmt.Errorf("timeout on reading data from Liquid API")
	}
}
