package models

import (
	"database/sql"
	"strconv"

	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/database"
	"git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB/response"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	createFavorite = "INSERT INTO Favorites(userId, accountId, Name) Values(?,?,?)"
	deleteFavorite = "DELETE FROM Favorites WHERE ID = ?"
	updateFavorite = "UPDATE Favorites SET name = ? WHERE ID = ?"
	listFavorites  = `SELECT f.ID, f.UserID, f.AccountID, f.Name, u.UserName, v.Name FROM Favorites as f
					JOIN Accounts as a on f.AccountID = a.ID
					LEFT JOIN Users as u on a.UserID = u.ID
					LEFT JOIN Vaults as v on a.VaultID = v.ID
					WHERE f.UserID = ?`
	readFavorite = `SELECT f.ID, f.UserID, f.AccountID, f.Name, u.UserName, v.Name FROM Favorites as f
					JOIN Accounts as a on f.AccountID = a.ID
					LEFT JOIN Users as u on a.UserID = u.ID
					LEFT JOIN Vaults as v on a.VaultID = v.ID
					WHERE f.ID = ?`
)

// Favorite defines information of a Favorite
type Favorite struct {
	ID           int            `form:"id" json:"id"`
	UserID       int            `form:"userId" json:"userId"`
	AccountID    int            `form:"accountId" json:"accountId"`
	FavoriteName string         `form:"favoriteName" json:"favoriteName"`
	Username     sql.NullString `json:"username"`
	Vaultname    sql.NullString `json:"vaultname"`
}

// RouteFavorite sets up the favorite model's HTTP routes
func RouteFavorite(engine *gin.Engine, auth *jwt.GinJWTMiddleware) {
	authGroup := engine.Group("/favorite", auth.MiddlewareFunc())
	authGroup.POST("", CreateFavorite)
	authGroup.PUT(":id", UpdateFavorite)
	authGroup.DELETE("/:id", DeleteFavorite)
	authGroup.GET("", ListFavorites)
}

// CreateFavorite is the handler for creating favorites
func CreateFavorite(c *gin.Context) {
	favorite := Favorite{}
	err := c.Bind(&favorite)

	if err != nil {
		log.Warningf("favorite: %s", err.Error())
		response.ServerError(c)
		return
	}

	if favorite.UserID == 0 || favorite.AccountID == 0 || len(favorite.FavoriteName) == 0 {
		log.Warningf("favorite: tried creating bad favorite %s", favorite)
		response.BadRequest(c)
		return
	}

	_, err = database.Exec(createFavorite, favorite.UserID, favorite.AccountID, favorite.FavoriteName)
	if err != nil {
		log.Warningf("favorite: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Created(c, favorite)
}

// UpdateFavorite is the handler for updating a given favorite by ID
func UpdateFavorite(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Errorf("favorite: %s", err.Error())
		response.ServerError(c)
		return
	}

	favorite := Favorite{}
	err = c.Bind(&favorite)

	rows, err := database.Query(readFavorite, id)
	if err != nil {
		log.Warningf("favorite: %s", err.Error())
		response.ServerError(c)
		return
	}

	if ok := rows.Next(); !ok {
		log.Warningf("favorite: could not find favorite(%d)", id)
		response.NotFound(c)
		return
	}

	rows, err = database.Query(updateFavorite, favorite.FavoriteName, id)
	if err != nil {
		log.Warningf("user: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Ok(c)
}

// ReadFavorite is the handler for reading a specific favorite
func ReadFavorite(c *gin.Context) {
	favorite := Favorite{}
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Errorf("favorite: %s", err.Error())
		response.ServerError(c)
		return
	}

	rows, err := database.Query(readFavorite, id)
	if err != nil {
		log.Warningf("favorite: %s", err.Error())
		response.ServerError(c)
		return
	}

	if err = ScanFavorite(&favorite, rows); err != nil {
		log.Warningf("favorite: couldn't scan: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Data(c, favorite)
}

// ListFavorites is the handler for reading favorites for a given user
func ListFavorites(c *gin.Context) {

	cUser, err := getUserFromContext(c)
	if err != nil {
		log.Warningf("favorite: %s", err.Error())
		response.NotFound(c)
		return
	}

	rows, err := database.Query(listFavorites, cUser.ID)
	defer rows.Close()
	if err != nil {
		log.Warningf("favorite: %s", err.Error())
		response.ServerError(c)
		return
	}

	favorites := []Favorite{}
	favorite := Favorite{}

	for rows.Next() {
		if err = ScanFavorite(&favorite, rows); err != nil {
			log.Warningf("favorite: couldn't scan: %s", err.Error())
			response.ServerError(c)
			return
		}
		favorites = append(favorites, favorite)
	}

	response.Data(c, favorites)
}

// DeleteFavorite is the handler for deleting a favorite
func DeleteFavorite(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if id == 0 {
		log.Warningf("favorite: tried deleting bad favorite")
		response.BadRequest(c)
		return
	}

	_, err := database.Exec(deleteFavorite, id)
	if err != nil {
		log.Warningf("favorite: %s", err.Error())
		response.ServerError(c)
		return
	}

	response.Ok(c)
}

// ScanFavorite scans the given sql row into the given Favorite
func ScanFavorite(f *Favorite, rows *sql.Rows) error {
	return rows.Scan(&f.ID, &f.UserID, &f.AccountID, &f.FavoriteName, &f.Username, &f.Vaultname)
}
