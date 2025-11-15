package main

import (
	"log"
	"psql-typesense/controllers"
	"psql-typesense/database"
	"psql-typesense/route"
	"psql-typesense/utils"
)

func main() {
	err := utils.LoadVaultSecretsWithRetry("/vault/secrets/secrets.txt", 5)
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

	router := route.Router()

	router.Run(":8007")
}
