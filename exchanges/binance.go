package exchanges

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/hedeqiang/cryptoexchange/types"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Binance struct {
	config types.ExchangeConfig
}

func NewBinance(config types.ExchangeConfig) *Binance {
	return &Binance{config: config}
}

func (b *Binance) Name() types.ExchangeName {
	return types.Binance
}

func (b *Binance) GetDefaultBaseURL() string {
	return "https://api.binance.com"
}

func (b *Binance) PrepareRequest(method, endpoint string, params map[string]interface{}, signed bool) (*http.Request, error) {
	baseURL := b.config.BaseURL
	if baseURL == "" {
		baseURL = b.GetDefaultBaseURL()
	}

	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, fmt.Sprint(v))
	}

	if signed {
		timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
		q.Set("timestamp", timestamp)

		mac := hmac.New(sha256.New, []byte(b.config.APISecret))
		mac.Write([]byte(q.Encode()))
		signature := hex.EncodeToString(mac.Sum(nil))
		q.Set("signature", signature)
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-MBX-APIKEY", b.config.APIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}
