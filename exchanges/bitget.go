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

type Bitget struct {
	config types.ExchangeConfig
}

func NewBitget(config types.ExchangeConfig) *Bitget {
	return &Bitget{config: config}
}

func (b *Bitget) Name() types.ExchangeName {
	return types.Bitget
}

func (b *Bitget) GetDefaultBaseURL() string {
	return "https://api.bitget.com"
}

func (b *Bitget) PrepareRequest(method, endpoint string, params map[string]interface{}, signed bool) (*http.Request, error) {
	baseURL := b.config.BaseURL
	if baseURL == "" {
		baseURL = b.GetDefaultBaseURL()
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
		timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
		message := timestamp + method + endpoint
		if method != "GET" {
			message += string(body)
		}

		mac := hmac.New(sha256.New, []byte(b.config.APISecret))
		mac.Write([]byte(message))
		signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		req.Header.Set("ACCESS-KEY", b.config.APIKey)
		req.Header.Set("ACCESS-SIGN", signature)
		req.Header.Set("ACCESS-TIMESTAMP", timestamp)
		req.Header.Set("ACCESS-PASSPHRASE", b.config.APIPassphrase)
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
