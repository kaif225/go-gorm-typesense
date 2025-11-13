package controllers

import (
	"context"
	"log"
	"net/http"
	"psql-typesense/database"
	"psql-typesense/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/typesense/typesense-go/v4/typesense/api"
)

func SyncSchemasImages(c *gin.Context) {

	var images []models.Images
	result := database.DB.Where("typesense_synced = ?", false).Find(&images)
	//result := database.DB.Find(&images)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not able to find image"})
		return
	}

	if len(images) == 0 {
		c.JSON(200, gin.H{"message": "All documents are already synced"})
		return
	}

	synced := 0
	failed := 0
	for _, img := range images {
		typesenseDoc := &models.TypesenseImage{
			ID:         strconv.Itoa(img.ID),
			Category:   img.Category,
			FileName:   img.FileName,
			S3Key:      img.S3Key,
			S3URL:      img.S3URL,
			UploadedAt: img.UploadedAt,
		}

		_, err := database.TypesenseClient.Collection("images").Documents().Create(context.Background(), typesenseDoc, &api.DocumentIndexParameters{})

		if err != nil {
			log.Printf("Failed to sync document %s: %v\n", strconv.Itoa(img.ID), err)
			failed++
			continue
		}

		err = database.DB.Model(&models.Images{}).Where("id = ?", img.ID).
			Update("typesense_synced", true).Error

		if err != nil {
			log.Printf("Falied to update sync status for document %s, %v\n", strconv.Itoa(img.ID), err)
			return
		}
		synced++

		c.JSON(200, gin.H{
			"message": "Sync completed",
			"synced":  synced,
			"failed":  failed,
		})
	}
}

func SyncSchemasUsers(c *gin.Context) {

	var users []models.Users
	result := database.DB.Where("typesense_synced = ?", false).Find(&users)
	//result := database.DB.Find(&images)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not able to find image"})
		return
	}

	if len(users) == 0 {
		c.JSON(200, gin.H{"message": "All documents are already synced"})
		return
	}

	synced := 0
	failed := 0
	for _, user := range users {
		typesenseDoc := &models.TypesenseUser{
			ID:        strconv.Itoa(user.ID),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Username:  user.Username,
			Role:      user.Role,
		}

		_, err := database.TypesenseClient.Collection("users").Documents().Create(context.Background(), typesenseDoc, &api.DocumentIndexParameters{})

		if err != nil {
			log.Printf("Failed to sync document %s: %v\n", strconv.Itoa(user.ID), err)
			failed++
			continue
		}

		err = database.DB.Model(&models.Users{}).Where("id = ?", user.ID).
			Update("typesense_synced", true).Error

		if err != nil {
			log.Printf("Falied to update sync status for document %s, %v\n", strconv.Itoa(user.ID), err)
			return
		}
		synced++

		c.JSON(200, gin.H{
			"message": "Sync completed",
			"synced":  synced,
			"failed":  failed,
		})
	}
}

// func ResetSchemaSync(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	err := utils.DeleteSchema("images", "users")
// 	if err != nil {
// 		log.Println(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting schema"})
// 		return
// 	}
// 	err = TypeSenseInit()
// 	if err != nil {
// 		log.Println(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "error Creating schema"})
// 		return
// 	}
// 	synced := 0
// 	failed := 0
// 	var images []models.Images
// 	result := database.DB.Find(&images)

// 	if result.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not able to find images"})
// 		return
// 	}

// 	for _, img := range images {
// 		typesenseDoc := &models.TypesenseImage{
// 			ID:         strconv.Itoa(img.ID),
// 			Category:   img.Category,
// 			FileName:   img.FileName,
// 			S3Key:      img.S3Key,
// 			S3URL:      img.S3URL,
// 			UploadedAt: img.UploadedAt,
// 		}

// 		_, err = database.TypesenseClient.Collection("images").Documents().Create(ctx, typesenseDoc, &api.DocumentIndexParameters{})
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "error syning indexes"})
// 			failed++
// 			continue
// 		}

// 		err = database.DB.Model(&models.Images{}).Where("typesense_synced = ?", false).
// 			Update("typesense_synced = ?", true).Error
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not able to find images"})
// 			return
// 		}
// 		synced++
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Reset and Sync completed",
// 		"synced":  synced,
// 		"failed":  failed,
// 	})
// }
