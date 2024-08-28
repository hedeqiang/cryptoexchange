package types

import (
	"net/http"
)

type ExchangeName string

const (
	Binance  ExchangeName = "BINANCE"
	OKX      ExchangeName = "OKX"
	Bitget   ExchangeName = "BITGET"
	Kucoin   ExchangeName = "KUCOIN"
	MEXC     ExchangeName = "MEXC"
	Gate     ExchangeName = "GATE"
	Kraken   ExchangeName = "KRAKEN"
	Bybit    ExchangeName = "BYBIT"
	Huobi    ExchangeName = "HUOBI"
	Coinbase ExchangeName = "COINBASE"
	BTSE     ExchangeName = "BTSE"
)

type ExchangeConfig struct {
	APIKey        string
	APISecret     string
	BaseURL       string
	APIPassphrase string
}

type Exchange interface {
	Name() ExchangeName
	GetDefaultBaseURL() string
	PrepareRequest(method, endpoint string, params map[string]interface{}, signed bool) (*http.Request, error)
}
