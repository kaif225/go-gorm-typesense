package controllers

import (
	"context"
	"log"
	"net/http"
	"psql-typesense/database"
	"psql-typesense/models"
	"psql-typesense/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/typesense/typesense-go/typesense/api/pointer"
	"github.com/typesense/typesense-go/v4/typesense/api"
)

func RegisterUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var user models.Users

	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": "Invalid Request Body"})
		return
	}

	validate := validator.New()

	err = validate.Struct(user)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	hashPass, err := utils.HashPassword(user.Password)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error encrypting password"})
		return
	}

	user.Password = hashPass
	user.UserCreatedAt = time.Now().Format(time.RFC3339)

	err = database.DB.Create(&user).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error in creating account"})
		return
	}

	tysenseDoc := models.TypesenseUser{
		ID:             strconv.Itoa(user.ID),
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		Username:       user.Username,
		UserCreatedAt:  user.UserCreatedAt,
		InactiveStatus: user.InactiveStatus,
		Role:           user.Role,
	}

	_, err = database.TypesenseClient.Collection("users").Documents().Create(ctx, tysenseDoc, &api.DocumentIndexParameters{})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user index is not synced"})
		return
	}
	err = database.DB.Model(&models.Users{}).Where("id = ?", user.ID).Update("typesense_synced", true).Error
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "registration completed"})
}

func SearchUsers(c *gin.Context) {

	query := c.DefaultQuery("q", "*")

	searchParameters := &api.SearchCollectionParams{
		Q:       pointer.String(query),
		Infix:   pointer.String("fallback"),
		QueryBy: pointer.String("first_name,last_name,user_created_at"),
		//FilterBy: pointer.String("num_employees:>100"),
		SortBy: pointer.String("user_created_at:desc"),
	}

	searchResult, err := database.TypesenseClient.Collection("users").Documents().Search(context.Background(), searchParameters)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "not able to get users"})
		return
	}

	c.JSON(http.StatusOK, searchResult)
}
