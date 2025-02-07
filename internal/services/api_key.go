package services

import (
	"context"
	"errors"
	"log"

	"github.com/mohammadhprp/zip-link/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIKetService struct {
	Collection *mongo.Collection
}

func NewAPIKetService(db *mongo.Database) *APIKetService {
	return &APIKetService{
		Collection: db.Collection("api_keys"),
	}
}

func (s *APIKetService) Create(ctx context.Context, apiKey *models.APIKey) error {
	apiKey.ID = primitive.NewObjectID()

	if _, err := s.Collection.InsertOne(ctx, apiKey); err != nil {
		return err
	}

	return nil
}

func (s *APIKetService) GetByKey(ctx context.Context, key string) (*models.APIKey, error) {
	var apiKey models.APIKey

	filter := bson.M{
		"key": key,
	}

	if err := s.Collection.FindOne(ctx, filter).Decode(&apiKey); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("URL not found")
		}
		return nil, err
	}

	return &apiKey, nil
}

func (s *APIKetService) IncreaseRequestCount(ctx context.Context, apiKey *models.APIKey) error {

	update := bson.D{
		{Key: "$inc", Value: bson.D{{Key: "request_count", Value: 1}}},
	}

	if _, err := s.Collection.UpdateByID(ctx, apiKey.ID, update); err != nil {
		log.Printf("Failed to update API Key request count: %v", err)
	}

	return nil
}
