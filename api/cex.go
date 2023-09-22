package api

import (
	"crypto-price/datatypes"
	"crypto-price/util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// currency exchange api for crypto currencies to USD
const apiUrl = "https://cex.io/api/ticker/%s/USD"

// currency exchange api for fiat currencies
var exchangeUrl string

func GetUsdGbpExchangeRate() (float64, error) {
	profileURL, err := util.GetFileUrl("profile")
	if err != nil {
		return 0, fmt.Errorf("error getting profile.json so cannot get API token")
	}
	b, err := os.ReadFile(profileURL)
	if err != nil {
		return 0, fmt.Errorf("error reading profile.json so cannot get API token")
	}

	var profile datatypes.Profile
	if err := json.Unmarshal(b, &profile); err != nil {
		return 0, fmt.Errorf("error unmarshalling profile.json so cannot get API token")
	}

	if profile.Token == "" {
		return 0, fmt.Errorf("please set your currency api token")
	}

	if profile.BaseCurrency != "GBP" && profile.BaseCurrency != "EUR" && profile.BaseCurrency != "USD" {
		return 0, fmt.Errorf("base currency must be GBP, EUR or USD, not %s", profile.BaseCurrency)
	}

	exchangeUrl = fmt.Sprintf("https://api.currencyapi.com/v3/latest?base_currency=USD&apikey=%s&currencies=%s", profile.Token, profile.BaseCurrency)

	res, err := http.Get(exchangeUrl)

	if err != nil {
		fmt.Println(err)
	}

	rateGBP := &datatypes.ExchangeRateGBP{}
	rateEUR := &datatypes.ExchangeRateEUR{}
	rateUSD := &datatypes.ExchangeRateUSD{}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}

		if profile.BaseCurrency == "GBP" {
			if err = json.Unmarshal(bodyBytes, &rateGBP); err != nil {
				fmt.Println(err)
			}
			return float64(rateGBP.Data.Gbp.Value), nil
		} else if profile.BaseCurrency == "EUR" {
			if err = json.Unmarshal(bodyBytes, &rateEUR); err != nil {
				fmt.Println(err)
			}
			return float64(rateEUR.Data.Eur.Value), nil
		} else {
			if err = json.Unmarshal(bodyBytes, &rateUSD); err != nil {
				fmt.Println(err)
			}

			return float64(rateUSD.Data.Usd.Value), nil
		}
	}
	return 0, fmt.Errorf("API returned status code %d", res.StatusCode)
}

func GetRate(currency string) (*datatypes.Rate, error) {
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
