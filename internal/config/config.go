package config

import (
	"github.com/shopspring/decimal"
	"time"
)

type (
	Config struct {
		Mixin  Mixin  `mapstructure:"mixin"`
		Oracle Oracle `mapstructure:"oracle"`
		Worker Worker `mapstructure:"worker"`
		DB     DB     `mapstructure:"db"`
		Debug  bool   `mapstructure:"debug"`
	}

	DB struct {
		Dsn string `json:"dsn"`
	}

	Mixin struct {
		ClientID   string `mapstructure:"client_id"`
		SessionID  string `mapstructure:"session_id"`
		PrivateKey string `mapstructure:"private_key"`
		PinToken   string `mapstructure:"pin_token"`
	}

	Oracle struct {
		RequestTimeout   time.Duration `mapstructure:"request_timeout"`
		ApproveThreshold int           `mapstructure:"approve_threshold"`
		PriceDeltaStr    string        `mapstructure:"price_delta"`
		PriceDelta       decimal.Decimal
	}

	Worker struct {
		TermTimeout time.Duration `mapstructure:"term_timeout"`
	}
)
