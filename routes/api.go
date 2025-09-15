package routes

import (
	"github.com/gin-gonic/gin"
	"kuhakuanime.com/controller"
	"kuhakuanime.com/middlewares"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/animes", controller.GetAnimes)
	r.GET("/animes/:animeSlug", controller.GetAnimeBySlug)
	r.GET("/animes/:animeSlug/episodes/:episodeNumber", controller.GetEpisodeBySlug)
	
	r.GET("/genres", controller.GetGenres)
	r.GET("/genres/:genreSlug", controller.GetGenreBySlug)
	
	r.GET("episodes/:episodeId/comments", controller.GetCommentsByEpisodeId)
	
	
	

	// Route With Admin Roles
	admin := r.Group("/admin", middlewares.AdminRoles)
	{
		// Manage Anime Data
		admin.POST("/animes", controller.CreateAnime)
		admin.PUT("/animes/:animeId", controller.UpdateAnime)
		admin.DELETE("/animes/:animeId", controller.DeleteAnime)
	
		// Manage Episodes
		admin.POST("/animes/:animeSlug/episodes", controller.CreateEpisode)
		admin.PUT("/animes/:animeSlug/episodes/:episodeNumber", controller.UpdateEpisode)
		admin.DELETE("/animes/:animeSlug/episodes/:episodeNumber", controller.DeleteEpisode)
	
		// Manage Genres
		admin.POST("/genres", controller.CreateGenre)
		admin.PUT("/genres/:genreSlug", controller.UpdateGenre)
		admin.DELETE("/genres/:genreSlug", controller.DeleteGenre)

		admin.POST("/animes/:animeId/genres", controller.AddAnimeGenres)
		admin.DELETE("/animes/:animeId/genres/:genreId", controller.RemoveAnimeGenres)
		admin.PUT("/change-roles", controller.ChangeRoles)
	}


	// Route With User Roles
	auth := r.Group("/", middlewares.Authenticate)
	{
		auth.GET("/favorites", controller.GetFavorites)
		auth.POST("/animes/:animeId/favorites", controller.AddFavorites)
		auth.DELETE("/animes/:animeId/favorites", controller.RemoveFavorites)
	
		auth.POST("/episodes/:episodeId/comments", controller.CreateComment)
		auth.PUT("/comments/:commentId", controller.UpdateComment)
		auth.DELETE("/comments/:commentId", controller.DeleteComment)
	}



	
	r.POST("/login", controller.Login)
	r.POST("/register", controller.Register)

	
	r.Static("/uploads", "./uploads")
}
