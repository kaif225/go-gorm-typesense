package main

import (
	"log"
	"psql-typesense/controllers"
	"psql-typesense/database"
	"psql-typesense/route"
	"psql-typesense/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	err := utils.LoadVaultSecretsWithRetry(".env", 5)
	//err := godotenv.Load()
	if err != nil {
		log.Println(err)
		return
	}
	err = database.Connect()
	if err != nil {
		log.Println(err)
		return
	}
	err = database.TsConnect()
	if err != nil {
		log.Println(err)
		return
	}
	controllers.TypeSenseInitImages()
	controllers.TypesenseInitUsers()
	controllers.S3Init()

	router := gin.Default()
	route.Protected(router)
	route.Unproctected(router)

	router.Run(":8007")
}
