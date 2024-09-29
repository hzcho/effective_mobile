package handler

import (
	"song_lib/internal/group"

	_ "song_lib/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           song library API API
// @version         1.0
// @description     This is an API server for working with songs
//
//	@host			localhost:8080
func InitRoutes(router *gin.Engine, groups group.Groups) {
	api := router.Group("/api/v1")
	{
		songs := api.Group("/songs")
		{
			songs.GET("/info", groups.Song.GetLib)
			songs.GET("/:id/verses", groups.Song.GetVerses)
			songs.POST("/", groups.Song.Add)
			songs.PUT("/:id", groups.Song.Update)
			songs.DELETE("/:id", groups.Song.Delete)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
