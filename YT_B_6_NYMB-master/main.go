package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/auth"
	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/database"
	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/gcus"
	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/models"
	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/transaction"

	"github.com/gin-gonic/gin"
	"github.com/igm/sockjs-go/sockjs"
	log "github.com/sirupsen/logrus"
)

var (
	dev    bool
	debug  bool
	remote bool
	port   int
	db     *sql.DB
)

func main() {

	// parse flags/arguments
	args()

	// check database connection
	checkDatabase()

	// start transaction queue
	go transaction.Run()

	// start GDAX Currency Update Service
	go gcus.Start()

	// setup HTTP router
	route()
}

func args() {

	flag.BoolVar(&dev, "dev", false, "Start server in dev mode.")
	flag.BoolVar(&debug, "debug", false, "Start server in debug mode.")
	flag.BoolVar(&remote, "remote", false, "Use remote database connection.")
	flag.IntVar(&port, "port", 8080, "Port on which to start server.")
	flag.Parse()

	log.SetOutput(os.Stdout)
	if !dev {
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(log.WarnLevel)
	}
	if debug {
		log.SetLevel(log.DebugLevel)
	}

}

func checkDatabase() {

	err := database.Open(remote)
	if err != nil {
		panic(err)
	}

}

func route() {

	r := gin.Default()

	authMiddleware := auth.JWTMiddleware(r)

	// setup websocket upgrade request endpoint
	transaction.SetAuthMiddleware(authMiddleware)
	handler := sockjs.NewHandler("/ws", sockjs.DefaultOptions, transaction.WebSocketHandler)
	r.Any("/ws/*path", gin.WrapH(handler))

	models.RouteAccount(r, authMiddleware)
	models.RouteBalance(r, authMiddleware)
	models.RouteCurrency(r, authMiddleware)
	models.RouteTransaction(r, authMiddleware)
	models.RouteUser(r, authMiddleware)
	models.RouteVault(r, authMiddleware)
	models.RouteFavorite(r, authMiddleware)
	models.RoutePermission(r, authMiddleware)

	r.Run(":" + fmt.Sprintf("%d", port))
}
