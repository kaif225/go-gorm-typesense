package route

import (
	"psql-typesense/controllers"

	"github.com/gin-gonic/gin"
)

func Unproctected(router *gin.Engine) {

	router.POST("/registration", controllers.RegisterUser)
	router.POST("/updatepassword/:id", controllers.UpdatePassword)
	router.POST("/login", controllers.Login)
	router.POST("/logout", controllers.Logout)
}
