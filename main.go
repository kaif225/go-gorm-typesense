package main

import (
	"log"
	"psql-typesense/controllers"
	"psql-typesense/database"
	"psql-typesense/route"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
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
