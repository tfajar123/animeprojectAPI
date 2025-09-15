package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "kuhakuanime.com/db/sqlc"
	"kuhakuanime.com/utils"
)

type GenreInput struct {
	Name string `form:"name" binding:"required"`
}

type AnimeGenre struct {
	GenreID int32 `form:"genre_id" binding:"required"`
}

func GetGenres(c *gin.Context) {
	queries := db.New(db.DBPool)

	genres, err := queries.GetGenres(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, genres)
}

func GetGenreBySlug(c *gin.Context) {
	genreSlug := c.Param("genreSlug")

	queries := db.New(db.DBPool)

	genre, err := queries.GetGenreBySlug(c, genreSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	animes, err := queries.GetAnimesByGenreId(c, int32(genre.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, animes)
}

func CreateGenre(c *gin.Context) {
	queries := db.New(db.DBPool)

	var genre GenreInput
	if err := c.ShouldBind(&genre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	genreSlug := utils.Slugify(genre.Name)
	if existing, err := queries.GetGenreBySlug(c, genreSlug); err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Genre name already exists",
			"genre": existing,
		})
		return
	}

	newGenre, err := queries.CreateGenre(c, db.CreateGenreParams{
		Name: genre.Name,
		Slug: genreSlug,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newGenre)
}

func UpdateGenre(c *gin.Context) {
	queries := db.New(db.DBPool)

	genreSlug := c.Param("genreSlug")
	oldGenre, err := queries.GetGenreBySlug(c, genreSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var genre GenreInput
	if err := c.ShouldBind(&genre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	genre.Name = utils.IfEmpty(genre.Name, oldGenre.Name, "")
		
	if genre.Name != oldGenre.Name {
		genreSlug = utils.Slugify(genre.Name)
		if existing, err := queries.GetGenreBySlug(c, genreSlug); err == nil && existing.ID != oldGenre.ID {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Genre name already exists",
				"genre": existing,
			})
			return
		}
	} else {
		genreSlug = oldGenre.Slug
	}

	_, err = queries.UpdateGenre(c, db.UpdateGenreParams{
		ID:   int32(oldGenre.ID),
		Name: genre.Name,
		Slug: genreSlug,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Genre updated successfully"})
}

func DeleteGenre(c *gin.Context) {
	queries := db.New(db.DBPool)

	genreSlug := c.Param("genreSlug")
	genre, err := queries.GetGenreBySlug(c, genreSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = queries.DeleteGenre(c, int32(genre.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Genre deleted successfully"})
}


func AddAnimeGenres(c *gin.Context) {
	queries := db.New(db.DBPool)

	idStr := c.Param("animeId")
	animeId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var animeGenre AnimeGenre
	if err := c.ShouldBind(&animeGenre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = queries.CheckAnimeGenreExists(c, db.CheckAnimeGenreExistsParams{
		AnimeID: int32(animeId),
		GenreID: int32(animeGenre.GenreID),
	})
	if err == nil {
		c.JSON(http.StatusNoContent, gin.H{"message": "Genre already added to anime"})
		return
	}

	_, err = queries.CreateAnimeGenre(c, db.CreateAnimeGenreParams{
		AnimeID: int32(animeId),
		GenreID: int32(animeGenre.GenreID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Genre added to anime successfully"})
}

func RemoveAnimeGenres(c *gin.Context) {
	queries := db.New(db.DBPool)
	genreIdStr := c.Param("genreId")
	genreId, err := strconv.Atoi(genreIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	animeIdStr := c.Param("animeId")
	animeId, err := strconv.Atoi(animeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid anime ID"})
		return
	}

	err = queries.DeleteAnimeGenre(c, db.DeleteAnimeGenreParams{
		AnimeID: int32(animeId),
		GenreID: int32(genreId),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Genre removed from anime successfully"})
}