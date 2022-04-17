package gcus

import (
	"os"
	"testing"
)

// Makes sure restricted currencies aren't included in currency list
func TestIsRestrictedsCurrency(t *testing.T) {
	for _, c := range restrictedCurrencies {
		result := isRestrictedCurrency(c)
		if !result {
			t.Errorf("%s is not treated as a restricted currency", c)
		}
	}
}

func TestGetClient(t *testing.T) {
	secret := os.Getenv("COINBASE_SECRET")
	key := os.Getenv("COINBASE_KEY")
	passphrase := os.Getenv("COINBASE_PASSPHRASE")

	client := getClient()
	if client.Secret != secret {
		t.Errorf("client secret impropertly set")
	}
	if client.Key != key {
		t.Errorf("client key impropertly set")
	}
	if client.Passphrase != passphrase {
		t.Errorf("client passphrase impropertly set")
	}

}

// Makes sure BTC, ETH, and LTC are created by the "create currency list" function
func TestCreateCurrencyList(t *testing.T) {
	list := createCurrencyList(getClient())
	if len(list) < 1 {
		t.Errorf("No currencies added to list")
	}
	if _, ok := list["BTC"]; !ok {
		t.Errorf("Bitcoin not added to list")
	}
	if _, ok := list["ETH"]; !ok {
		t.Errorf("Etherium not added to list")
	}
	if _, ok := list["LTC"]; !ok {
		t.Errorf("Litecoin not added to list")
	}

}