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

	r.GET("/episodes/:id", controller.GetEpisodeByID)
	r.POST("/episodes", controller.CreateEpisode)
	r.PUT("/episodes/:id", controller.UpdateEpisode)
	r.DELETE("/episodes/:id", controller.DeleteEpisode)
	r.GET("/episodes/animes/:id", controller.GetEpisodesByAnimeID)
}
