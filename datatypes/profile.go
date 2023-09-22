package datatypes

import (
	"crypto-price/util"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Profile struct {
	Name         string `json:"name"`
	Token        string `json:"token"`
	BaseCurrency string `json:"baseCurrency"`
}

func (p *Profile) GetFiatSymbol() string {
	switch p.BaseCurrency {
	case "GBP":
		return "£"
	case "EUR":
		return "€"
	case "USD":
		return "$"
	}

	return "£"
}

func (p *Profile) GetFromJsonFile() error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	profileURL, err := util.GetFileUrl("profile")
	if err != nil {
		return err
	}

	b, err := os.ReadFile(profileURL)
	if err != nil {
		newProfile := []byte(`{
			"name": "Chief",
			"token": "",
			"baseCurrency": "GBP"
		}`)

		if err = os.Mkdir(fmt.Sprintf("%s/gm", dirname), 0755); err != nil {
			panic(err)
		}

		if err = os.WriteFile(profileURL, newProfile, 0644); err != nil {
			panic(err)
		}

		b, err = os.ReadFile(profileURL)
		if err != nil {
			panic(err)
		}

	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	return nil
}

func (p *Profile) SetTokenAndSave(newToken string) error {
	p.Token = newToken

	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	profileURL, err := util.GetFileUrl("profile")
	if err != nil {
		return err
	}

	if err = os.WriteFile(profileURL, b, 0644); err != nil {
		return err
	}

	return nil
}

func (p *Profile) SetNameAndSave(newName string) error {
	p.Name = newName

	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	profileURL, err := util.GetFileUrl("profile")
	if err != nil {
		return err
	}

	if err = os.WriteFile(profileURL, b, 0644); err != nil {
		return err
	}

	return nil
}

func (p *Profile) SetCurrencyAndSave(newCurrency string) error {
	nc := strings.ToUpper(newCurrency)
	profileURL, err := util.GetFileUrl("profile")
	if err != nil {
		return err
	}

	if nc != "GBP" && nc != "USD" && nc != "EUR" {
		return fmt.Errorf("base currency must be GBP, USD or EUR, not %s", nc)
	}

	p.BaseCurrency = nc

	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	os.WriteFile(profileURL, b, 0644)
	return nil
}
