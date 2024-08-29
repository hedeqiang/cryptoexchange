package cryptoexchange

import (
	"testing"

	"github.com/hedeqiang/cryptoexchange/types"
	"github.com/stretchr/testify/assert"
)

func TestCryptoExchangeClient_AddExchange(t *testing.T) {
	client := NewCryptoExchangeClient()

	err := client.AddExchange(types.Binance, types.ExchangeConfig{
		APIKey:    "test_key",
		APISecret: "test_secret",
	})

	assert.NoError(t, err)
	assert.NotNil(t, client.exchange)
	assert.Equal(t, types.Binance, client.exchange.Name())
}

func TestCryptoExchangeClient_SendRequest(t *testing.T) {
	client := NewCryptoExchangeClient()
	err := client.AddExchange(types.Binance, types.ExchangeConfig{
		APIKey:    "test_key",
		APISecret: "test_secret",
	})
	assert.NoError(t, err)

	params := map[string]interface{}{
		"symbol": "BTCUSDT",
	}

	// 定义一个结构体来匹配预期的响应格式
	type TickerPrice struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}

	var response TickerPrice
	err = client.SendRequest("GET", "/api/v3/ticker/price", params, false, &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Symbol)
	assert.NotEmpty(t, response.Price)

	// 测试数组响应
	var arrayResponse []TickerPrice
	err = client.SendRequest("GET", "/api/v3/ticker/price", nil, false, &arrayResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, arrayResponse)
	if len(arrayResponse) > 0 {
		assert.NotEmpty(t, arrayResponse[0].Symbol)
		assert.NotEmpty(t, arrayResponse[0].Price)
	}
}
