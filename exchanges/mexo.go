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

type MEXC struct {
	config types.ExchangeConfig
}

func NewMEXC(config types.ExchangeConfig) *MEXC {
	return &MEXC{config: config}
}

func (m *MEXC) Name() types.ExchangeName {
	return types.MEXC
}

func (m *MEXC) GetDefaultBaseURL() string {
	return "https://api.mexc.com"
}

func (m *MEXC) PrepareRequest(method, endpoint string, params map[string]interface{}, signed bool) (*http.Request, error) {
	baseURL := m.config.BaseURL
	if baseURL == "" {
		baseURL = m.GetDefaultBaseURL()
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

		mac := hmac.New(sha256.New, []byte(m.config.APISecret))
		mac.Write([]byte(q.Encode()))
		signature := hex.EncodeToString(mac.Sum(nil))
		q.Set("signature", signature)
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-MEXC-APIKEY", m.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
