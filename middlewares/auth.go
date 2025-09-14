package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "kuhakuanime.com/db/sqlc"
	"kuhakuanime.com/utils"
)

func Authenticate(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == ""{
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	userID, err := utils.VerifyToken(token)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	queries := db.New(db.DBPool)
	_, err = queries.GetUser(c, userID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	c.Set("user_id", userID)

	c.Next()
}