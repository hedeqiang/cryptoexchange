package exchanges

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/hedeqiang/cryptoexchange/types"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Kraken struct {
	config types.ExchangeConfig
}

func NewKraken(config types.ExchangeConfig) *Kraken {
	return &Kraken{config: config}
}

func (k *Kraken) Name() types.ExchangeName {
	return types.Kraken
}

func (k *Kraken) GetDefaultBaseURL() string {
	return "https://api.kraken.com"
}

func (k *Kraken) PrepareRequest(method, endpoint string, params map[string]interface{}, signed bool) (*http.Request, error) {
	baseURL := k.config.BaseURL
	if baseURL == "" {
		baseURL = k.GetDefaultBaseURL()
	}

	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	if params == nil {
		params = make(map[string]interface{})
	}

	if signed {
		nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
		params["nonce"] = nonce

		postData := url.Values{}
		for key, value := range params {
			postData.Set(key, fmt.Sprint(value))
		}

		signature, err := k.getKrakenSignature(endpoint, postData.Encode(), nonce)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(method, u.String(), strings.NewReader(postData.Encode()))
		if err != nil {
			return nil, err
		}

		req.Header.Set("API-Key", k.config.APIKey)
		req.Header.Set("API-Sign", signature)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		return req, nil
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, fmt.Sprint(v))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (k *Kraken) getKrakenSignature(endpoint, postData, nonce string) (string, error) {
	sha256Sum := sha256.Sum256([]byte(nonce + postData))
	pathBytes := []byte(endpoint)
	sha512Input := append(pathBytes, sha256Sum[:]...)

	decodedSecret, err := base64.StdEncoding.DecodeString(k.config.APISecret)
	if err != nil {
		return "", fmt.Errorf("failed to decode API secret: %v", err)
	}

	mac := hmac.New(sha512.New, decodedSecret)
	mac.Write(sha512Input)
	macSum := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(macSum), nil
}
