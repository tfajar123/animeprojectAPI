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
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error() })
		return
	}

	queries := db.New(db.DBPool)
	_, err = queries.GetUser(c, userID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error() })
		return
	}

	c.Set("user_id", userID)

	c.Next()
}

func AdminRoles(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == ""{
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	userID, err := utils.VerifyToken(token)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error() })
		return
	}

	queries := db.New(db.DBPool)
	user, err := queries.GetUser(c, userID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error() })
		return
	}

	if user.Role != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You Are Not Admin"})
		return
	}

	c.Set("user_id", userID)

	c.Next()
}