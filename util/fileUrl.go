package util

import (
	"fmt"
	"os"
)

func GetFileUrl(file string) (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if file == "profile" {
		return fmt.Sprintf("%s/gm/profile.json", dirname), nil
	}
	if file == "wallet" {
		return fmt.Sprintf("%s/gm/wallet.json", dirname), nil
	}

	return "", fmt.Errorf("file %s not found", file)
}
