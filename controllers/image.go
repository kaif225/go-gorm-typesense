package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"psql-typesense/database"
	"psql-typesense/models"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/typesense/typesense-go/typesense/api/pointer"
	"github.com/typesense/typesense-go/v4/typesense/api"
)

func AddImages(c *gin.Context) {
	bucket := os.Getenv("BUCKET_NAME")
	region := os.Getenv("AWS_REGION")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	category := c.Param("category")
	file, err := c.FormFile("image")

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not provided"})
		return
	}
	// Validate file type (optional but recommended)
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
	fileHeader := file.Header.Get("Content-Type")
	isValid := false
	for _, t := range allowedTypes {
		if fileHeader == t {
			isValid = true
			break
		}
	}
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only images allowed"})
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is empty"})
		return
	}
	defer fileContent.Close()

	s3Key := fmt.Sprintf("%s/%s", category, file.Filename)

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &s3Key,
		Body:   fileContent,
	})

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading images"})
		return
	}

	s3Url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, s3Key)

	imageDoc := models.Images{
		Category:   category,
		FileName:   file.Filename,
		S3URL:      s3Url,
		S3Key:      s3Key,
		UploadedAt: time.Now().Format(time.RFC3339),
	}

	database.DB.Create(&imageDoc)
	instertedId := imageDoc.ID
	id := strconv.Itoa(instertedId)
	typesenseDoc := models.TypesenseImage{

		ID:         id,
		Category:   category,
		FileName:   file.Filename,
		S3URL:      s3Url,
		S3Key:      s3Key,
		UploadedAt: imageDoc.UploadedAt,
	}

	_, err = database.TypesenseClient.Collection("images").Documents().Create(ctx, typesenseDoc, &api.DocumentIndexParameters{})

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading images"})
		return
	}
	err = database.DB.Model(&models.Images{}).Where("id = ?", instertedId).Update("typesense_synced", true).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Updating image sync status failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "image uploaded successfully"})

}

func SearchImages(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var imageResponse []models.ImageResponse
	bucket := os.Getenv("BUCKET_NAME")
	query := c.DefaultQuery("q", "*")
	//category := c.Query("category")
	searchParameters := &api.SearchCollectionParams{
		Q:       pointer.String(query),
		Infix:   pointer.String("fallback"),
		QueryBy: pointer.String("file_name,category"),
		SortBy:  pointer.String("uploaded_at:desc"),
		PerPage: pointer.Int(8),
	}

	// Add category filter if provided
	// if category != "" {
	// 	searchParameters.FilterBy = pointer.String(fmt.Sprintf("category:=%s", category))
	// }

	searchResult, err := database.TypesenseClient.Collection("images").Documents().Search(context.Background(), searchParameters)

	if err != nil {
		log.Println("Typesense search error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	if searchResult.Hits != nil {

		for _, hit := range *searchResult.Hits {
			if doc := hit.Document; doc != nil {
				s3Key := (*doc)["s3_key"].(string)
				fileName := (*doc)["file_name"].(string)
				category := (*doc)["category"].(string)
				s3URL := (*doc)["s3_url"].(string)
				uploadedAt := (*doc)["uploaded_at"].(string)

				request, err := presignedClient.PresignGetObject(ctx, &s3.GetObjectInput{
					Bucket: &bucket,
					Key:    &s3Key,
				}, func(opts *s3.PresignOptions) {
					opts.Expires = time.Duration(60 * time.Second)
				})
				signedURL := s3URL
				if err == nil {
					signedURL = request.URL
				} else {
					log.Println("Error generating pre-signed URL:", err)
				}

				imageResponse = append(imageResponse, models.ImageResponse{
					Category:   category,
					FileName:   fileName,
					S3Key:      s3Key,
					S3URL:      s3URL,
					SignedURL:  signedURL,
					UploadedAt: uploadedAt,
				})
			}

		}
	}

	c.JSON(http.StatusOK, imageResponse)
}
