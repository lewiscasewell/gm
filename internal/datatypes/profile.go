package datatypes

import (
	"crypto-price/internal/util"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Profile struct {
	Name              string    `json:"name"`
	Token             string    `json:"token"`
	BaseCurrency      string    `json:"baseCurrency"`
	FileUrl           string    `json:"fileUrl"`
	LastCheckedTime   time.Time `json:"lastCheckedTime"`
	LastCheckedAmount float64   `json:"lastCheckedAmount"`
}

var DefaultProfile = Profile{
	Name:              "Chief",
	Token:             "",
	BaseCurrency:      "GBP",
	FileUrl:           "",
	LastCheckedTime:   time.Now(),
	LastCheckedAmount: 0,
}

func (p Profile) String() string {
	return fmt.Sprintf("Name: %s\nToken: %s\nBase Currency: %s", p.Name, p.Token, p.BaseCurrency)
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

	fileUrl := p.FileUrl

	b, err := os.ReadFile(fileUrl)
	if err != nil {
		newProfile := []byte(`{
			"name": "Chief",
			"token": "",
			"baseCurrency": "GBP"
		}`)

		// check if gm directory exists
		if _, err := os.Stat(fmt.Sprintf("%s/gm", dirname)); os.IsNotExist(err) {
			if err = os.Mkdir(fmt.Sprintf("%s/gm", dirname), 0755); err != nil {
				panic(err)
			}
		}

		if err = os.WriteFile(fileUrl, newProfile, 0644); err != nil {
			panic(err)
		}

		b, err = os.ReadFile(fileUrl)
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

	profileURL := p.FileUrl

	if err = os.WriteFile(profileURL, b, 0644); err != nil {
		return err
	}

	return nil
}

func (p *Profile) SetNameAndSave(newName string) error {
	if newName == "" {
		return fmt.Errorf("name cannot be empty")
	}

	p.Name = newName

	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	profileURL := p.FileUrl

	if err = os.WriteFile(profileURL, b, 0644); err != nil {
		return err
	}

	return nil
}

func (p *Profile) SetLastCheckedAndSave(amount float64) error {
	p.LastCheckedTime = time.Now()
	p.LastCheckedAmount = amount

	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	profileURL := p.FileUrl

	if err = os.WriteFile(profileURL, b, 0644); err != nil {
		return err
	}

	return nil
}

func (p *Profile) SetCurrencyAndSave(newCurrency string) error {
	nc := strings.ToUpper(newCurrency)
	profileURL := p.FileUrl

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

func (p *Profile) Reset(newP Profile) error {
	newFileUrl := p.FileUrl

	newProfile, err := json.Marshal(newP)
	if err != nil {
		return err
	}

	if err := os.WriteFile(newFileUrl, newProfile, 0644); err != nil {
		return err
	}

	return nil
}

func (p *Profile) RestoreFileUrl() error {
	profileUrl, err := util.GetFileUrl("profile")
	if err != nil {
		return err
	}

	p.FileUrl = profileUrl

	newProfile, err := json.Marshal(p)
	if err != nil {
		return err
	}

	if err := os.WriteFile(profileUrl, newProfile, 0644); err != nil {
		return err
	}

	return nil
}
