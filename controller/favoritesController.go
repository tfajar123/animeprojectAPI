package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "kuhakuanime.com/db/sqlc"
)

func GetFavorites(c *gin.Context) {
	queries := db.New(db.DBPool)

	userID, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID not found"})
		return
	}

	favorites, err := queries.GetFavoritesByUserId(c, int32(userID.(int32)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, favorites)
}

func AddFavorites(c *gin.Context) {
	queries := db.New(db.DBPool)

	idStr := c.Param("animeId")
	animeId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	
	userID, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID not found"})
		return
	}

	_, err = queries.CheckFavoriteExists(c, db.CheckFavoriteExistsParams{
		AnimeID: int32(animeId),
		UserID:  int32(userID.(int32)),
	})
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Favorite already exists"})
		return
	}

	_, err = queries.CreateFavorite(c, db.CreateFavoriteParams{
		AnimeID: int32(animeId),
		UserID:  int32(userID.(int32)),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite added successfully"})
}

func RemoveFavorites(c *gin.Context) {
	queries := db.New(db.DBPool)

	idStr := c.Param("animeId")
	animeId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	err = queries.DeleteFavorite(c, int32(animeId))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite removed successfully"})
}