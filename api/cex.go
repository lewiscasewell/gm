package api

import (
	"crypto-price/datatypes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// currency exchange api for crypto currencies to USD
const apiUrl = "https://cex.io/api/ticker/%s/USD"

// currency exchange api for fiat currencies
const exchangeUrl = "https://api.currencyapi.com/v3/latest?apikey=cur_live_DInycfAzvFAqnXqqyAywtl3wXYcJkOVH8AGWZBiz&currencies=GBP&base_currency=USD"

func GetUsdGbpExchangeRate() float64 {
	res, err := http.Get(exchangeUrl)

	if err != nil {
		fmt.Println(err)
	}

	rate := &datatypes.ExchangeRate{}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}

		json.Unmarshal(bodyBytes, rate)
	}
	return float64(rate.Data.Gbp.Value)
}

func GetRate(currency string) (*datatypes.Rate, error) {
	if len(currency) != 3 {
		return nil, fmt.Errorf("currency code must be 3 characters")
	}

	ucc := strings.ToUpper(currency)

	rate := &datatypes.Rate{
		Currency: ucc,
	}

	res, err := http.Get(fmt.Sprintf(apiUrl, ucc))

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(bodyBytes, &rate)

		if err != nil {
			return nil, err
		}

	} else {
		return nil, fmt.Errorf("API returned status code %d", res.StatusCode)
	}

	return rate, err
}
