package routes

import (
	"github.com/gin-gonic/gin"
	"kuhakuanime.com/controller"
	"kuhakuanime.com/middlewares"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/animes/:animeSlug", controller.GetAnimeBySlug)
	r.GET("/animes", controller.GetAnimes)

	r.GET(":animeSlug/episodes/:episodeNumber", controller.GetEpisodeBySlug)
	

	r.GET("/genres", controller.GetGenres)
	r.GET("/genres/:genreSlug", controller.GetGenreBySlug)

	


	adminAuthenticated := r.Group("/")
	adminAuthenticated.Use(middlewares.AdminRoles)
	// Manage Anime Data
	adminAuthenticated.POST("/animes", controller.CreateAnime)
	adminAuthenticated.PUT("/animes/:animeId", controller.UpdateAnime)
	adminAuthenticated.DELETE("/animes/:animeId", controller.DeleteAnime)

	// Manage Episodes
	adminAuthenticated.POST(":animeSlug/episodes", controller.CreateEpisode)
	adminAuthenticated.PUT(":animeSlug/episodes/:episodeNumber", controller.UpdateEpisode)
	adminAuthenticated.DELETE(":animeSlug/episodes/:episodeNumber", controller.DeleteEpisode)

	// Manage Genres
	adminAuthenticated.POST("/genres", controller.CreateGenre)
	adminAuthenticated.PUT("/genres/:genreSlug", controller.UpdateGenre)
	adminAuthenticated.DELETE("/genres/:genreSlug", controller.DeleteGenre)
	adminAuthenticated.POST("/genre", controller.AddAnimeGenres)
	adminAuthenticated.DELETE("/genre/:genreId/:animeId", controller.RemoveAnimeGenres)

	authenticated := r.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.PUT("/change-roles", controller.ChangeRoles)


	
	r.POST("/login", controller.Login)
	r.POST("/register", controller.Register)

	
	r.Static("/uploads", "./uploads")
}
