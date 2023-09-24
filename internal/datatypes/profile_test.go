package datatypes

import (
	"os"
	"testing"
)

const testUrlBase = "../../testdata/"

func TestProfile_GetFromJsonFile(t *testing.T) {
	testProfileUrl := testUrlBase + "profile.json"
	expectedProfile := Profile{
		Name:         "Test",
		Token:        "",
		BaseCurrency: "USD",
	}
	profile := Profile{
		FileUrl: testProfileUrl,
	}
	err := profile.GetFromJsonFile()
	if err != nil {
		t.Error(err)
	}

	if profile.Name != expectedProfile.Name {
		t.Errorf("expected name to be test but got %s", profile.Name)
	}

	if profile.BaseCurrency != expectedProfile.BaseCurrency {
		t.Errorf("expected base currency to be GBP but got %s", profile.BaseCurrency)
	}
}

func TestProfile_GetFromNonExistentJsonFile(t *testing.T) {
	testProfileUrl := testUrlBase + "non-existent.json"
	profile := Profile{
		FileUrl: testProfileUrl,
	}
	err := profile.GetFromJsonFile()
	if err != nil {
		t.Error(err)
	}

	// check file now exists
	if _, err := os.Stat(testProfileUrl); os.IsNotExist(err) {
		t.Errorf("expected file to exist but it doesn't")
	}

	if profile.Name != "Chief" {
		t.Errorf("expected name to be set to Chief but got %s", profile.Name)
	}

	// remove test file called empty-profile.json
	if err := os.Remove(testProfileUrl); err != nil {
		t.Error(err)
	}
}

func TestProfile_GetFiatSymbol(t *testing.T) {
	testProfileUrl := testUrlBase + "profile.json"
	profile := Profile{
		FileUrl: testProfileUrl,
	}
	err := profile.GetFromJsonFile()
	if err != nil {
		t.Error(err)
	}

	if profile.GetFiatSymbol() != "$" {
		t.Errorf("expected fiat symbol to be $ for USD but got %s", profile.GetFiatSymbol())
	}
}

func TestProfile_SetNameAndSave(t *testing.T) {
	testProfileUrl := testUrlBase + "profile.json"
	profile := Profile{
		FileUrl: testProfileUrl,
	}

	t.Log(profile)
	err := profile.GetFromJsonFile()
	if err != nil {
		t.Error(err)
	}

	err = profile.SetNameAndSave("Chief")
	if err != nil {
		t.Error(err)
	}

	if profile.Name != "Chief" {
		t.Errorf("expected name to be Test but got %s", profile.Name)
	}

	// reset name to Test
	err = profile.SetNameAndSave("Test")
	if err != nil {
		t.Error(err)
	}

	if profile.Name != "Test" {
		t.Errorf("expected name to be Test but got %s", profile.Name)
	}
}

func TestProfile_SetCurrencyAndSave(t *testing.T) {
	testProfileUrl := testUrlBase + "profile.json"
	profile := Profile{
		FileUrl: testProfileUrl,
	}
	err := profile.GetFromJsonFile()
	if err != nil {
		t.Error(err)
	}

	err = profile.SetCurrencyAndSave("GBP")
	if err != nil {
		t.Error(err)
	}

	if profile.BaseCurrency != "GBP" {
		t.Errorf("expected base currency to be GBP but got %s", profile.BaseCurrency)
	}

	// try invalid currency
	err = profile.SetCurrencyAndSave("NOT-VALID")
	if err == nil {
		t.Error("expected error but got nil")
	}

	// reset base currency to USD
	err = profile.SetCurrencyAndSave("USD")
	if err != nil {
		t.Error(err)
	}

	if profile.BaseCurrency != "USD" {
		t.Errorf("expected base currency to be USD but got %s", profile.BaseCurrency)
	}

}

func TestProfile_SetTokenAndSave(t *testing.T) {
	testProfileUrl := testUrlBase + "profile.json"
	profile := Profile{
		FileUrl: testProfileUrl,
	}
	err := profile.GetFromJsonFile()
	if err != nil {
		t.Error(err)
	}

	err = profile.SetTokenAndSave("test-token")
	if err != nil {
		t.Error(err)
	}

	if profile.Token != "test-token" {
		t.Errorf("expected token to be test-token but got %s", profile.Token)
	}

	// reset token
	err = profile.SetTokenAndSave("")
	if err != nil {
		t.Error(err)
	}

	if profile.Token != "" {
		t.Errorf("expected token to be empty but got %s", profile.Token)
	}
}

func TestProfile_Reset(t *testing.T) {
	testProfileUrl := testUrlBase + "profile.json"
	profile := Profile{
		FileUrl: testProfileUrl,
	}

	newProfile := Profile{
		Name:         "Test",
		Token:        "",
		BaseCurrency: "USD",
		FileUrl:      testProfileUrl,
	}
	err := profile.GetFromJsonFile()
	if err != nil {
		t.Error(err)
	}

	err = profile.Reset(newProfile)
	if err != nil {
		t.Error(err)
	}
	t.Log(profile)
	if profile.Name != newProfile.Name {
		t.Errorf("expected name to be Chief but got %s", profile.Name)
	}

	if profile.Token != "" {
		t.Errorf("expected token to be empty but got %s", profile.Token)
	}

	if profile.BaseCurrency != newProfile.BaseCurrency {
		t.Errorf("expected base currency to be GBP but got %s", profile.BaseCurrency)
	}
}
