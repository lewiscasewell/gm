package api_test

import (
	"crypto-price/api"
	"testing"
)

func TestApiCallError(t *testing.T) {
	_, err := api.GetRate("")
	if err == nil {
		t.Error("Expected error but got nil")
	}
}
