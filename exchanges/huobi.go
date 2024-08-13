package exchanges

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/hedeqiang/cryptoexchange/types"
)

type Huobi struct {
	config types.ExchangeConfig
}

func NewHuobi(config types.ExchangeConfig) *Huobi {
	return &Huobi{config: config}
}

func (h *Huobi) Name() types.ExchangeName {
	return types.Huobi
}

func (h *Huobi) GetDefaultBaseURL() string {
	return "https://api.huobi.pro"
}

func (h *Huobi) PrepareRequest(method, endpoint string, params map[string]interface{}, signed bool) (*http.Request, error) {
	baseURL := h.config.BaseURL
	if baseURL == "" {
		baseURL = h.GetDefaultBaseURL()
	}

	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	if params == nil {
		params = make(map[string]interface{})
	}

	if signed {
		timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")
		params["AccessKeyId"] = h.config.APIKey
		params["SignatureMethod"] = "HmacSHA256"
		params["SignatureVersion"] = "2"
		params["Timestamp"] = timestamp

		payload := h.buildPayload(method, u.Host, u.Path, params)
		signature := h.sign(payload)
		params["Signature"] = signature
	}

	queryString := h.buildQueryString(params)
	u.RawQuery = queryString

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (h *Huobi) buildPayload(method, host, path string, params map[string]interface{}) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var p []string
	for _, k := range keys {
		p = append(p, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(fmt.Sprint(params[k]))))
	}

	payload := strings.Join([]string{method, host, path, strings.Join(p, "&")}, "\n")
	return payload
}

func (h *Huobi) sign(payload string) string {
	hash := hmac.New(sha256.New, []byte(h.config.APISecret))
	hash.Write([]byte(payload))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func (h *Huobi) buildQueryString(params map[string]interface{}) string {
	values := url.Values{}
	for k, v := range params {
		values.Add(k, fmt.Sprint(v))
	}
	return values.Encode()
}
