# Crypto Exchange Client

A Go framework for accessing multiple cryptocurrency exchanges' APIs.

## Supported Exchanges
- Binance
- OKX
- Bitget
- Kucoin
- MEXC
- Gate.io
- Kraken
- Bybit

## Installation

```bash
go get github.com/hedeqiang/cryptoexchange
```

## Usage
First, import the necessary packages and set up your API configurations:
```go
package main

import (
	"log"
	"os"

	"github.com/hedeqiang/cryptoexchange"
	"github.com/hedeqiang/cryptoexchange/types"
)

func main() {
	// Create a client instance
	c := cryptoexchange.NewCryptoExchangeClient()

	// Set up your API configurations
	binanceConfig := types.ExchangeConfig{
		APIKey:    os.Getenv("BINANCE_API_KEY"),
		APISecret: os.Getenv("BINANCE_API_SECRET"),
	}

	// Add exchange to the client
	err := c.AddExchange(types.Binance, binanceConfig)
	if err != nil {
		log.Fatalf("Failed to add Binance exchange: %v", err)
	}

	// Now you can use the client to send requests
	// ...
}
```
### Explanation of the signed Parameter
- Public endpoints (such as market data) do not require authentication and can be accessed with signed=false.
- Private endpoints (such as account data or trading) require authentication and must be accessed with signed=true.

## Examples
### Binance
#### Public Endpoint (Market Data)

```go
// Get the latest price of BTC/USDT
params := map[string]interface{}{
    "symbol": "BTCUSDT",
}
response, err := c.SendRequest("GET", "/api/v3/ticker/price", params, false)
if err != nil {
    log.Fatalf("Failed to send request: %v", err)
}
fmt.Printf("Binance BTC/USDT price: %v\n", response)
```

#### Private Endpoint (Withdraw Funds)
```go
// Withdraw funds from your Binance account
params := map[string]interface{}{
    "coin": "USDT",
    "withdrawOrderId": "123456",
    "amount": 10,
    "network": "BSC",
    "address": "your_usdt_address",
}
response, err := c.SendRequest("POST", "/sapi/v1/capital/withdraw/apply", params, true)
if err != nil {
    log.Fatalf("Failed to send request: %v", err)
}
fmt.Printf("Withdrawal response: %v\n", response)
```

#### Get Account Balance
```go
response, err := c.SendRequest("GET", "/api/v3/account", nil, true)
if err != nil {
    log.Fatalf("Failed to get account balance: %v", err)
}
fmt.Printf("Binance account balance: %v\n", response)
````

### OKX
#### Public Endpoint (Market Data)
```go
// Get the ticker information for BTC/USDT
params := map[string]interface{}{
    "instId": "BTC-USDT",
}
response, err := c.SendRequest("GET", "/api/v5/market/ticker", params, false)
if err != nil {
    log.Fatalf("Failed to send request: %v", err)
}
fmt.Printf("OKX BTC/USDT ticker: %v\n", response)
```
#### Private Endpoint (Withdraw Funds)
```go
// Withdraw funds from your OKX account
params := map[string]interface{}{
    "ccy": "USDT",
    "amt": "10",
    "dest": "4",  // 4 means withdrawal to external address
    "toAddr": "your_usdt_address",
    "fee": "0.5",
}
response, err := c.SendRequest("POST", "/api/v5/asset/withdrawal", params, true)
if err != nil {
    log.Fatalf("Failed to send request: %v", err)
}
fmt.Printf("Withdrawal response: %v\n", response)
```

#### Get Account Balance
```go
params := map[string]interface{}{
    "ccy": "USDT",  // 可选，如果不提供则返回所有币种余额
}
response, err := c.SendRequest("GET", "/api/v5/account/balance", params, true)
if err != nil {
    log.Fatalf("Failed to get account balance: %v", err)
}
fmt.Printf("OKX account balance: %v\n", response)
```

### Bitget
#### Get Account Balance
```go
response, err := c.SendRequest("GET", "/api/v2/spot/account/assets", nil, true)
if err != nil {
    log.Fatalf("Failed to get account balance: %v", err)
}
fmt.Printf("Bitget account balance: %v\n", response)
```

### Kucoin
#### Get Account Balance
```go
response, err := c.SendRequest("GET", "/api/v1/accounts", nil, true)
if err != nil {
    log.Fatalf("Failed to get account balance: %v", err)
}
fmt.Printf("Kucoin account balance: %v\n", response)
```

### MEXC
#### Get Account Balance
```go
response, err := c.SendRequest("GET", "/api/v3/account", nil, true)
if err != nil {
    log.Fatalf("Failed to get account balance: %v", err)
}
fmt.Printf("MEXC account balance: %v\n", response)
```

### Gate.io
#### Get Account Balance
```go
response, err := c.SendRequest("GET", "/api/v4/spot/accounts", nil, true)
if err != nil {
    log.Fatalf("Failed to get account balance: %v", err)
}
fmt.Printf("Gate.io account balance: %v\n", response)
```

### Kraken
#### Get Account Balance
```go
response, err := c.SendRequest("POST", "/0/private/Balance", nil, true)
if err != nil {
    log.Fatalf("Failed to get account balance: %v", err)
}
fmt.Printf("Kraken account balance: %v\n", response)
```

### Bybit

#### Public Endpoint (Market Data)
```go
// Get the latest price of BTC/USDT
params := map[string]interface{}{
    "symbol": "BTCUSDT",
    "category": "spot",
}
response, err := c.SendRequest("GET", "/v5/market/tickers", params, false)
if err != nil {
    log.Fatalf("Failed to send request: %v", err)
}
fmt.Printf("Bybit BTC/USDT price: %v\n", response)
```

#### Private Endpoint (Get Account Balance)
```go
params := map[string]interface{}{
    "accountType": "UNIFIED",
}
response, err := c.SendRequest("GET", "/v5/account/wallet-balance", params, true)
if err != nil {
    log.Fatalf("Failed to get account balance: %v", err)
}
fmt.Printf("Bybit account balance: %v\n", response)
```
... (Similar examples for other exchanges)


## License
This project is licensed under the MIT License.
