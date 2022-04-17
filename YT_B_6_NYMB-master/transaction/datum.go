package transaction

import "git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/models"

const (
	// WSTokenType defines the token type for WSDatums
	WSTokenType = "token"
	// WSCurrencyType defines the currencies type for WSDatums
	WSCurrencyType = "currency"
	// WSCurrenciesType defines the datum for initial currency prices
	WSCurrenciesType = "currencies"
	// WSTransactionType defines the transaction type for WSDatums
	WSTransactionType = "transaction"
	// WSMessageType defines the message type for WSDatums
	WSMessageType = "message"
)

// WSDatum defines the basic message struct for all tcp/websocket communication
type WSDatum struct {
	Type        string              `json:"type"`
	Token       *string             `json:"token,omitempty"`
	Currency    *WSCurrency         `json:"currency,omitempty"`
	Currencies  map[string]float64  `json:"currencies,omitempty"`
	Message     *WSMessage          `json:"message,omitempty"`
	Transaction *models.Transaction `json:"transaction,omitempty"`
}

// WSCurrency defines the currency price update datum
type WSCurrency struct {
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	PricePrediction string  `json:"priceprediction"`
}

// WSMessage defines the user to user messaging datum
type WSMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}
