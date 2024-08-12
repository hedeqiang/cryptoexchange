package exchanges

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hedeqiang/cryptoexchange/types"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Gate struct {
	config types.ExchangeConfig
}

func NewGate(config types.ExchangeConfig) *Gate {
	return &Gate{config: config}
}

func (g *Gate) Name() types.ExchangeName {
	return types.Gate
}

func (g *Gate) GetDefaultBaseURL() string {
	return "https://api.gateio.ws"
}

func (g *Gate) PrepareRequest(method, endpoint string, params map[string]interface{}, signed bool) (*http.Request, error) {
	baseURL := g.config.BaseURL
	if baseURL == "" {
		baseURL = g.GetDefaultBaseURL()
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
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		payloadToSign := method + "\n" + endpoint + "\n" + u.RawQuery + "\n"
		if method != "GET" {
			hash := sha512.New()
			hash.Write(body)
			payloadToSign += hex.EncodeToString(hash.Sum(nil)) + "\n"
		}
		payloadToSign += timestamp

		mac := hmac.New(sha512.New, []byte(g.config.APISecret))
		mac.Write([]byte(payloadToSign))
		signature := hex.EncodeToString(mac.Sum(nil))

		req.Header.Set("KEY", g.config.APIKey)
		req.Header.Set("Timestamp", timestamp)
		req.Header.Set("SIGN", signature)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
