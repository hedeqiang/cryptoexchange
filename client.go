package cryptoexchange

import (
	"encoding/json"
	"fmt"
	"github.com/hedeqiang/cryptoexchange/exchanges"
	"github.com/hedeqiang/cryptoexchange/types"
	"io/ioutil"
	"net/http"
)

type CryptoExchangeClient struct {
	exchange types.Exchange
	client   *http.Client
}

func NewCryptoExchangeClient() *CryptoExchangeClient {
	return &CryptoExchangeClient{
		exchange: nil,
		client:   &http.Client{},
	}
}

func (c *CryptoExchangeClient) AddExchange(name types.ExchangeName, config types.ExchangeConfig) error {

	switch name {
	case types.Binance:
		c.exchange = exchanges.NewBinance(config)
	case types.OKX:
		c.exchange = exchanges.NewOKX(config)
	case types.Bitget:
		c.exchange = exchanges.NewBitget(config)
	case types.Kucoin:
		c.exchange = exchanges.NewKucoin(config)
	case types.MEXC:
		c.exchange = exchanges.NewMEXC(config)
	case types.Gate:
		c.exchange = exchanges.NewGate(config)
	case types.Kraken:
		c.exchange = exchanges.NewKraken(config)
	case types.Bybit:
		c.exchange = exchanges.NewBybit(config)
	case types.Huobi:
		c.exchange = exchanges.NewHuobi(config)
	default:
		return &ExchangeError{Exchange: name, Message: "unsupported exchanges"}
	}

	return nil
}

func (c *CryptoExchangeClient) GetExchange() types.Exchange {
	return c.exchange
}

func (c *CryptoExchangeClient) SendRequest(method, endpoint string, params map[string]interface{}, signed bool) (map[string]interface{}, error) {
	exchange := c.exchange

	req, err := exchange.PrepareRequest(method, endpoint, params, signed)
	if err != nil {
		return nil, &ExchangeError{Exchange: c.exchange.Name(), Message: err.Error()}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, &ExchangeError{Exchange: c.exchange.Name(), Message: err.Error()}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &ExchangeError{Exchange: c.exchange.Name(), Message: err.Error()}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, &ExchangeError{Exchange: c.exchange.Name(), Message: fmt.Sprintf("failed to parse response: %s", err.Error())}
	}

	return result, nil
}
