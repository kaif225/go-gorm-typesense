package route

import (
	"psql-typesense/controllers"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	router := gin.Default()
	router.POST("/upload/:category", controllers.AddImages)
	router.GET("/images/search", controllers.SearchImages)
	return router
}
