package auth

import (
	"log"
	"time"

	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/database"
	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/models"
	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/response"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

const jwtSecretKey = "secretKeyShouldBeReplaced"

// JWTMiddleware sets up and returns JWT authentication middleware
func JWTMiddleware(r *gin.Engine) *jwt.GinJWTMiddleware {
	middleware := &jwt.GinJWTMiddleware{
		Realm:         "test zone",
		Key:           []byte(jwtSecretKey),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		Authenticator: authenticator,
		Authorizator:  authorizator,
		Unauthorized: func(c *gin.Context, code int, message string) {
			response.Unauthorized(c)
		},
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",

		TimeFunc: time.Now,
	}

	r.POST("/login", middleware.LoginHandler)
	auth := r.Group("/auth", middleware.MiddlewareFunc())
	auth.GET("/token", middleware.RefreshHandler)
	auth.GET("/user", userHandler)

	return middleware
}

// authenticator checks the context's email and password credentials
// and returns the email and true if the credentials are valid
func authenticator(email, password string, c *gin.Context) (string, bool) {
	rows, err := database.Query("Select * FROM Users WHERE Email=?", email)
	if err != nil {
		log.Fatal(err)
		response.ServerError(c)
		return email, false
	}
	defer rows.Close()

	usr := &models.User{}
	for rows.Next() {
		if err = models.ScanUser(usr, rows); err != nil {
			log.Fatal(err)
			response.ServerError(c)
			return email, false
		}
		if email == usr.Email && password == usr.Password {
			return email, true
		}
	}

	return email, false
}

// authorizator checks the context's jwt validity
// when a user reaches a protected api endpoint
func authorizator(email string, c *gin.Context) bool {
	if err := jwt.ExtractClaims(c).Valid(); err != nil {
		return false
	}

	return true
}

// userHandler defines the handler which returns a user's account information.
func userHandler(c *gin.Context) {
	email, ok := c.Get("userID")
	if !ok {
		response.ServerError(c)
		log.Fatal("main.authUserHandler: couldn't get email from context")
		return
	}
	rows, err := database.Query(models.ListUsers+" WHERE u.Email=?", email)
	if err != nil {
		log.Fatal(err)
		response.ServerError(c)
		return
	}
	defer rows.Close()

	usr := &models.User{}
	for rows.Next() {
		if err = models.ScanUserWithRoleAndAccount(usr, rows); err != nil {
			log.Fatal(err)
			response.ServerError(c)
			return
		}
		if email == usr.Email {
			response.Data(c, usr)
			return
		}
	}

	response.ServerError(c)
	return
}
