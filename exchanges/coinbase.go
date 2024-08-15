package exchanges

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hedeqiang/cryptoexchange/types"
)

type Coinbase struct {
	config types.ExchangeConfig
}

func NewCoinbase(config types.ExchangeConfig) *Coinbase {
	return &Coinbase{config: config}
}

func (c *Coinbase) Name() types.ExchangeName {
	return types.Coinbase
}

func (c *Coinbase) GetDefaultBaseURL() string {
	return "https://api.exchange.coinbase.com"
}

func (c *Coinbase) PrepareRequest(method, endpoint string, params map[string]interface{}, signed bool) (*http.Request, error) {
	baseURL := c.config.BaseURL
	if baseURL == "" {
		baseURL = c.GetDefaultBaseURL()
	}

	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	var body []byte
	if method == "GET" {
		q := u.Query()
		for k, v := range params {
			q.Set(k, fmt.Sprint(v))
		}
		u.RawQuery = q.Encode()
	} else {
		body, err = json.Marshal(params)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if signed {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		message := timestamp + method + endpoint
		if method != "GET" {
			message += string(body)
		}

		signature := c.sign(message)

		req.Header.Set("CB-ACCESS-KEY", c.config.APIKey)
		req.Header.Set("CB-ACCESS-SIGN", signature)
		req.Header.Set("CB-ACCESS-TIMESTAMP", timestamp)
		req.Header.Set("CB-ACCESS-PASSPHRASE", c.config.APIPassphrase)
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Coinbase) sign(message string) string {
	h := hmac.New(sha256.New, []byte(c.config.APISecret))
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
