package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "kuhakuanime.com/db/sqlc"
	"kuhakuanime.com/utils"
)

type User struct {
	Username string `form:"username" binding:"required,alphanum"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type LoginInput struct {
    Username    string `form:"username" binding:"required,alphanum"`
    Password string `form:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input User
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}

	queries := db.New(db.DBPool)
	_ , err = queries.CreateUser(c, db.CreateUserParams{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queries := db.New(db.DBPool)
	user, err := queries.GetUserByUsername(c, input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}
	token, err := utils.GenerateToken(user.Username, int32(user.ID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

