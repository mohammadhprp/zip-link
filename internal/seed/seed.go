package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/mohammadhprp/zip-link/configs"
	"github.com/mohammadhprp/zip-link/internal/models"
	"github.com/mohammadhprp/zip-link/internal/services"
)

func main() {
	db := configs.ConnectMongoDB()
	defer configs.MongoClient.Disconnect(nil)

	apiKeyService := services.NewAPIKetService(db)

	seedAPIKeys(apiKeyService)
}

func seedAPIKeys(apiKeyService *services.APIKetService) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	apiKeys := []models.APIKey{
		{
			Key:          uuid.New().String(),
			ExpiresAt:    time.Now().Add(30),
			RequestCount: 0,
			Limit:        1000,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Key:          uuid.New().String(),
			ExpiresAt:    time.Now().Add(30),
			RequestCount: 98,
			Limit:        100,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Key:          uuid.New().String(),
			ExpiresAt:    time.Now(),
			RequestCount: 0,
			Limit:        100,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	for _, apiKey := range apiKeys {
		err := apiKeyService.Create(ctx, &apiKey)
		if err != nil {
			log.Printf("Failed to seed API key %+v: %v\n", apiKey, err)
		} else {
			log.Printf("Seeded API key: %+v\n", apiKey)
		}
	}
}
