package route

import (
	"psql-typesense/controllers"
	"psql-typesense/middlewares"

	"github.com/gin-gonic/gin"
)

func Protected(router *gin.Engine) {

	protected := router.Group("/")
	protected.Use(middlewares.JWT())
	protected.POST("/upload/:category", controllers.AddImages)
	protected.GET("/users/search", controllers.SearchUsers)
	protected.GET("/images/search", controllers.SearchImages)
	protected.POST("/syncimages", controllers.SyncSchemasImages)
	protected.POST("/syncusers", controllers.SyncSchemasImages)

}
