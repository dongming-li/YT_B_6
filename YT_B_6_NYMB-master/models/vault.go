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
	listVaults = "SELECT v.ID, v.Name, u2.UserName as Owner, v.OwnerID as OwnerID, acc.ID as AccountID, p.UserID as UserID  "+ 
	"From users as u "+
	"Join Permissions as p "+
	"On u.id=p.UserID "+
	"Join Vaults as v "+
	"On v.id=p.VaultID "+
	"JOIN Accounts as acc "+ 
	"ON v.id=acc.VaultID "+
	"Join Vaults as v2 "+
	"On v2.id=p.VaultID "+
 	"Join Users as u2 "+
	"On u2.ID=v2.OwnerID "+
	"Where p.UserID=?"

	
		
	createVault = "INSERT INTO Vaults(Name, OwnerID) VALUES(?, ?)"
	
	readVault   = "SELECT * FROM Vaults WHERE ID=?"
	updateVault = "UPDATE Vaults SET Name=? WHERE ID=?"
	deleteVault = "DELETE FROM Vaults WHERE ID=?"
	
)

// Vault defines information of a Vault
type Vault struct {
	ID		int    `form:"id" json:"id"`
	Name	string `form:"name" json:"name"`
	Owner	string `json:"owner"`
	OwnerID	int    `form:"ownerID" json:"ownerID"`
	AccountID	int  `form:"account" json:"account"`
	UserID	int    `form:"userID" json:"userID"`
}



// RouteVault sets up the Vault model's HTTP routes
func RouteVault(engine *gin.Engine, auth *jwt.GinJWTMiddleware) {
	authGroup := engine.Group("/vault", auth.MiddlewareFunc())
	authGroup.GET("", ListVault)
	authGroup.POST("", CreateVault)
	authGroup.GET("/:id", ReadVault)
	authGroup.PUT("/:id", UpdateVault)
	authGroup.DELETE("/:id", DeleteVault)
	authGroup.GET("/:id/permission", getPermissionsForAVault)
	authGroup.GET("/:id/allpermissions", getAllPermissionsForAVault)
}

// ListVault is the handler for listing vaults
func ListVault(c *gin.Context) {
	cUser, err := getUserFromContext(c)
	if err != nil {
		log.Errorf("vault: %s", err.Error())
		response.ServerError(c)
		return
	}
	rows, err := database.Query(listVaults, cUser.ID)
	// query := GetQuery(listVaults, c.Request.URL.Query())
	// rows, err := database.Query(query)
	defer rows.Close()
	if err != nil {
		log.Errorf("vault: %s", err.Error())
		response.ServerError(c)
		return
	}

	var vaults []Vault
	for rows.Next() {
		vault := Vault{}
		if err = ScanVaultWithOwner(&vault, rows); err != nil {
			log.Errorf("vault: %s", err.Error())
			response.ServerError(c)
			return
		}
		vaults = append(vaults, vault)
	}

	if err != nil {
		log.Warningf("vault: %s", err.Error())
		response.ServerError(c)
		return
	}
	fmt.Printf("\nvaults: %+v\n",vaults )
	response.Data(c, vaults)
}

// CreateVault is the handler for creating vaults
func CreateVault(c *gin.Context) {
	vault := Vault{}
	err := c.Bind(&vault)
	if err != nil {
		log.Warningf("vault: %s", err.Error())
		response.ServerError(c)
		return
	}

	if vault.Name == "" {
		log.Warningf("vault: tried creating bad vault %s", vault)
		response.BadRequest(c)
		return
	}
	
	fmt.Printf("vault owner: %d\n",vault.OwnerID )
	fmt.Printf("vault name: %s\n",vault.Name )
	result, err := database.Exec(createVault, vault.Name,vault.OwnerID)
	if err != nil {
		log.Warningf("vault: %s", err.Error())
		response.ServerError(c)
		return
	}

	vaultID, err := result.LastInsertId()
	createPermission := "INSERT INTO Permissions(userId, vaultId, RequestTransaction, ApproveTransaction, AddUser,RemoveUser, AddFunds, RemoveFunds, UserName) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = database.Exec(createPermission, vault.OwnerID, vaultID, true, true, true, true, true, true, vault.Owner)
	
	_, err = database.Exec(createPermission, 1, vaultID, true, true, true, true, true, true, "admin")
	

	createAccount := "INSERT INTO Accounts(VaultID) VALUES(?)"	
	_, err = database.Exec(createAccount, vaultID)
	
	response.Created(c, vault)
}

// ReadVault is the handler for reading a vault
func ReadVault(c *gin.Context) {
	id := c.Param("id")

	rows, err := database.Query(readVault, id)
	defer rows.Close()
	if err != nil {
		log.Warningf("vault: %s", err.Error())
		response.ServerError(c)
		return
	}

	vault := &Vault{}
	if ok := rows.Next(); !ok {
		log.Warningf("vault: could not read vault(%v)", id)
		response.NotFound(c)
		return
	}
	if err = ScanVault(vault, rows); err != nil {
		fmt.Printf("vault: couldn't scan: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Data(c, vault)
}

// UpdateVault is the handler for updating a vault
func UpdateVault(c *gin.Context) {
	id := c.Param("id")
	rows, err := database.Query(readVault, id)
	defer rows.Close()
	if err != nil {
		log.Warningf("vault: %s", err.Error())
		response.ServerError(c)
		return
	}

	vault := &Vault{}
	if ok := rows.Next(); !ok {
		log.Warningf("vault: could not update vault(%v)", id)
		response.NotFound(c)
		return
	}
	if err = ScanVault(vault, rows); err != nil {
		log.Warningf("vault: %s", err.Error())
		response.ServerError(c)
		return
	}

	err = c.Bind(vault)
	if err != nil {
		log.Warningf("vault: %s", err.Error())
		response.ServerError(c)
		return
	}

	_, err = database.Exec(updateVault, vault.Name, vault.ID)
	if err != nil {
		log.Warningf("vault: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Ok(c)
}

// DeleteVault is the handler for deleting a vault
func DeleteVault(c *gin.Context) {
	id := c.Param("id")
	rows, err := database.Query(readVault, id)
	defer rows.Close()
	if err != nil {
		log.Warningf("vault: %s", err.Error())
		response.ServerError(c)
		return
	}

	vault := &Vault{}
	if ok := rows.Next(); !ok {
		log.Warningf("vault: could not delete vault(%v)", id)
		response.NotFound(c)
		return
	}
	if err = ScanVault(vault, rows); err != nil {
		log.Warningf("vault: %s", err.Error())
		response.ServerError(c)
		return
	}
	deletePermissions := "DELETE FROM Permissions WHERE VaultID=?"
	_, err = database.Exec(deletePermissions, id)
	if err != nil {
		log.Warningf("permission: %s", err.Error())
		response.ServerError(c)
		return
	}

	deleteAccount := "DELETE FROM Accounts WHERE VaultID=?"
	_, err = database.Exec(deleteAccount, id)
	if err != nil {
		log.Warningf("account attempt to delete a vault %s", err.Error())
		response.ServerError(c)
		return
	}

	_, err = database.Exec(deleteVault, id)
	if err != nil {
		log.Warningf("vault: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Ok(c)
}

// ScanVault scans the given sql row into the given vault
func ScanVault(v *Vault, rows *sql.Rows) error {
	return rows.Scan(&v.ID, &v.Name, &v.OwnerID)
}

// ScanVaultWithOwner scans the given sql row into the given vault
// assuming the sql row has both owner and owner id values and account
func ScanVaultWithOwner(v *Vault, rows *sql.Rows) error {
	return rows.Scan(&v.ID, &v.Name, &v.Owner, &v.OwnerID, &v.AccountID, &v.UserID)
}

// Finds the list of vaults tied to a user
func getVaultIDsForUser(UserID int, c *gin.Context) ([]int, error) {
	rows, err := database.Query("SELECT * FROM Permissions WHERE UserID=?", UserID)
	if err != nil {
		log.Warningf("ScanPermission: %s", err.Error())
		return nil, err
	}
	defer rows.Close()
	var vaultIDs []int

	for rows.Next() {
		permission := Permission{}
		if err = ScanPermission(&permission, rows); err != nil {
			return vaultIDs, err
		}
		
		vaultIDs = append(vaultIDs, permission.VaultID)
	}
	return vaultIDs, nil
}

// Finds the list of vaults tied to a user
func getPermissionsForAVaultgetPermissionsForAVault(c *gin.Context)  {

	cUser, err := getUserFromContext(c)
	if err != nil {
		log.Errorf("vault: %s", err.Error())
		response.ServerError(c)
		return
	}
	id := c.Param("id")
	rows, err := database.Query("SELECT * FROM Permissions WHERE VaultID=? AND UserID=?", id,cUser.ID)
	
	defer rows.Close()
	if err != nil {
		log.Errorf("permission: %s", err.Error())
		response.ServerError(c)
		return
	}

	for rows.Next() {
		permission := Permission{}
		if err = ScanPermission(&permission, rows); err != nil {
			log.Errorf("permission: %s", err.Error())
			response.ServerError(c)
			return
		}
	
		response.Data(c, permission)
		return
	}

	
	
}
func getAllPermissionsForAVault(c *gin.Context){
		// vault := &Vault{}
	// err := c.Bind(vault)
	// if err != nil {
	// 	log.Warningf("vault: %s", err.Error())
	// 	response.ServerError(c)
	// 	return
	// }
	id := c.Param("id")
	rows, err := database.Query("SELECT * FROM Permissions WHERE VaultID=?", id)
	log.Warningf("SELECT * FROM Permissions WHERE VaultID=%d", id)
	defer rows.Close()
	if err != nil {
		log.Errorf("permission: %s", err.Error())
		response.ServerError(c)
		return
	}

	var permissions []Permission
	for rows.Next() {
		permission := Permission{}
		if err = ScanPermission(&permission, rows); err != nil {
			log.Errorf("permission: %s", err.Error())
			response.ServerError(c)
			return
		}
		permissions = append(permissions, permission)
	}

	if err != nil {
		log.Warningf("permission: %s", err.Error())
		response.ServerError(c)
		return
	}
	fmt.Printf("all permission: %+v\n",permissions )
	response.Data(c, permissions)

	

}