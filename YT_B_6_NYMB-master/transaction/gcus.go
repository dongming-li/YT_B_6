package transaction

import (
	"encoding/json"
	"net"

	log "github.com/sirupsen/logrus"
)

func connectGCUS() {

	listener, err := net.Listen("tcp", GCUSServiceAddress)
	if err != nil {
		log.Warningf("transaction: %s", err.Error())
		return
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Warningf("transaction: %s", err.Error())
		return
	}
	log.Debug("transaction: accepted gcus service connection")

	gcusReader(conn)

}

func gcusReader(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			log.Warningf("transaction: %s", err.Error())
			return
		}
		go updateCurrency(buffer[0:length])
	}
}

func updateCurrency(msg []byte) {
	c := &Currency{}
	if err := json.Unmarshal(msg, c); err != nil {
		log.Warningf("transaction: %s", err.Error())
		return
	}

	var currency *Currency
	switch c.ShortName {
	case btc:
		currency = currencies[btc]
		break
	case eth:
		currency = currencies[eth]
		break
	case ltc:
		currency = currencies[ltc]
		break
	}

	currency.mux.Lock()
	defer currency.mux.Unlock()
	currency.Price = c.Price
	currency.PricePrediction = c.PricePrediction
	data := &WSDatum{
		Type: WSCurrencyType,
		Currency: &WSCurrency{
			Name:            currency.ShortName,
			Price:           currency.Price,
			PricePrediction: currency.PricePrediction,
		},
	}
	msg, err := json.Marshal(data)
	if err != nil {
		log.Warningf("transaction: %s", err.Error())
	}
	priceUpdate <- msg
}
