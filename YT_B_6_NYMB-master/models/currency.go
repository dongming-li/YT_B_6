package models

import (
	"database/sql"

	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/database"
	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/response"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	listCurrencies = "SELECT * FROM Currencies"
	createCurrency = "INSERT INTO Currencies(Name, UnitPrice) VALUES(?, ?)"
	readCurrency   = "SELECT * FROM Currencies WHERE ID=?"
	updateCurrency = "UPDATE Currencies SET UnitPrice=? WHERE ID=?"
	deleteCurrency = "DELETE FROM Currencies WHERE ID=?"

	checkIfCurrencyExists = "SELECT COUNT(*) FROM Currencies WHERE ID=?"
)

// Currency defines information of a Currency
type Currency struct {
	ID        int    `form:"id" json:"id"`
	Name      string `form:"name" json:"name"`
	ShortName string `form:"shortname" json:"shortname"`
	UnitPrice int    `form:"unitPrice" json:"unitPrice"`
}

// RouteCurrency sets up the Currency model's HTTP routes
func RouteCurrency(engine *gin.Engine, auth *jwt.GinJWTMiddleware) {
	authGroup := engine.Group("/currency", auth.MiddlewareFunc())
	authGroup.GET("", ListCurrency)
	authGroup.POST("", CreateCurrency)
	authGroup.GET("/:id", ReadCurrency)
	authGroup.PUT("/:id", UpdateCurrency)
	authGroup.DELETE("/:id", DeleteCurrency)
}

// ListCurrency is the handler for listing Currencies
func ListCurrency(c *gin.Context) {
	query := GetQuery(listCurrencies, c.Request.URL.Query())
	rows, err := database.Query(query)
	defer rows.Close()
	if err != nil {
		log.Errorf("currency: %s", err.Error())
		response.ServerError(c)
		return
	}

	var currencies []Currency
	for rows.Next() {
		currency := Currency{}
		if err = ScanCurrency(&currency, rows); err != nil {
			log.Errorf("currency: %s", err.Error())
			response.ServerError(c)
			return
		}
		currencies = append(currencies, currency)
	}

	if err != nil {
		log.Warningf("currency: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Data(c, currencies)
}

// CreateCurrency is the handler for creating Currencies
//
// This will probably never be used because the only three
// rows in the database should be for Bitcoin, Litecoin, and Ethurium
// which should exist on application start up.
func CreateCurrency(c *gin.Context) {
	currency := Currency{}
	err := c.Bind(&currency)
	if err != nil {
		log.Warningf("currency: %s", err.Error())
		response.ServerError(c)
		return
	}

	if currency.Name == "" {
		log.Warningf("currency: tried creating bad currency %s", currency)
		response.BadRequest(c)
		return
	}

	_, err = database.Exec(createCurrency, currency.Name, currency.UnitPrice)
	if err != nil {
		log.Warningf("currency: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Created(c, currency)
}

// ReadCurrency is the handler for reading Currencies
func ReadCurrency(c *gin.Context) {
	id := c.Param("id")

	rows, err := database.Query(readCurrency, id)
	defer rows.Close()
	if err != nil {
		log.Warningf("currency: %s", err.Error())
		response.ServerError(c)
		return
	}

	currency := &Currency{}
	if ok := rows.Next(); !ok {
		log.Warningf("currency: could not read currency(%v)", id)
		response.NotFound(c)
		return
	}
	if err = ScanCurrency(currency, rows); err != nil {
		log.Warningf("currency: couldn't scan: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Data(c, currency)
}

// UpdateCurrency is the handler for updating Currencies
func UpdateCurrency(c *gin.Context) {
	id := c.Param("id")
	rows, err := database.Query(readCurrency, id)
	defer rows.Close()
	if err != nil {
		log.Warningf("currency: %s", err.Error())
		response.ServerError(c)
		return
	}

	currency := &Currency{}
	if ok := rows.Next(); !ok {
		log.Warningf("currency: could not update currency(%v)", id)
		response.NotFound(c)
		return
	}
	if err = ScanCurrency(currency, rows); err != nil {
		log.Warningf("currency: %s", err.Error())
		response.ServerError(c)
		return
	}

	err = c.Bind(currency)
	if err != nil {
		log.Warningf("currency: %s", err.Error())
		response.ServerError(c)
		return
	}

	_, err = database.Exec(updateCurrency, currency.UnitPrice)
	if err != nil {
		log.Warningf("currency: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Ok(c)
}

// DeleteCurrency is the handler for Deleting Currencies
func DeleteCurrency(c *gin.Context) {
	id := c.Param("id")
	rows, err := database.Query(readCurrency, id)
	defer rows.Close()
	if err != nil {
		log.Warningf("currency: %s", err.Error())
		response.ServerError(c)
		return
	}

	currency := &Currency{}
	if ok := rows.Next(); !ok {
		log.Warningf("currency: could not delete currency(%v)", id)
		response.NotFound(c)
		return
	}
	if err = ScanCurrency(currency, rows); err != nil {
		log.Warningf("currency: %s", err.Error())
		response.ServerError(c)
		return
	}

	_, err = database.Exec(deleteCurrency, id)
	if err != nil {
		log.Warningf("currency: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Ok(c)
}

// ScanCurrency scans the given sql row into the given currency
func ScanCurrency(c *Currency, rows *sql.Rows) error {
	return rows.Scan(&c.ID, &c.Name, &c.ShortName, &c.UnitPrice)
}

// Checks to see if a currency exists. If an sql error happens, return false
func currencyExists(CurrencyID int) (bool, error) {
	currencyCount := 0
	row, err := database.QueryRow(checkIfCurrencyExists, CurrencyID)
	row.Scan(&currencyCount)
	if err != nil {
		log.Warningf("currency: %s", err.Error())
		return false, err
	}

	if currencyCount > 0 {
		return true, err
	}

	return false, err
}
