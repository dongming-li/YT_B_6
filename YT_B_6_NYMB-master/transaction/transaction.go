package transaction

import (
	"sync"
)

const (
	// GCUSServiceAddress is the address which the gcus service writes currency data
	// and this transaction service reads the data.
	GCUSServiceAddress = ":8082"
	// APIServiceAddress is the address which the api service reads database
	// transaction update requests from this transaction service.
	APIServiceAddress = ":8083"

	btc = "BTC"
	eth = "ETH"
	ltc = "LTC"
	usd = "USD"
)

// Currency defines the realtime price information for a particular currency.
type Currency struct {
	ShortName       string  `json:"name"`
	Price           float64 `json:"price"`
	PricePrediction string  `json:"priceprediction"`
	mux             sync.Mutex
}

var (
	q           queue
	currencies  map[string]*Currency
	priceUpdate chan []byte
)

// Run begins the transaction queue and all related processes.
func Run() {

	sessions = make(map[int]*Session)
	q = *newQueue()
	currencies = make(map[string]*Currency)
	priceUpdate = make(chan []byte)

	currencies[btc] = &Currency{ShortName: btc, Price: 0.0}
	currencies[eth] = &Currency{ShortName: eth, Price: 0.0}
	currencies[ltc] = &Currency{ShortName: ltc, Price: 0.0}
	currencies[usd] = &Currency{ShortName: usd, Price: 1}

	go connectGCUS()

	sessionPricesUpdate()

}

func getCurrencies() map[string]float64 {
	var currs = make(map[string]float64)
	for k, c := range currencies {
		c.mux.Lock()
		defer c.mux.Unlock()
		currs[k] = c.Price
	}
	return currs
}

func sessionPricesUpdate() {
	for {
		select {
		case data, ok := <-priceUpdate:
			if ok {
				sessionsMux.RLock()
				for _, session := range sessions {
					session.sock.Send(string(data))
				}
				sessionsMux.RUnlock()
			}
		}
	}
}
