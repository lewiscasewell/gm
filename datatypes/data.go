package datatypes

import "time"

type Rate struct {
	Currency string  `json:"currency"`
	Price    float64 `json:"ask"`
}

type ExchangeRate struct {
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

type Profile struct {
	Name         string `json:"name"`
	Token        string `json:"token"`
	BaseCurrency string `json:"baseCurrency"`
}

type Wallet map[string]float64
