package exchanges

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hedeqiang/cryptoexchange/types"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Kucoin struct {
	config types.ExchangeConfig
}

func NewKucoin(config types.ExchangeConfig) *Kucoin {
	return &Kucoin{config: config}
}

func (k *Kucoin) Name() types.ExchangeName {
	return types.Kucoin
}

func (k *Kucoin) GetDefaultBaseURL() string {
	return "https://api.kucoin.com"
}

func (k *Kucoin) PrepareRequest(method, endpoint string, params map[string]interface{}, signed bool) (*http.Request, error) {
	baseURL := k.config.BaseURL
	if baseURL == "" {
		baseURL = k.GetDefaultBaseURL()
	}

	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	body := []byte{}
	if method == "GET" {
		q := u.Query()
		for key, value := range params {
			q.Set(key, fmt.Sprint(value))
		}
		u.RawQuery = q.Encode()
	} else {
		body, err = json.Marshal(params)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	if signed {
		timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
		message := timestamp + method + endpoint
		if method != "GET" {
			message += string(body)
		}

		mac := hmac.New(sha256.New, []byte(k.config.APISecret))
		mac.Write([]byte(message))
		signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		req.Header.Set("KC-API-KEY", k.config.APIKey)
		req.Header.Set("KC-API-SIGN", signature)
		req.Header.Set("KC-API-TIMESTAMP", timestamp)
		req.Header.Set("KC-API-PASSPHRASE", k.config.APIPassphrase)
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
