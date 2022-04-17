package models

import (
	"database/sql"
	"sort"
	"strconv"

	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/database"
	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/response"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	// CRUD for balances
	listBalances  = "SELECT * FROM Balances WHERE AccountID=?"
	createBalance = "INSERT INTO Balances(AccountID, CurrencyID, Amount) VALUES(?, ?, ?)"
	readBalance   = "SELECT * FROM Balances WHERE ID=?"
	updateBalance = "UPDATE Balances SET Amount=? WHERE ID=?"
	deleteBalance = "DELETE FROM Balances WHERE ID=?"

	checkIfBalanceExists = "SELECT COUNT(*) FROM Balances WHERE AccountID=? AND CurrencyID=?"
)

// Balance defines information of a Balance
type Balance struct {
	ID         int     `form:"id" json:"id"`
	AccountID  int     `form:"accountId" json:"accountId"`
	CurrencyID int     `form:"currencyId" json:"currencyId"`
	Amount     float64 `form:"amount" json:"amount"`
}

//TODO: weird amounts for currency (like strings) result in getting an amount of zero. Include error handling code for this

// RouteBalance sets up the Balance model's HTTP routes
func RouteBalance(engine *gin.Engine, auth *jwt.GinJWTMiddleware) {
	authGroup := engine.Group("/balance", auth.MiddlewareFunc())
	authGroup.GET("", ListBalance)
	authGroup.POST("", CreateBalance)
	authGroup.GET("/:id", ReadBalance)
	authGroup.PUT("/:id", UpdateBalance)
	authGroup.DELETE("/:id", DeleteBalance)
	adminGroup := engine.Group("/admin", auth.MiddlewareFunc())
	adminGroup.POST("/addFunds", AddFunds)
}

// ListBalance is the handler for listing Balances
func ListBalance(c *gin.Context) {
	cUser, err := getUserFromContext(c)
	if err != nil {
		log.Warningf("balance: %s", err.Error())
		response.NotFound(c)
		return
	}

	accountIDs, err := getAllAccountIDsForUser(cUser.ID, c)
	if err != nil {
		log.Errorf("balance: unable to get account IDs corresponding to user, error %s", err.Error())
		response.ServerError(c)
		return
	}

	// Get list of balances using Account IDs
	var totalBalances []Balance

	// loop through every account ID to add its balances to the list
	for i := 0; i < len(accountIDs); i++ {
		rows, err := database.Query(listBalances, accountIDs[i])
		if err != nil {
			log.Errorf("balance: %s", err.Error())
			response.ServerError(c)
			return
		}

		var balances []Balance
		for rows.Next() {
			balance := Balance{}
			if err = ScanBalance(&balance, rows); err != nil {
				log.Errorf("balance: %s", err.Error())
				response.ServerError(c)
				return
			}
			balances = append(balances, balance)
		}

		if err != nil {
			log.Warningf("balance: %s", err.Error())
			response.ServerError(c)
			return
		}

		rows.Close()
		totalBalances = append(totalBalances, balances...)
	}
	response.Data(c, totalBalances)
}

// CreateBalance is the handler for creating Balances
func CreateBalance(c *gin.Context) {

	balance := Balance{}
	err := c.Bind(&balance)

	if err != nil {
		log.Errorf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}

	// Checks to see if the balance has a valid currency
	result, err := currencyExists(balance.CurrencyID)
	if err != nil {
		log.Errorf("balance: error finding currency: %s", err.Error())
		response.ServerError(c)
		return

	}
	if !result {
		log.Warningf("balance: tried creating bad balance %s", balance)
		response.BadRequest(c)
		return
	}

	// Check if the user has permission to create the balance
	cUser, err := getUserFromContext(c)
	if err != nil {
		log.Warningf("balance: %s", err.Error())
		response.NotFound(c)
		return
	}

	accountIDs, err := getAllAccountIDsForUser(cUser.ID, c)
	if err != nil && cUser.RoleID != adminRoleID {
		log.Errorf("balance: unable to get account IDs corresponding to user, error %s", err.Error())
		response.ServerError(c)
		return
	}

	sort.Ints(accountIDs)
	if index := sort.SearchInts(accountIDs, balance.AccountID); (index >= len(accountIDs) || accountIDs[index] != balance.AccountID) && cUser.RoleID != adminRoleID {
		log.Warningf("balance: user %s does not have permission to create balance %s", cUser, balance)
		response.Unauthorized(c)
		return
	}
	// Check to see if a balance with the same AccountID and CurrencyID exists
	result, err = balanceExists(balance.AccountID, balance.CurrencyID)
	if err != nil {
		log.Errorf("balance: error finding balance: %s", err.Error())
		response.ServerError(c)
		return
	}
	if result {
		log.Warningf("balance: tried to create balance with duplicate account and currency %s", balance)
		response.BadRequest(c)
		return
	}

	// Create the balance
	idresult, err := database.Exec(createBalance, balance.AccountID, balance.CurrencyID, balance.Amount)
	if err != nil {
		log.Warningf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}

	intresult, err := idresult.LastInsertId()
	if err != nil {
		log.Errorf("balance: error finding balance ID: %s", err.Error())
		response.ServerError(c)
		return
	}

	balance.ID = int(intresult)

	response.Created(c, balance)
}

// ReadBalance is the handler for reading Balances
func ReadBalance(c *gin.Context) {
	BalanceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}

	// Read the balance
	rows, err := database.Query(readBalance, BalanceID)
	if err != nil {
		log.Errorf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}

	if ok := rows.Next(); !ok {
		log.Warningf("balance: could not find balance(%d)", BalanceID)
		response.NotFound(c)
		return
	}

	balance := Balance{}
	err = ScanBalance(&balance, rows)
	if err != nil {
		log.Errorf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}

	// Check if the user has permission to read the balance
	cUser, err := getUserFromContext(c)
	if err != nil {
		log.Warningf("balance: %s", err.Error())
		response.NotFound(c)
		return
	}

	accountIDs, err := getAllAccountIDsForUser(cUser.ID, c)
	if err != nil {
		log.Errorf("balance: unable to get account IDs corresponding to user, error %s", err.Error())
		response.ServerError(c)
		return
	}

	// Checks to see if the balance's accountID is an accountID bound to the user
	sort.Ints(accountIDs)
	index := sort.SearchInts(accountIDs, balance.AccountID)
	if (index >= len(accountIDs) || accountIDs[index] != balance.AccountID) && cUser.RoleID != adminRoleID {
		log.Warningf("balance: user %v does not have permission to read balance with ID %d", cUser)
		response.Unauthorized(c)
		return
	}

	response.Data(c, balance)
}

// UpdateBalance is the handler for updating Balances
func UpdateBalance(c *gin.Context) {
	// the balance to update
	balanceMessage := Balance{}
	err := c.Bind(&balanceMessage)

	if err != nil {
		log.Warningf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}

	// the ID of the balance to update
	BalanceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}

	// Read the balance from the database
	rows, err := database.Query(readBalance, BalanceID)
	if err != nil {
		log.Errorf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}

	if ok := rows.Next(); !ok {
		log.Warningf("balance: could not find balance(%d)", BalanceID)
		response.NotFound(c)
		return
	}

	balanceDatabase := Balance{}
	err = ScanBalance(&balanceDatabase, rows)
	if err != nil {
		log.Errorf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}

	// Check if this balance's accountID and currencyID are the same as the provided balance
	if balanceDatabase.AccountID != balanceMessage.AccountID ||
		balanceDatabase.CurrencyID != balanceMessage.CurrencyID {
		log.Warningf("balance: AccountID and/or CurrencyID do not match balance in database for %v", balanceMessage)
		response.BadRequest(c)
		return
	}

	// Check if the user has permission to read the balance
	cUser, err := getUserFromContext(c)
	if err != nil {
		log.Warningf("balance: %s", err.Error())
		response.NotFound(c)
		return
	}

	accountIDs, err := getAllAccountIDsForUser(cUser.ID, c)
	if err != nil {
		log.Errorf("balance: unable to get account IDs corresponding to user, error %s", err.Error())
		response.ServerError(c)
		return
	}

	// Checks to see if the balance's accountID is an accountID bound to the user
	sort.Ints(accountIDs)
	index := sort.SearchInts(accountIDs, balanceMessage.AccountID)
	if (index >= len(accountIDs) || accountIDs[index] != balanceMessage.AccountID) && cUser.RoleID != adminRoleID {
		log.Warningf("balance: user %v does not have permission to update balance with ID %d", cUser, BalanceID)
		response.Unauthorized(c)
		return
	}

	// Updates the balance
	_, err = database.Exec(updateBalance, balanceMessage.Amount, BalanceID)
	if err != nil {
		log.Warningf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Ok(c)
}

// DeleteBalance is the handler for Deleting Balances
func DeleteBalance(c *gin.Context) {
	// the ID of the balance to update
	BalanceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Errorf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}

	// Check if the user has permission to read the balance
	cUser, err := getUserFromContext(c)
	if err != nil {
		log.Warningf("balance: %s", err.Error())
		response.NotFound(c)
		return
	}

	accountIDs, err := getAllAccountIDsForUser(cUser.ID, c)
	if err != nil {
		log.Errorf("balance: unable to get account IDs corresponding to user, error %s", err.Error())
		response.ServerError(c)
		return
	}

	// Read the balance from the database
	rows, err := database.Query(readBalance, BalanceID)
	if err != nil {
		log.Errorf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}

	if ok := rows.Next(); !ok {
		log.Warningf("balance: could not find balance(%d)", BalanceID)
		response.NotFound(c)
		return
	}

	balanceDatabase := Balance{}
	err = ScanBalance(&balanceDatabase, rows)
	if err != nil {
		log.Errorf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}

	// Checks to see if the balance's accountID is an accountID bound to the user
	sort.Ints(accountIDs)
	index := sort.SearchInts(accountIDs, balanceDatabase.AccountID)
	if (index >= len(accountIDs) || accountIDs[index] != balanceDatabase.AccountID) && cUser.RoleID != adminRoleID {
		log.Warningf("balance: user %v does not have permission to delete balance with ID %d", cUser, BalanceID)
		response.Unauthorized(c)
		return
	}

	// Updates the balance
	_, err = database.Exec(deleteBalance, BalanceID)
	if err != nil {
		log.Warningf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Ok(c)
}

// ScanBalance scans the given sql row into the given currency
func ScanBalance(c *Balance, rows *sql.Rows) error {
	return rows.Scan(&c.ID, &c.AccountID, &c.CurrencyID, &c.Amount)
}

// Returns whether a balance exists. If an sql error happens, return false
func balanceExists(AccountID int, CurrencyID int) (bool, error) {
	// check if account ID and currency ID are present in table
	count := 0
	row, err := database.QueryRow(checkIfBalanceExists, AccountID, CurrencyID)
	row.Scan(&count)
	if err != nil {
		log.Warningf("balance: %s", err.Error())
		return false, err
	}

	if count > 0 {
		return true, err
	}

	return false, err
}

// AddFunds is the handler for admins to add funds to another user's account
func AddFunds(c *gin.Context) {

	// check if the current user is admin
	email, _ := c.Get("userID")
	user := &User{}
	rows, err := database.Query(readUserWithRoleByEmail, email)
	if err != nil {
		log.Warningf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}
	if rows.Next() {
		err = ScanUserWithEmailAndRole(user, rows)
		if err != nil {
			log.Warningf("balance: %s", err.Error())
			response.ServerError(c)
			return
		}
	}
	rows.Close()
	if user.RoleID != 1 {
		log.Warningf("balance: User(%s) tried adding funds", email)
		response.Unauthorized(c)
		return
	}

	// check for userId, currencyId, and amount
	type fund struct {
		UserID     int     `json:"userId"`
		CurrencyID int     `json:"currencyId"`
		Amount     float64 `json:"amount"`
	}
	req := &fund{}
	err = c.Bind(req)
	if err != nil {
		log.Warningf("balance: %s", err.Error())
		response.BadRequest(c)
		return
	}
	if req.UserID == 0 || req.CurrencyID == 0 || req.Amount == 0.0 {
		response.BadRequest(c)
		return
	}

	// get the user's account
	account := &Account{}
	rows, err = database.Query("SELECT * FROM Accounts WHERE UserID=?", req.UserID)
	if err != nil {
		log.Warningf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}
	if rows.Next() {
		err = ScanAccount(account, rows)
		if err != nil {
			log.Warningf("balance: %s", err.Error())
			response.ServerError(c)
			return
		}
	}
	rows.Close()

	// check if the balance exists
	balance := &Balance{}
	rows, err = database.Query("SELECT * FROM Balances WHERE AccountID=? AND CurrencyID=?",
		account.ID, req.CurrencyID)
	if err != nil {
		log.Warningf("balance: %s", err.Error())
		response.ServerError(c)
		return
	}
	// if the balance exists, increment
	if rows.Next() {
		err = rows.Scan(&balance.ID, &balance.AccountID, &balance.CurrencyID, &balance.Amount)
		if err != nil {
			log.Warningf("balance: %s", err.Error())
			response.ServerError(c)
			return
		}
		rows.Close()
		_, err = database.Exec(updateBalance, balance.Amount+req.Amount, balance.AccountID, balance.CurrencyID)
		if err != nil {
			log.Warningf("balance: %s", err.Error())
			response.ServerError(c)
			return
		}
	} else {
		// else create with amount
		_, err = database.Exec(createBalance, account.ID, req.CurrencyID, req.Amount)
		if err != nil {
			log.Warningf("balance: %s", err.Error())
			response.ServerError(c)
			return
		}
	}

	response.Ok(c)
}
