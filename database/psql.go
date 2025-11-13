package database

import (
	"log"
	"psql-typesense/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {

	dsn := "host=localhost user=postgres password=example dbname=store port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Println(err)
		return err
	}

	err = DB.AutoMigrate(&models.Images{}, &models.Users{})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Connected to Postgres DB")

	return nil
}
