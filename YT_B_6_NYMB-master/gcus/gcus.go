package gcus

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/transaction"

	exchange "github.com/preichenberger/go-coinbase-exchange"
	log "github.com/sirupsen/logrus"
)

const (
	//the SQL table currencies are added to
	currencyTableName = "Currencies"
	//used to find the conversion between USD and a cryptocurrency
	productStringTrailerUSD = "-USD"
)

var (
	//tcp connection
	tConn net.Conn
	//list of the "real" currencies we're not interested in tracking
	restrictedCurrencies = [...]string{"USD", "GBP", "EUR"}
)

// Start the GDAX Currency Update Service. Is meant to be run as a subroutine. Updates the database currency table periodically
func Start() {
	// check GDAX API credentials and connectivity
	checkGDAXapiCredentials()
	connectToTransactionQueue()

	client := getClient()
	currencies := createCurrencyList(client)

	// TODO: generate this currency list automatically
	markovCurrencies := []string{"BTC", "ETH", "LTC"}
	setupMarkov(markovCurrencies)

	for {
		pollCurrencyPrices(currencies, client)
		time.Sleep(time.Duration(1) * time.Second)
	}
}

// Determines whether a string refers to one of the "restricted" currencies
// (i.e. the real-world currencies we don't track)
func isRestrictedCurrency(shortName string) bool {
	for _, checkAgainst := range restrictedCurrencies {
		if shortName == checkAgainst {
			return true
		}
	}

	return false
}

// Tests GDAX API connectivity
func checkGDAXapiCredentials() {
	client := getClient()
	_, err := client.GetCurrencies()

	if err != nil {
		panic(err)
	}
}

// Gets a client with the system Secret, Key, and Passphrase
func getClient() *exchange.Client {
	secret := os.Getenv("COINBASE_SECRET")
	key := os.Getenv("COINBASE_KEY")
	passphrase := os.Getenv("COINBASE_PASSPHRASE")

	client := exchange.NewClient(secret, key, passphrase)
	return client
}

// Resets the database currency table using values from the GDAX API.
// Returns a client to the API and the list of currencies stored in the SQL table
func createCurrencyList(client *exchange.Client) map[string]float64 {
	var currencies map[string]float64
	currencies = make(map[string]float64)

	exchangeCurrencies, err := client.GetCurrencies()
	if err != nil {
		log.Warningf("gcus: %s", err.Error())
		return currencies
	}

	for _, currency := range exchangeCurrencies {
		currencies[currency.Id] = 0
	}

	return currencies
}

// Updates the "currencies" map with new price data
func pollCurrencyPrices(currencies map[string]float64, client *exchange.Client) {
	for currencyName, currentPrice := range currencies {
		tick, err := (*client).GetTicker(currencyName + productStringTrailerUSD)
		if err != nil {
			continue
		}

		if currentPrice != tick.Price {
			prediction := addPriceData(currencyName, tick.Price)
			currencies[currencyName] = tick.Price
			writeCurrency(currencyName, currencies[currencyName], prediction)
		}
	}
}

func connectToTransactionQueue() {
	err := fmt.Errorf("")
	for err != nil {
		time.Sleep(time.Second)
		tConn, err = net.Dial("tcp", transaction.GCUSServiceAddress)
		if err != nil {
			log.Warningf("gcus: %s", err.Error())
		}
	}
	log.Debug("gcus: accepted transaction service connection")
}

//Writes currency information to internal TCP connection
func writeCurrency(shortName string, price float64, pricePrediction string) {
	tempCurrency := &transaction.Currency{
		ShortName:       shortName,
		Price:           price,
		PricePrediction: pricePrediction,
	}

	msg, err := json.Marshal(tempCurrency)
	if err != nil {
		log.Warningf("gcus: %s", err.Error())
		return
	}

	_, err = tConn.Write(msg)
	if err != nil {
		log.Warningf("gcus: %s", err.Error())
	}
}
