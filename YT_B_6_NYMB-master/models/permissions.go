package models

import (
	"database/sql"
	// "fmt"

	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/database"
	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/response"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	listPermissions  = "SELECT * FROM Permissions"
	createPermission = "INSERT INTO Permissions(userId, vaultId, RequestTransaction, ApproveTransaction, AddUser,RemoveUser, AddFunds, RemoveFunds, UserName) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	createPermissionSimple = "INSERT INTO Permissions(userId, vaultId) VALUES(?, ?)"
	readPermission   = "SELECT * FROM Permissions WHERE ID=?"
	updatePermission = "UPDATE Permissions SET RequestTransaction=?, ApproveTransaction=?, " +
		"AddUser=?, RemoveUser=?, AddFunds=?, RemoveFunds=? WHERE ID=?"
	deletePermission = "DELETE FROM Permissions WHERE ID=?"
)

// Permission defines information of a permission
type Permission struct {
	ID                 int  `form:"ID" json:"ID"`
	UserID             int  `form:"userID" json:"userID"`
	VaultID            int  `form:"vaultID" json:"vaultID"`
	RequestTransaction bool `form:"requestTransaction" json:"requestTransaction"`
	ApproveTransaction bool `form:"approveTransaction" json:"approveTransaction"`
	AddUser            bool `form:"addUser" json:"addUser"`
	RemoveUser         bool `form:"removeUser" json:"removeUser"`
	AddFunds           bool `form:"addFunds" json:"addFunds"`
	RemoveFunds        bool `form:"removeFunds" json:"removeFunds"`
	UserName           string `form:"userName" json:"userName"`
}

// RoutePermission sets up the permission model's HTTP routes
func RoutePermission(engine *gin.Engine, auth *jwt.GinJWTMiddleware) {
	authGroup := engine.Group("/permissions", auth.MiddlewareFunc())
	authGroup.GET("", ListPermission)
	authGroup.POST("", CreatePermission)
	authGroup.GET("/:id", ReadPermission)
	authGroup.PUT("/:id", UpdatePermission)
	authGroup.DELETE("/:id", DeletePermission)
	
}

// ListPermission is the handler for listing permissions
func ListPermission(c *gin.Context) {
	query := GetQuery(listPermissions, c.Request.URL.Query())
	rows, err := database.Query(query)
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

	response.Data(c, permissions)
}

// CreatePermission is the handler for creating permissions
func CreatePermission(c *gin.Context) {
	permission := Permission{}
	err := c.Bind(&permission)
	if err != nil {
		log.Warningf("permission: %s", err.Error())
		response.ServerError(c)
		return
	}

	if permission.UserID == 0 || permission.VaultID == 0 {
		log.Warningf("permission: tried creating bad permission %s", permission)
		response.BadRequest(c)
		return
	}
	log.Warningf("vault id: %d", permission.VaultID)
	_, err = database.Exec(createPermission, permission.UserID, permission.VaultID, permission.RequestTransaction, permission.ApproveTransaction, permission.AddUser, permission.RemoveUser, permission.AddFunds, permission.RemoveFunds, permission.UserName)
	if err != nil {
		log.Warningf("permission: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Created(c, permission)
}

// ReadPermission is the handler for reading a permission
//
// In the URL, "id" represents the permission's user id.
func ReadPermission(c *gin.Context) {
	id := c.Param("id")

	rows, err := database.Query(readPermission, id)
	defer rows.Close()
	if err != nil {
		log.Warningf("permission: %s", err.Error())
		response.ServerError(c)
		return
	}

	permission := &Permission{}
	if ok := rows.Next(); !ok {
		log.Warningf("permission: could not read permission(%v)", id)
		response.NotFound(c)
		return
	}
	if err = ScanPermission(permission, rows); err != nil {
		log.Warningf("permission: couldn't scan: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Data(c, permission)
}

// UpdatePermission is the handler for updating a permission
//
// In the URL, "id" represents the permission's user id.
func UpdatePermission(c *gin.Context) {
	// Bind request body to permission in order to retrieve userId and vaultId
	permission := &Permission{}
	err := c.Bind(permission)
	if err != nil {
		log.Warningf("permission: %s", err.Error())
		response.ServerError(c)
		return
	}
	if permission.UserID == 0 || permission.VaultID == 0 {
		log.Warningf("permission: tried updating bad permission %s", permission)
		response.BadRequest(c)
		return
	}

	// Query database for given permission
	rows, err := database.Query(readPermission, permission.UserID, permission.VaultID)
	defer rows.Close()
	if err != nil {
		log.Warningf("permission: %s", err.Error())
		response.ServerError(c)
		return
	}
	if ok := rows.Next(); !ok {
		log.Warningf("permission: could not update %s", permission)
		response.NotFound(c)
		return
	}
	if err = ScanPermission(permission, rows); err != nil {
		log.Warningf("permission: %s", err.Error())
		response.ServerError(c)
		return
	}

	// rebind request body to overwrite changed values and update the database
	c.Bind(permission)
	_, err = database.Exec(updatePermission, permission.RequestTransaction, permission.ApproveTransaction,
		permission.AddUser, permission.RemoveUser, permission.AddFunds, permission.RemoveFunds,
		permission.UserID, permission.VaultID)
	if err != nil {
		log.Warningf("permission: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Ok(c)
}

// DeletePermission is the handler for deleting a permission
func DeletePermission(c *gin.Context) {
	id := c.Param("id")
	rows, err := database.Query(readPermission, id)
	defer rows.Close()
	if err != nil {
		log.Warningf("permission: %s", err.Error())
		response.ServerError(c)
		return
	}

	permission := &Permission{}
	if ok := rows.Next(); !ok {
		log.Warningf("permission: could not delete permission(%v)", id)
		response.NotFound(c)
		return
	}
	if err = ScanPermission(permission, rows); err != nil {
		log.Warningf("permission: %s", err.Error())
		response.ServerError(c)
		return
	}

	_, err = database.Exec(deletePermission, id)
	if err != nil {
		log.Warningf("permission: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Ok(c)
}

// ScanPermission scans the given sql row into the given permission
func ScanPermission(p *Permission, rows *sql.Rows) error {
	return rows.Scan(&p.ID, &p.UserID, &p.VaultID, &p.RequestTransaction, &p.ApproveTransaction,
		&p.AddUser, &p.RemoveUser, &p.AddFunds, &p.RemoveFunds, &p.UserName)
}


