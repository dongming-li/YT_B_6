package models

import (
	"database/sql"
	"fmt"

	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/database"
	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/response"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	listAccounts  = "SELECT * FROM Accounts"
	createAccount = "INSERT INTO Accounts(%s) VALUES(?)"
	readAccount   = "SELECT * FROM Accounts WHERE ID=?"
	updateAccount = "UPDATE Accounts SET %s=? WHERE ID=?"
	deleteAccount = "DELETE FROM Accounts WHERE ID=?"

	getUserAccount  = "SELECT * FROM Accounts WHERE UserID=?"
	getVaultAccount = "SELECT * FROM Accounts WHERE VaultID=?"
)

// Account defines information of a Account
type Account struct {
	ID      int `form:"id" json:"id"`
	UserID  int `form:"userId" json:"userId"`
	VaultID int `form:"vaultId" json:"vaultId"`
}

// RouteAccount sets up the Account model's HTTP routes
func RouteAccount(engine *gin.Engine, auth *jwt.GinJWTMiddleware) {
	authGroup := engine.Group("/account", auth.MiddlewareFunc())
	authGroup.GET("", ListAccount)
	authGroup.POST("", CreateAccount)
	authGroup.GET("/:id", ReadAccount)
	authGroup.PUT("/:id", UpdateAccount)
	authGroup.DELETE("/:id", DeleteAccount)
}

// ListAccount is the handler for listing Accounts
func ListAccount(c *gin.Context) {
	query := GetQuery(listAccounts, c.Request.URL.Query())
	rows, err := database.Query(query)
	if err != nil {
		log.Errorf("account: %s", err.Error())
		response.ServerError(c)
		return
	}

	var accounts []Account
	for rows.Next() {
		account := Account{}
		if err = ScanAccount(&account, rows); err != nil {
			log.Errorf("account: %s", err.Error())
			response.ServerError(c)
			return
		}
		accounts = append(accounts, account)
	}

	if err != nil {
		log.Warningf("account: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Data(c, accounts)
}

// CreateAccount is the handler for creating Accounts
func CreateAccount(c *gin.Context) {
	account := Account{}
	err := c.Bind(&account)
	if err != nil {
		log.Warningf("account: %s", err.Error())
		response.ServerError(c)
		return
	}

	if account.UserID != 0 {
		query := fmt.Sprintf(createAccount, "UserID")
		_, err = database.Exec(query, account.UserID)
	} else if account.VaultID != 0 {
		query := fmt.Sprintf(createAccount, "VaultID")
		_, err = database.Exec(query, account.VaultID)
	} else {
		log.Warningf("account: tried creating bad account %s", account)
		response.BadRequest(c)
		return
	}

	if err != nil {
		log.Warningf("account: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Created(c, account)
}

// ReadAccount is the handler for reading Accounts
func ReadAccount(c *gin.Context) {
	id := c.Param("id")

	rows, err := database.Query(readAccount, id)
	if err != nil {
		log.Warningf("account: %s", err.Error())
		response.ServerError(c)
		return
	}

	account := &Account{}
	rows.Next()
	if err = ScanAccount(account, rows); err != nil {
		log.Warningf("account: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Data(c, account)
}

// UpdateAccount is the handler for updating Accounts
//
// This will probably never be used. Every account is
// associated with either a User or Vault. Updating
// an Account implies transfering an Account from
// one User or Vault to another. The other User or
// Vault probably already has an Account and no
// User or Vault should have more than one Account.
func UpdateAccount(c *gin.Context) {

	// retrieve account row by given url id
	id := c.Param("id")
	rows, err := database.Query(readAccount, id)
	if err != nil {
		log.Warningf("account: %s", err.Error())
		response.ServerError(c)
		return
	}

	// scan row into account
	account := &Account{}
	rows.Next()
	if err = ScanAccount(account, rows); err != nil {
		log.Warningf("account: %s", err.Error())
		response.ServerError(c)
		return
	}

	// check kind of Account (User or Vault)
	if account.UserID != 0 {

		// retrieve the User accociated with the Account
		user := &User{}
		rows, err = database.Query(readUser, account.UserID)
		rows.Next()
		if err = ScanUser(user, rows); err != nil {
			log.Warningf("account: %s", err.Error())
			response.ServerError(c)
			return
		}

		// get the client's email
		email, ok := c.Get("userID")
		if !ok {
			log.Fatal("account: couldn't get email from context")
			response.ServerError(c)
			return
		}

		// make sure the client's email matches the Account's User email
		if user.Email != email {
			log.Warningf("account: user(%s) tried to update %s", email, account)
			response.Unauthorized(c)
			return
		}

		// bind the context's form data to account
		err = c.Bind(account)
		if err != nil {
			log.Warningf("account: %s", err.Error())
			response.ServerError(c)
			return
		}

		// update the account
		query := fmt.Sprintf(updateAccount, "UserID")
		_, err = database.Exec(query, account.UserID, account.ID)
		if err != nil {
			log.Warningf("account: %s", err.Error())
			response.ServerError(c)
			return
		}

	} else if account.VaultID != 0 {

		// TODO: add vault lookup
		response.NotFound(c)
		return

	} else {
		log.Warningf("account: invalid %s", account)
		response.BadRequest(c)
		return
	}

	response.Ok(c)
}

// DeleteAccount is the handler for Deleting Accounts
//
// This would need to be called before a User or Vault is deleted
func DeleteAccount(c *gin.Context) {
	id := c.Param("id")
	rows, err := database.Query(readAccount, id)
	if err != nil {
		log.Warningf("account: %s", err.Error())
		response.ServerError(c)
		return
	}

	account := &Account{}
	rows.Next()
	if err = ScanAccount(account, rows); err != nil {
		log.Warningf("account: %s", err.Error())
		response.ServerError(c)
		return
	}

	if account.UserID != 0 {

		// retrieve the User accociated with the Account
		user := &User{}
		rows, err = database.Query(readUser, account.UserID)
		rows.Next()
		if err = ScanUser(user, rows); err != nil {
			log.Warningf("account: %s", err.Error())
			response.ServerError(c)
			return
		}

		// get the client's email
		email, ok := c.Get("userID")
		if !ok {
			log.Fatal("account: couldn't get email from context")
			response.ServerError(c)
			return
		}

		// make sure the client's email matches the Account's User email
		if user.Email != email {
			log.Warningf("account: user(%s) tried to delete %s", email, account)
			response.Unauthorized(c)
			return
		}

		// try to delete the account
		_, err = database.Exec(deleteAccount, id)
		if err != nil {
			log.Warningf("account: %s", err.Error())
			response.ServerError(c)
			return
		}

	} else if account.VaultID != 0 {

		// TODO: add vault lookup
		response.NotFound(c)
		return

	} else {
		log.Warningf("account: invalid %s", account)
		response.BadRequest(c)
		return
	}

	response.Ok(c)
}

// ScanAccount scans the given sql row into the given account
func ScanAccount(a *Account, rows *sql.Rows) error {
	var userID sql.NullInt64
	var vaultID sql.NullInt64
	if err := rows.Scan(&a.ID, &userID, &vaultID); err != nil {
		return err
	}
	if userID.Valid {
		a.UserID = int(userID.Int64)
	}
	if vaultID.Valid {
		a.VaultID = int(vaultID.Int64)
	}
	return nil
}

func (a Account) String() string {
	return fmt.Sprintf("Account(id=%d, userId=%d, vaultId=%d)", a.ID, a.UserID, a.VaultID)
}

// Finds the list of accounts tied to a user's personal account
func getAccountIDsForUser(UserID int, c *gin.Context) ([]int, error) {
	rows, err := database.Query("SELECT * FROM Accounts WHERE UserID=?", UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accountIDs []int

	for rows.Next() {
		account := Account{}
		if err = ScanAccount(&account, rows); err != nil {
			return accountIDs, err
		}
		accountIDs = append(accountIDs, account.ID)
	}

	return accountIDs, nil
}

// Finds the list of accounts tied to a vault
func getAccountIDsForVault(VaultID int, c *gin.Context) ([]int, error) {
	rows, err := database.Query("SELECT * FROM Accounts WHERE VaultID=?", VaultID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accountIDs []int

	for rows.Next() {
		account := Account{}
		if err = ScanAccount(&account, rows); err != nil {
			return accountIDs, err
		}
		accountIDs = append(accountIDs, account.ID)
	}

	return accountIDs, nil
}

// Finds the list of accounts tied to an admin account (i.e., all of them)
func getAccountIDsForAdmin(UserID int, c *gin.Context) ([]int, error) {
	rows, err := database.Query(listAccounts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accountIDs []int

	for rows.Next() {
		account := Account{}
		if err = ScanAccount(&account, rows); err != nil {
			return accountIDs, err
		}
		accountIDs = append(accountIDs, account.ID)
	}

	return accountIDs, nil
}

// Finds the list of accounts tied to a user's personal account and all vaults they own
func getAllAccountIDsForUser(UserID int, c *gin.Context) ([]int, error) {
	if UserID == adminRoleID {
		return getAccountIDsForAdmin(UserID, c)
	}

	accountIDs, err := getAccountIDsForUser(UserID, c)
	if err != nil {
		return accountIDs, err
	}

	vaultIDs, err := getVaultIDsForUser(UserID, c)
	if err != nil {
		return accountIDs, err
	}

	for _, vaultID := range vaultIDs {

		vaultAccountIDs, err := getAccountIDsForVault(vaultID, c)
		accountIDs = append(accountIDs, vaultAccountIDs...)
		if err != nil {

			return accountIDs, err
		}
	}

	return accountIDs, nil
}
