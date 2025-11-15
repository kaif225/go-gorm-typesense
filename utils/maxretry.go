package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func LoadVaultSecretsWithRetry(path string, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		err := godotenv.Load(path)
		if err == nil {
			log.Printf("✅ Secrets loaded in %d seconds", i*2)
			return nil
		}

		if i < maxRetries-1 {
			log.Printf("⏳ Waiting for secrets... (attempt %d/%d)", i+1, maxRetries)
			time.Sleep(2 * time.Second)
		}
	}
	return fmt.Errorf("❌ failed to load secrets after %d attempts", maxRetries)
}
