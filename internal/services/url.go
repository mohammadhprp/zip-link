package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadhprp/zip-link/internal/models"
	"github.com/mohammadhprp/zip-link/internal/requests"
	"github.com/mohammadhprp/zip-link/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	urlCacheKeyPrefix = "url:"
	urlCacheDuration  = 24 * time.Hour
)

type URLService struct {
	Collection          *mongo.Collection
	AnalyticsCollection *mongo.Collection
	CacheService        CacheService
}

func NewURLService(db *mongo.Database, cache *CacheService) *URLService {
	return &URLService{
		Collection:          db.Collection("urls"),
		AnalyticsCollection: db.Collection("url_analytics"),
		CacheService:        *cache,
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

func (s *URLService) Get(ctx context.Context, filter bson.M) (*models.URL, error) {
	var cacheKey string
	if shortCode, ok := filter["short_code"].(string); ok {
		cacheKey = urlCacheKeyPrefix + shortCode
	}

	if cacheKey != "" {
		if cachedURL, err := s.getFromCache(ctx, cacheKey); err == nil {
			if cachedURL.ExpiresAt != nil && cachedURL.ExpiresAt.Before(time.Now()) {
				_ = s.CacheService.Delete(ctx, cacheKey)
				return nil, errors.New("Invalid URL")
			}
			return cachedURL, nil
		}
	}

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

	if err := s.cacheURL(ctx, &url); err != nil {
		log.Fatalf("Failed to cache the url: %v", err)
	}

	return &url, nil
}

func (s *URLService) SetAnalytics(c *fiber.Ctx, url models.URL) error {
	ipAddress := utils.GetClientIP(c)

	analytics := models.URLAnalytics{
		ID:        primitive.NewObjectID(),
		URLID:     url.ID,
		IPAddress: ipAddress,
		UserAgent: c.Get("User-Agent"),
		Referrer:  c.Get("Referer"),
		CreatedAt: time.Now(),
	}

	if _, err := s.AnalyticsCollection.InsertOne(c.Context(), analytics); err != nil {
		return errors.New("failed to log URL analytics: " + err.Error())
	}

	return nil
}

func (s *URLService) cacheURL(ctx context.Context, url *models.URL) error {
	urlJSON, err := json.Marshal(url)
	if err != nil {
		return err
	}

	cacheKey := urlCacheKeyPrefix + url.ShortCode
	return s.CacheService.Set(ctx, cacheKey, string(urlJSON), urlCacheDuration)
}

func (s *URLService) getFromCache(ctx context.Context, key string) (*models.URL, error) {
	cachedData, err := s.CacheService.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var url models.URL
	if err := json.Unmarshal([]byte(cachedData), &url); err != nil {
		return nil, err
	}

	return &url, nil
}
