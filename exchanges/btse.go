package exchanges

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hedeqiang/cryptoexchange/types"
)

type BTSE struct {
	config types.ExchangeConfig
}

func NewBTSE(config types.ExchangeConfig) *BTSE {
	return &BTSE{config: config}
}

func (b *BTSE) Name() types.ExchangeName {
	return types.BTSE
}

func (b *BTSE) GetDefaultBaseURL() string {
	return "https://api.btse.com/spot"
}

func (b *BTSE) PrepareRequest(method, endpoint string, params map[string]interface{}, signed bool) (*http.Request, error) {
	baseURL := b.config.BaseURL
	if baseURL == "" {
		baseURL = b.GetDefaultBaseURL()
	}

	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	headers := http.Header{
		"Content-Type": []string{"application/json"},
	}

	var bodyStr string
	if method == "GET" && params != nil {
		q := u.Query()
		for k, v := range params {
			q.Set(k, fmt.Sprint(v))
		}
		u.RawQuery = q.Encode()
	} else if method == "POST" && params != nil {
		jsonBody, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		bodyStr = string(jsonBody)
	}

	if signed {
		requestNonce := fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
		concatenatedStr := endpoint + requestNonce + bodyStr
		signature := b.sign(concatenatedStr)

		headers.Set("request-api", b.config.APIKey)
		headers.Set("request-nonce", requestNonce)
		headers.Set("request-sign", signature)
	}

	req, err := http.NewRequest(method, u.String(), strings.NewReader(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header = headers

	return req, nil
}

func (b *BTSE) sign(message string) string {
	h := hmac.New(sha512.New384, []byte(b.config.APISecret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
