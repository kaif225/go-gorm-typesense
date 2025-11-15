package database

import (
	"fmt"
	"log"
	"os"
	"psql-typesense/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")   // optional
	timezone := os.Getenv("DB_TIMEZONE") // optional

	// fallback defaults (optional)
	if sslmode == "" {
		sslmode = "disable"
	}
	if timezone == "" {
		timezone = "Asia/Kolkata"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone,
	)

	//dsn := "host=localhost user=postgres password=example dbname=store port=5432 sslmode=disable TimeZone=Asia/Kolkata"
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
