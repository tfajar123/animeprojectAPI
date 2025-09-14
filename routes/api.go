package routes

import (
	"github.com/gin-gonic/gin"
	"kuhakuanime.com/controller"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/animes/:id", controller.GetAnimeByID)
	r.POST("/animes", controller.CreateAnime)
	r.PUT("/animes/:id", controller.UpdateAnime)
	r.DELETE("/animes/:id", controller.DeleteAnime)
	r.GET("/animes", controller.GetAnimes)

	r.GET("/:animeSlug/episodes/:episodeNumber", controller.GetEpisodeBySlug)
	r.POST(":animeSlug/episodes", controller.CreateEpisode)
	r.PUT(":animeSlug/episodes/:episodeNumber", controller.UpdateEpisode)
	r.DELETE(":animeSlug/episodes/:episodeNumber", controller.DeleteEpisode)
	r.GET(":animeSlug/episodes", controller.GetEpisodesByAnimeID)
}
