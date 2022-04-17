package models

import (
	"database/sql"
	"time"

	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/database"
	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/response"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	mysql "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

// TransactionStatus is the status of the transaction
type TransactionStatus string

const (
	listTransactions  = "SELECT * FROM Transactions"
	createTransaction = "INSERT INTO Transactions(FromID, ToID, CurrencyID, Amount, Created) VALUES(?, ?, ?, ?, ?)"
	readTransaction   = "SELECT * FROM Transactions WHERE ID=?"
	updateTransaction = "UPDATE Transactions SET Completed=? WHERE ID=?"
	deleteTransaction = "DELETE FROM Transactions WHERE ID=?"

	// Approve denotes a transaction's status as a request to approve
	Approve TransactionStatus = "approve"
	// Approved denotes a transaction's status as an approved request
	Approved TransactionStatus = "approved"
	// Deny denotes a transaction's status as a request to deny
	Deny TransactionStatus = "deny"
	// Denied denotes a transaction's status as a denied request
	Denied TransactionStatus = "denied"
	// New denotes a transaction's status as a new request
	New TransactionStatus = "new"
)

// Transaction defines information of a Transaction
type Transaction struct {
	ID         int               `form:"id" json:"id"`
	FromID     int               `form:"fromId" json:"fromId"`
	ToID       int               `form:"toId" json:"toId"`
	CurrencyID int               `form:"currencyId" json:"currencyId"`
	Amount     float64           `form:"amount" json:"amount"`
	Created    time.Time         `form:"created" json:"created"`
	Completed  *time.Time        `form:"completed" json:"completed"`
	Status     TransactionStatus `form:"status" json:"status"`
}

// RouteTransaction sets up the Transaction model's HTTP routes
func RouteTransaction(engine *gin.Engine, auth *jwt.GinJWTMiddleware) {
	authGroup := engine.Group("/transaction", auth.MiddlewareFunc())
	authGroup.GET("", ListTransaction)
	authGroup.POST("", CreateTransaction)
	authGroup.GET("/:id", ReadTransaction)
	authGroup.PUT("/:id", UpdateTransaction)
	authGroup.DELETE("/:id", DeleteTransaction)
}

// ListTransaction is the handler for listing Transactions
func ListTransaction(c *gin.Context) {
	query := GetQuery(listTransactions, c.Request.URL.Query())
	rows, err := database.Query(query)
	defer rows.Close()
	if err != nil {
		log.Errorf("transaction: %s", err.Error())
		response.ServerError(c)
		return
	}

	var transactions []Transaction
	for rows.Next() {
		transaction := Transaction{}
		if err = ScanTransaction(&transaction, rows); err != nil {
			log.Errorf("transaction: %s", err.Error())
			response.ServerError(c)
			return
		}
		transactions = append(transactions, transaction)
	}

	if err != nil {
		log.Warningf("transaction: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Data(c, transactions)
}

// CreateTransaction is the handler for creating Transactions
func CreateTransaction(c *gin.Context) {
	response.Ok(c)
}

// ReadTransaction is the handler for reading Transactions
func ReadTransaction(c *gin.Context) {
	response.Ok(c)
}

// UpdateTransaction is the handler for updating Transactions
func UpdateTransaction(c *gin.Context) {
	response.Ok(c)
}

// DeleteTransaction is the handler for Deleting Transactions
func DeleteTransaction(c *gin.Context) {
	response.Ok(c)
}

// ScanTransaction scans the given sql row into the given transaction
func ScanTransaction(t *Transaction, rows *sql.Rows) error {
	var completed mysql.NullTime
	if err := rows.Scan(&t.ID, &t.FromID, &t.ToID, &t.CurrencyID, &t.Amount,
		&t.Created, &completed, &t.Status); err != nil {
		return err
	}
	if completed.Valid {
		t.Completed = &completed.Time
	}
	return nil
}
