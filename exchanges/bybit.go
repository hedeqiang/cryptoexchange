package exchanges

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hedeqiang/cryptoexchange/types"
)

type Bybit struct {
	config types.ExchangeConfig
}

func NewBybit(config types.ExchangeConfig) *Bybit {
	return &Bybit{config: config}
}

func (b *Bybit) Name() types.ExchangeName {
	return types.Bybit
}

func (b *Bybit) GetDefaultBaseURL() string {
	return "https://api.bybit.com"
}

func (b *Bybit) PrepareRequest(method, endpoint string, params map[string]interface{}, signed bool) (*http.Request, error) {
	baseURL := b.config.BaseURL
	if baseURL == "" {
		baseURL = b.GetDefaultBaseURL()
	}

	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	if params == nil {
		params = make(map[string]interface{})
	}

	if signed {
		timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		params["api_key"] = b.config.APIKey
		params["timestamp"] = timestamp

		signature := b.generateSignature(params)
		params["sign"] = signature
	}

	queryString := b.buildQueryString(params)
	if method == "GET" || method == "DELETE" {
		u.RawQuery = queryString
	}

	var body []byte
	if method == "POST" || method == "PUT" {
		body = []byte(queryString)
	}

	req, err := http.NewRequest(method, u.String(), strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	if method == "POST" || method == "PUT" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return req, nil
}

func (b *Bybit) generateSignature(params map[string]interface{}) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signaturePayload strings.Builder
	for _, k := range keys {
		signaturePayload.WriteString(k)
		signaturePayload.WriteString("=")
		signaturePayload.WriteString(fmt.Sprint(params[k]))
		signaturePayload.WriteString("&")
	}
	payload := strings.TrimSuffix(signaturePayload.String(), "&")

	h := hmac.New(sha256.New, []byte(b.config.APISecret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

func (b *Bybit) buildQueryString(params map[string]interface{}) string {
	values := url.Values{}
	for k, v := range params {
		values.Add(k, fmt.Sprint(v))
	}
	return values.Encode()
}
