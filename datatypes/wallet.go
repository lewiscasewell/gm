package datatypes

import (
	"crypto-price/util"
	"encoding/json"
	"os"
)

type Wallet map[string]float64

func (w *Wallet) GetFromJsonFile() error {
	walletURL, err := util.GetFileUrl("wallet")
	if err != nil {
		return err
	}

	b, err := os.ReadFile(walletURL)
	if err != nil {
		newEmptyWallet := []byte(`{}`)

		os.WriteFile(walletURL, newEmptyWallet, 0644)
		*w = make(Wallet)

		if err := json.Unmarshal(newEmptyWallet, &w); err != nil {
			return err
		}

		b, err = os.ReadFile(walletURL)
		if err != nil {
			return err
		}
	}

	if err := json.Unmarshal(b, &w); err != nil {
		return err
	}

	return nil
}

func (w *Wallet) Reset() error {
	walletURL, err := util.GetFileUrl("wallet")
	if err != nil {
		return err
	}
	newEmptyWallet := []byte(`{}`)
	os.WriteFile(walletURL, newEmptyWallet, 0644)

	return nil
}
