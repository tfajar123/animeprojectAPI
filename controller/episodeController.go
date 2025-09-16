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
}

type UpdateEpisodeInput struct {
	Episode_url    string `form:"episode_url"`
}

func GetEpisodeBySlug(c *gin.Context) {
    animeSlug := c.Param("animeSlug")
    episodeNumStr := c.Param("episodeNumber")

    episodeNum, err := strconv.Atoi(episodeNumStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode number"})
        return
    }

    queries := db.New(db.DBPool)

    episode, err := queries.GetEpisodeBySlugAndNumber(c, db.GetEpisodeBySlugAndNumberParams{
        Slug:          animeSlug,
        EpisodeNumber: int32(episodeNum),
    })
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Episode not found"})
        return
    }

    c.JSON(http.StatusOK, episode)
}

func CreateEpisode(c *gin.Context) {
	idStr := c.Param("animeId")
	animeId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}


	anime, err := db.New(db.DBPool).GetAnimeById(c, int32(animeId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}
	
	queries := db.New(db.DBPool)

	var episode CreateEpisodeInput
	if err := c.ShouldBind(&episode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = queries.CheckEpisodeExists(c, db.CheckEpisodeExistsParams{
		EpisodeNumber: int32(episode.Episode_number),
		AnimeID:       int32(anime.ID),
	})
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Episode already exists"})
		return
	}

	newEpisode, err := queries.CreateEpisode(c, db.CreateEpisodeParams{
		EpisodeNumber: int32(episode.Episode_number),
		EpisodeUrl:    episode.Episode_url,
		AnimeID:       int32(anime.ID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, newEpisode)
}

func UpdateEpisode(c *gin.Context) {
	queries := db.New(db.DBPool)

	idStr := c.Param("animeId")
	animeId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	anime, err := db.New(db.DBPool).GetAnimeById(c, int32(animeId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}

    episodeNumStr := c.Param("episodeNumber")

    episodeNum, err := strconv.Atoi(episodeNumStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode number"})
        return
    }

	oldEpisode , err := queries.GetEpisodeBySlugAndNumber(c, db.GetEpisodeBySlugAndNumberParams{
        Slug:          anime.Slug,
        EpisodeNumber: int32(episodeNum),
    })
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var episode UpdateEpisodeInput
	if err := c.ShouldBind(&episode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	episode.Episode_url = utils.IfEmpty(episode.Episode_url, oldEpisode.EpisodeUrl, "")

	updatedEpisode, err := queries.UpdateEpisode(c, db.UpdateEpisodeParams{
		ID:            int32(oldEpisode.EpisodeID),
		EpisodeUrl:    episode.Episode_url,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, updatedEpisode)
}

func DeleteEpisode(c *gin.Context) {
	
	idStr := c.Param("animeId")
	animeId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	anime, err := db.New(db.DBPool).GetAnimeById(c, int32(animeId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found"})
		return
	}

    episodeNumStr := c.Param("episodeNumber")
	
    episodeNum, err := strconv.Atoi(episodeNumStr)
    if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode number"})
        return
    }
	queries := db.New(db.DBPool)
	
    episode, err := queries.GetEpisodeBySlugAndNumber(c, db.GetEpisodeBySlugAndNumberParams{
        Slug:          anime.Slug,
        EpisodeNumber: int32(episodeNum),
    })

	err = queries.DeleteEpisode(c, int32(episode.EpisodeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Episode deleted successfully"})
}
