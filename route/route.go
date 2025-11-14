package route

import (
	"psql-typesense/controllers"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	router := gin.Default()
	router.POST("/upload/:category", controllers.AddImages)
	router.GET("/image/search", controllers.SearchImages)
	router.GET("/user/search", controllers.SearchUsers)
	router.POST("/syncimages", controllers.SyncSchemasImages)
	router.POST("/syncusers", controllers.SyncSchemasUsers)
	router.POST("/registration", controllers.RegisterUser)
	return router
}
