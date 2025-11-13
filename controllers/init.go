package controllers

import (
	"context"
	"fmt"
	"log"
	"psql-typesense/database"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/typesense/typesense-go/typesense/api/pointer"
	"github.com/typesense/typesense-go/v4/typesense/api"
)

var s3Client *s3.Client
var presignedClient *s3.PresignClient

func S3Init() {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}

	s3Client = s3.NewFromConfig(cfg)
	presignedClient = s3.NewPresignClient(s3Client)
}

func TypeSenseInitImages() error {
	ctx := context.Background()

	schema := &api.CollectionSchema{
		Name: "images",
		Fields: []api.Field{
			{
				Name: "id",
				Type: "string",
			}, {
				Name:  "category",
				Type:  "string",
				Facet: pointer.True(),
			},
			{
				Name:  "file_name",
				Type:  "string",
				Infix: pointer.True(),
			},
			{
				Name: "s3_key",
				Type: "string",
			}, {
				Name: "uploaded_at",
				Type: "string",
				Sort: pointer.True(),
			},
		},
		//DefaultSortingField: pointer.String("category"),
	}
	_, err := database.TypesenseClient.Collection("images").Retrieve(ctx)
	if err == nil {
		log.Println("Typesense Schema already exists")
		return nil
	}
	_, err = database.TypesenseClient.Collections().Create(ctx, schema)
	if err != nil {
		log.Println("The Error is :", err)
		return err
	}
	log.Println("Typesense Schema created successfully")

	return nil
}

func TypesenseInitUsers() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	schemaUser := &api.CollectionSchema{
		Name: "users",
		Fields: []api.Field{
			{
				Name: "id",
				Type: "string",
			}, {
				Name:  "first_name",
				Type:  "string",
				Facet: pointer.True(),
			},
			{
				Name:  "last_name",
				Type:  "string",
				Infix: pointer.True(),
			},
			{
				Name: "email",
				Type: "string",
			}, {
				Name: "username",
				Type: "string",
				Sort: pointer.True(),
			}, {
				Name: "password_changed_at",
				Type: "string",
				Sort: pointer.True(),
			}, {
				Name: "user_created_at",
				Type: "string",
				Sort: pointer.True(),
			}, {
				Name: "inactive_status",
				Type: "string",
			}, {
				Name: "role",
				Type: "string",
			},
		},
		//DefaultSortingField: pointer.String("category"),
	}
	_, err := database.TypesenseClient.Collection("images").Retrieve(ctx)
	if err == nil {
		log.Println("Typesense Schema already exists")
		return nil
	}
	_, err = database.TypesenseClient.Collections().Create(ctx, schemaUser)
	if err != nil {
		log.Println("The Error is :", err)
		return err
	}

	log.Println("Typesense User Schema created successfully")
	return nil
}
