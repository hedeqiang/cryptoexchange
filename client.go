package cryptoexchange

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/hedeqiang/cryptoexchange/exchanges"
	"github.com/hedeqiang/cryptoexchange/types"
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
	case types.Coinbase:
		c.exchange = exchanges.NewCoinbase(config)
	case types.BTSE:
		c.exchange = exchanges.NewBTSE(config)
	default:
		return &ExchangeError{Exchange: name, Message: "unsupported exchanges"}
	}

	return nil
}

func (c *CryptoExchangeClient) GetExchange() types.Exchange {
	return c.exchange
}

func (c *CryptoExchangeClient) SendRequest(method, endpoint string, params map[string]interface{}, signed bool, result interface{}) error {
	exchange := c.exchange

	req, err := exchange.PrepareRequest(method, endpoint, params, signed)
	if err != nil {
		return &ExchangeError{Exchange: c.exchange.Name(), Message: err.Error()}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return &ExchangeError{Exchange: c.exchange.Name(), Message: err.Error()}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &ExchangeError{Exchange: c.exchange.Name(), Message: err.Error()}
	}

	if resp.StatusCode != http.StatusOK {
		return &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	// 使用反射来确定结果类型并相应地解析
	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() != reflect.Ptr || resultValue.IsNil() {
		return fmt.Errorf("result must be a non-nil pointer")
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return &ExchangeError{Exchange: c.exchange.Name(), Message: fmt.Sprintf("failed to parse response: %s", err.Error())}
	}

	return nil
}
