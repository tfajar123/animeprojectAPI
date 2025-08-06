package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "kuhakuanime.com/db/sqlc"
	"kuhakuanime.com/utils"
)

type CreateEpisodeInput struct {
	Episode_number int    `form:"episode_number" binding:"required"`
	Episode_url    string `form:"episode_url" binding:"required"`
	Anime_id       int    `form:"anime_id"`
}

type UpdateEpisodeInput struct {
	Episode_number int    `form:"episode_number"`
	Episode_url    string `form:"episode_url"`
	Anime_id       int    `form:"anime_id"`
}

func GetEpisodeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	queries := db.New(db.DBPool)

	episode, err := queries.GetAnime(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, episode)
}

func GetEpisodesByAnimeID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	
	
	queries := db.New(db.DBPool) 

	episodes, err := queries.GetEpisodesByAnimeId(c, int32(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, episodes)
}

func CreateEpisode(c *gin.Context) {
	queries := db.New(db.DBPool)

	var episode CreateEpisodeInput
	if err := c.ShouldBind(&episode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newEpisode, err := queries.CreateEpisode(c, db.CreateEpisodeParams{
		EpisodeNumber: int32(episode.Episode_number),
		EpisodeUrl:    episode.Episode_url,
		AnimeID:       int32(episode.Anime_id),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, newEpisode)
}

func UpdateEpisode(c *gin.Context) {
	queries := db.New(db.DBPool)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	oldEpisode , err := queries.GetEpisode(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var episode UpdateEpisodeInput
	if err := c.ShouldBind(&episode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	episode.Episode_number = utils.IfEmpty(episode.Episode_number, int(oldEpisode.EpisodeNumber), 0)
	episode.Episode_url = utils.IfEmpty(episode.Episode_url, oldEpisode.EpisodeUrl, "")

	updatedEpisode, err := queries.UpdateEpisode(c, db.UpdateEpisodeParams{
		ID:            int32(id),
		EpisodeNumber: int32(episode.Episode_number),
		EpisodeUrl:    episode.Episode_url,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, updatedEpisode)
}

func DeleteEpisode(c *gin.Context) {
	queries := db.New(db.DBPool)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = queries.DeleteEpisode(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Episode deleted successfully"})
}
