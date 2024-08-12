package cryptoexchange

import (
	"fmt"
	"github.com/hedeqiang/cryptoexchange/types"
)

type ExchangeError struct {
	Exchange types.ExchangeName
	Message  string
}

func (e *ExchangeError) Error() string {
	return fmt.Sprintf("Exchange %s error: %s", e.Exchange, e.Message)
}

type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API request failed with status %d: %s", e.StatusCode, e.Body)
}
