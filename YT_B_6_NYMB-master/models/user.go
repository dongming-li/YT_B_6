package models

import (
	"database/sql"
	"fmt"

	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/database"
	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/response"

	"github.com/VividCortex/mysqlerr"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	// ListUsers is exported for use with auth.go
	ListUsers = "SELECT u.ID, u.Username, u.Email, u.FirstName, u.LastName, r.RoleID, a.ID FROM Users AS u " +
		"JOIN UserRoles AS r ON u.ID=r.UserID " +
		"LEFT JOIN Accounts AS a ON u.ID=a.UserID"
	createUser = "INSERT INTO Users(username, email, password, firstname, lastname) VALUES(?, ?, ?, ?, ?)"
	readUser   = "SELECT * FROM Users WHERE ID=?"
	updateUser = "UPDATE Users SET username=?, email=?, password=?, firstname=?, lastname=? WHERE ID=?"
	deleteUser = "DELETE FROM Users WHERE ID=?"

	readUserWithRoleByEmail = "SELECT u.ID, Email, RoleID FROM Users AS u " +
		"JOIN UserRoles AS r ON u.ID=r.UserID WHERE u.Email=?"

	readUserWithAccountByEmail = "SELECT ID, UserName, Email, Password, FirstName, LastName, a.ID " +
		"AS VaultID FROM Users AS u LEFT JOIN Accounts AS a ON u.ID=a.UserID WHERE u.Email=?;"

	createUserRole = "INSERT INTO UserRoles(UserID, RoleID) VALUES(?, 2)"
	readUserRole   = "SELECT * FROM UserRoles WHERE UserID=?"
	updateUserRole = "UPDATE UserRoles SET RoleID=? WHERE UserID=?"
	deleteUserRole = "DELETE FROM UserRoles WHERE UserID=?"

	adminRoleID = 1
	userRoleID  = 2
)

// User defines information of a user
type User struct {
	ID        int    `json:"id"`
	Username  string `form:"username" json:"username"`
	Email     string `form:"email" json:"email"`
	Password  string `form:"password" json:"password"`
	FirstName string `form:"firstname" json:"firstname"`
	LastName  string `form:"lastname" json:"lastname"`
	RoleID    int    `form:"role" json:"role"`
	Account   int64  `form:"account" json:"account"`
}

// RouteUser sets up the user model's HTTP routes
func RouteUser(engine *gin.Engine, auth *jwt.GinJWTMiddleware) {
	authGroup := engine.Group("/user", auth.MiddlewareFunc())
	authGroup.GET("", ListUser)
	engine.POST("/user", CreateUser)
	authGroup.GET("/:id", ReadUser)
	authGroup.PUT("/:id", UpdateUser)
	authGroup.DELETE("/:id", DeleteUser)
}

// ListUser is the handler for listing users
func ListUser(c *gin.Context) {
	query := GetQuery(ListUsers, c.Request.URL.Query())
	rows, err := database.Query(query)
	defer rows.Close()
	if err != nil {
		dErr, _ := err.(database.Error)
		log.Warningf("user: %s", dErr.Error())
		if dErr.Code == mysqlerr.ER_BAD_FIELD_ERROR {
			response.BadRequestMessage(c, "Invalid field specified")
			return
		}
		response.ServerError(c)
		return
	}

	var users []User
	for rows.Next() {
		user := User{}
		if err = ScanUserWithRoleAndAccount(&user, rows); err != nil {
			log.Errorf("user: %s", err.Error())
			response.ServerError(c)
			return
		}
		user.Password = ""
		users = append(users, user)
	}

	if err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Data(c, users)
}

// CreateUser is the handler for creating users
func CreateUser(c *gin.Context) {
	user := User{}
	err := c.Bind(&user)
	if err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		log.Warningf("user: tried creating bad user %s", user)
		response.BadRequest(c)
		return
	}

	result, err := database.Exec(createUser,
		user.Username, user.Email, user.Password, user.FirstName, user.LastName)
	if err != nil {
		dErr, _ := err.(database.Error)
		log.Warningf("user: %s", dErr.Error())
		if dErr.Code == mysqlerr.ER_DUP_ENTRY {
			response.BadRequestMessage(c, "Username or Email already taken.")
			return
		}
		response.ServerError(c)
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}

	_, err = database.Exec(createUserRole, userID)
	if err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}

	user.RoleID = 2
	response.Created(c, user)
}

// ReadUser is the handler for reading a user
func ReadUser(c *gin.Context) {
	id := c.Param("id")

	rows, err := database.Query(readUserWithAccountByEmail, id)
	defer rows.Close()
	if err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}

	user := &User{}
	if ok := rows.Next(); !ok {
		log.Warningf("user: could not read user(%v)", id)
		response.NotFound(c)
		return
	}
	if err = ScanUserWithAccount(user, rows); err != nil {
		log.Warningf("user: couldn't scan: %s", err.Error())
		response.ServerError(c)
		return
	}

	user.Password = ""
	response.Data(c, user)
}

// UpdateUser is the handler for updating a user
func UpdateUser(c *gin.Context) {

	// get the user needing update by url paramater
	id := c.Param("id")
	rows, err := database.Query(readUser, id)
	if err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}
	user := &User{}
	if ok := rows.Next(); !ok {
		log.Warningf("user: could not update user(%v)", id)
		response.NotFound(c)
		return
	}
	if err = ScanUser(user, rows); err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}
	rows.Close()

	// get the user associated with the current context to check their role
	email, _ := c.Get("userID")
	rows, err = database.Query(readUserWithRoleByEmail, email)
	if err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}
	cUser := &User{}
	if ok := rows.Next(); !ok {
		log.Warningf("user: could not update user(%v)", id)
		response.NotFound(c)
		return
	}
	if err = ScanUserWithEmailAndRole(cUser, rows); err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}
	rows.Close()

	// check if the context user is the user needing update or is an admin
	if user.Email != email && cUser.RoleID != adminRoleID {
		log.Warningf("user: user(%s) tried to update user (%s)", email, user.Email)
		response.Unauthorized(c)
		return
	}

	password := user.Password
	// overwrite the old user needing update with the request fields
	err = c.Bind(user)
	if err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}

	if user.Password == "" {
		user.Password = password
	}

	// update the user table
	_, err = database.Exec(updateUser,
		user.Username, user.Email, user.Password, user.FirstName, user.LastName, user.ID)
	if err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}

	// if the user role has changed and the context user is admin, update the user role table
	if user.RoleID != 0 && cUser.RoleID == adminRoleID {
		_, err = database.Exec(updateUserRole, user.RoleID, user.ID)
		if err != nil {
			log.Warningf("user: %s", err.Error())
			response.ServerError(c)
			return
		}
	}

	response.Ok(c)
}

// DeleteUser is the handler for deleting a user
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	rows, err := database.Query(readUser, id)
	if err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}
	user := &User{}
	if ok := rows.Next(); !ok {
		log.Warningf("user: could not delete user(%v)", id)
		response.NotFound(c)
		return
	}
	if err = ScanUser(user, rows); err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}
	rows.Close()

	// get the user associated with the current context to check their role
	email, _ := c.Get("userID")
	rows, err = database.Query(readUserWithRoleByEmail, email)
	if err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}
	cUser := &User{}
	if ok := rows.Next(); !ok {
		log.Warningf("user: could not delete user(%v)", id)
		response.NotFound(c)
		return
	}
	if err = ScanUserWithEmailAndRole(cUser, rows); err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}
	rows.Close()

	// check if the context user is the user needing delete or is an admin
	if user.Email != email && cUser.RoleID != adminRoleID {
		log.Warningf("user: user(%s) tried to delete user (%s)", email, user.Email)
		response.Unauthorized(c)
		return
	}

	// delete the user role
	_, err = database.Exec(deleteUserRole, id)
	if err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}

	// delete the user
	_, err = database.Exec(deleteUser, id)
	if err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Ok(c)
}

// ScanUser scans the given sql row into the given user
func ScanUser(u *User, rows *sql.Rows) error {
	return rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.FirstName, &u.LastName)
}

// ScanUserWithAccount scans the given sql row into the given user with an account ID corresponding to the user
func ScanUserWithAccount(u *User, rows *sql.Rows) error {
	var accountID sql.NullInt64
	err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.FirstName, &u.LastName, &accountID)
	if accountID.Valid {
		u.Account = accountID.Int64
	}

	return err
}

// ScanUserWithRoleAndAccount scans the given sql row into the user
// assuming the returned row also has role id
func ScanUserWithRoleAndAccount(u *User, rows *sql.Rows) error {
	var accountID sql.NullInt64
	err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FirstName, &u.LastName, &u.RoleID, &accountID)
	if accountID.Valid {
		u.Account = accountID.Int64
	}

	return err
}

// ScanUserWithEmailAndRole scans the given sql row into the user
// assuming the returned row only has id, email, and role specified
func ScanUserWithEmailAndRole(u *User, rows *sql.Rows) error {
	return rows.Scan(&u.ID, &u.Email, &u.RoleID)
}

func (u User) String() string {
	return fmt.Sprintf("User(username=%s, email=%s, firstname=%s, lastname=%s, role=%d)",
		u.Username, u.Email, u.FirstName, u.LastName, u.RoleID)
}

// Gets the user from the context
func getUserFromContext(c *gin.Context) (User, error) {
	// Get User from context
	email, _ := c.Get("userID")
	rows, err := database.Query(readUserWithRoleByEmail, email)
	cUser := User{}
	if ok := rows.Next(); !ok {
		err = fmt.Errorf("could not find user(%v)", email)
		return cUser, err
	}
	err = ScanUserWithEmailAndRole(&cUser, rows)
	rows.Close()

	return cUser, err
}
