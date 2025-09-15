package controller

import (
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
	Slug		string `form:"slug"`
}

type AnimeResponse struct {
    ID          int32          `json:"id"`
    Title       string         `json:"title"`
    Description string         `json:"description"`
    Slug        string         `json:"slug"`
    Type        string         `json:"type"`
    Image       string         `json:"image"`
    ImagePort   string         `json:"image_port"`
	Genres      []db.Genre     `json:"genres"`
    Episodes    []db.Episode   `json:"episodes"`
}

type AnimeWithGenre struct {
	db.Anime
	Genres []db.Genre `json:"genres"`
}


func GetAnimeBySlug(c *gin.Context) {
	animeSlug := c.Param("animeSlug")

	queries := db.New(db.DBPool)

	anime, err := queries.GetAnimeBySlug(c, animeSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	episodes, err := queries.GetEpisodesByAnimeSlug(c, animeSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	genres, err := queries.GetGenresByAnimeId(c, int32(anime.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := AnimeResponse{
        ID:          anime.ID,
        Title:       anime.Title,
        Description: anime.Description.String,
        Slug:        anime.Slug,
        Type:        anime.Type.String,
        Image:       anime.Image.String,
        ImagePort:   anime.ImagePort.String,
		Genres:      genres,
        Episodes:    episodes,
		
    }

	c.JSON(http.StatusOK, response)
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

	var anime CreateAnimeInput
	if err := c.ShouldBind(&anime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	image, err := utils.ProcessImageUpload(c, "image", "", "animes_landscape")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	image_port, err := utils.ProcessImageUpload(c, "image_port", "", "animes_portrait")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	slug := utils.Slugify(anime.Title)
	if existing, err := queries.GetAnimeBySlug(c, slug); err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Judul anime ini sudah ada",
			"anime": existing,
		})
		return
	}

	newAnime, err := queries.CreateAnime(c, db.CreateAnimeParams{
		Title:       anime.Title,
		Description: pgtype.Text{String: anime.Description, Valid: true},
		Image:       pgtype.Text{String: image, Valid: true},
		Type:        pgtype.Text{String: anime.Type, Valid: true},
		Slug:        slug,
		ImagePort:   pgtype.Text{String: image_port, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, newAnime)
}

func UpdateAnime(c *gin.Context) {
	queries := db.New(db.DBPool)

	idStr := c.Param("animeId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	oldAnime, err := queries.GetAnimeById(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var anime UpdateAnimeInput
	if err := c.ShouldBind(&anime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	anime.Title       = utils.IfEmpty(anime.Title, oldAnime.Title, "")
	anime.Description = utils.IfEmpty(anime.Description, oldAnime.Description.String, "")
	anime.Type        = utils.IfEmpty(anime.Type, oldAnime.Type.String, "")
	
	if anime.Title != oldAnime.Title {
		slug := utils.Slugify(anime.Title)
		if existing, err := queries.GetAnimeBySlug(c, slug); err == nil && existing.ID != oldAnime.ID {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Judul anime ini sudah ada",
				"anime": existing,
			})
			return
		}
		anime.Slug = slug
	} else {
		anime.Slug = oldAnime.Slug
	}

	image, err := utils.ProcessImageUpload(c, "image", oldAnime.Image.String, "animes_landscape")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	image_port, err := utils.ProcessImageUpload(c, "image_port", oldAnime.ImagePort.String, "animes_portrait")
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
		Slug:        pgtype.Text{String: anime.Slug, Valid: true}.String,
		ImagePort:   pgtype.Text{String: image_port, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedAnime)
}

func DeleteAnime(c *gin.Context) {
	queries := db.New(db.DBPool)

	idStr := c.Param("animeId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	oldAnime, err := queries.GetAnimeById(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	os.Remove(oldAnime.Image.String)
	os.Remove(oldAnime.ImagePort.String)
	err = queries.DeleteAnime(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Anime deleted successfully"})
}