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
	if err != nil {
		return
	}

	params := map[string]interface{}{
		"symbol": "BTCUSDT",
	}

	response, err := client.SendRequest("GET", "/api/v3/ticker/price", params, false)

	assert.NoError(t, err)
	assert.NotNil(t, response)

}
