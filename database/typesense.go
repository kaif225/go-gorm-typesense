package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/typesense/typesense-go/v4/typesense"
)

var TypesenseClient *typesense.Client

func TsConnect() error {
	typesenseUrl := os.Getenv("TYPESENSE_URI")
	apiKey := os.Getenv("TYPESENSE_API_KEY")
	if typesenseUrl == "" || apiKey == "" {
		return errors.New("typesense Error: Api key or URL is missing")
	}
	client := typesense.NewClient(
		typesense.WithServer(typesenseUrl),
		typesense.WithAPIKey(apiKey))

	resp, err := client.Health(context.TODO(), 5*time.Second)

	if err != nil {
		log.Printf("Error connecting to typesense %v:  %v", typesenseUrl, err)
		return err
	}

	if !resp {
		return fmt.Errorf("typesense health check failed: %+v", resp)
	}
	TypesenseClient = client
	log.Println("Connected to Typesense")
	return nil
}
