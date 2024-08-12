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
	"time"
)

type OKX struct {
	config types.ExchangeConfig
}

func NewOKX(config types.ExchangeConfig) *OKX {
	return &OKX{config: config}
}

func (o *OKX) Name() types.ExchangeName {
	return types.OKX
}

func (o *OKX) GetDefaultBaseURL() string {
	return "https://www.okx.com"
}

func (o *OKX) PrepareRequest(method, endpoint string, params map[string]interface{}, signed bool) (*http.Request, error) {
	baseURL := o.config.BaseURL
	if baseURL == "" {
		baseURL = o.GetDefaultBaseURL()
	}

	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	body := []byte{}
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

	req, err := http.NewRequest(method, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	if signed {
		timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
		message := timestamp + method + endpoint
		if method != "GET" {
			message += string(body)
		}

		mac := hmac.New(sha256.New, []byte(o.config.APISecret))
		mac.Write([]byte(message))
		signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		req.Header.Set("OK-ACCESS-KEY", o.config.APIKey)
		req.Header.Set("OK-ACCESS-SIGN", signature)
		req.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Set("OK-ACCESS-PASSPHRASE", o.config.APIPassphrase)
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
