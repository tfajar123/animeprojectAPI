package controller

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"kuhakuanime.com/db/sqlc"
	"kuhakuanime.com/utils"
)

type CreateAnimeInput struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description"`
	Type        string `form:"type"`
}

type UpdateAnimeInput struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	Type        string `form:"type"`
}

func GetAnimeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	queries := db.New(db.DBPool)

	anime, err := queries.GetAnime(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, anime)
}

func GetAnimes(c *gin.Context) {
	queries := db.New(db.DBPool)

	animes, err := queries.GetAnimes(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, animes)
}

func CreateAnime(c *gin.Context) {
	queries := db.New(db.DBPool)


	fmt.Println("title:", c.PostForm("title"))
	fmt.Println("description:", c.PostForm("description"))
	fmt.Println("type:", c.PostForm("type"))
	var anime CreateAnimeInput
	if err := c.ShouldBind(&anime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	image, err := utils.ProcessImageUpload(c, "", "animes")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	newAnime, err := queries.CreateAnime(c, db.CreateAnimeParams{
		Title:       anime.Title,
		Description: pgtype.Text{String: anime.Description, Valid: true},
		Image:       pgtype.Text{String: image, Valid: true},
		Type:        pgtype.Text{String: anime.Type, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, newAnime)
}

func UpdateAnime(c *gin.Context) {
	queries := db.New(db.DBPool)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	oldAnime, err := queries.GetAnime(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var anime UpdateAnimeInput
	if err := c.ShouldBind(&anime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	image, err := utils.ProcessImageUpload(c, oldAnime.Image.String, "animes")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedAnime, err := queries.UpdateAnime(c, db.UpdateAnimeParams{
		ID:          int32(id),
		Title:       anime.Title,
		Description: pgtype.Text{String: anime.Description, Valid: true},
		Image:       pgtype.Text{String: image, Valid: true},
		Type:        pgtype.Text{String: anime.Type, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedAnime)
}

func DeleteAnime(c *gin.Context) {
	queries := db.New(db.DBPool)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	oldAnime, err := queries.GetAnime(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	os.Remove(oldAnime.Image.String)
	err = queries.DeleteAnime(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Anime deleted successfully"})
}