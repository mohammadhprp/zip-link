package services

import (
	"context"
	"errors"
	"time"

	"github.com/mohammadhprp/zip-link/internal/models"
	"github.com/mohammadhprp/zip-link/internal/requests"
	"github.com/mohammadhprp/zip-link/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type URLService struct {
	Collection *mongo.Collection
}

func NewURLService(db *mongo.Database) *URLService {
	return &URLService{
		Collection: db.Collection("urls"),
	}
}

func (s *URLService) Create(ctx context.Context, request requests.StoreURLRequest) (*models.URL, error) {
	shortCode := utils.GenerateShortCode()

	url := &models.URL{
		ID:          primitive.NewObjectID(),
		OriginalURL: request.URL,
		ShortCode:   shortCode,
		ClickCount:  0,
		Metadata:    make(models.Map),
		ExpiresAt:   request.ExpiresAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := s.Collection.InsertOne(ctx, url)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *URLService) Get(ctx context.Context, code string) (*models.URL, error) {
	filter := bson.M{"short_code": code}
	var url models.URL

	err := s.Collection.FindOne(ctx, filter).Decode(&url)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("URL not found")
		}
		return nil, err
	}

	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("Invalid URL")
	}

	return &url, nil
}
