package datatypes

import (
	"time"
)

type Rate struct {
	Currency string  `json:"currency"`
	Price    float64 `json:"ask"`
}

type ExchangeRateGBP struct {
	Meta struct {
		LastUpdatedAt time.Time `json:"last_updated_at"`
	} `json:"meta"`
	Data struct {
		Gbp struct {
			Code  string  `json:"code"`
			Value float64 `json:"value"`
		} `json:"GBP"`
	} `json:"data"`
}

type ExchangeRateEUR struct {
	Meta struct {
		LastUpdatedAt time.Time `json:"last_updated_at"`
	} `json:"meta"`
	Data struct {
		Eur struct {
			Code  string  `json:"code"`
			Value float64 `json:"value"`
		} `json:"EUR"`
	} `json:"data"`
}

type ExchangeRateUSD struct {
	Meta struct {
		LastUpdatedAt time.Time `json:"last_updated_at"`
	} `json:"meta"`
	Data struct {
		Usd struct {
			Code  string  `json:"code"`
			Value float64 `json:"value"`
		} `json:"USD"`
	} `json:"data"`
}
