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
	CacheService        *CacheService
}

func NewURLService(db *mongo.Database, cache *CacheService) *URLService {
	return &URLService{
		Collection:          db.Collection("urls"),
		AnalyticsCollection: db.Collection("url_analytics"),
		CacheService:        cache,
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

func (s *URLService) Get(c *fiber.Ctx, filter bson.M) (*models.URL, error) {
	ctx := c.Context()
	shortCode, _ := filter["short_code"].(string)
	cacheKey := buildCacheKey(shortCode)

	ipAddress := utils.GetClientIP(c)
	userAgent := c.Get("User-Agent")

	if cacheKey != "" {
		if cachedURL, err := s.getValidCachedURL(ctx, cacheKey); err == nil {
			go s.processAnalytics(ctx, ipAddress, userAgent, cachedURL)
			return cachedURL, nil
		}
	}

	url, err := s.fetchURLFromDB(ctx, filter)
	if err != nil {
		return nil, err
	}

	go s.processAnalytics(ctx, ipAddress, userAgent, url)

	return url, nil
}

func buildCacheKey(shortCode string) string {
	if shortCode == "" {
		return ""
	}
	return urlCacheKeyPrefix + shortCode
}

func (s *URLService) getValidCachedURL(ctx context.Context, cacheKey string) (*models.URL, error) {
	cachedURL, err := s.getFromCache(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	if cachedURL.ExpiresAt != nil && cachedURL.ExpiresAt.Before(time.Now()) {
		_ = s.CacheService.Delete(ctx, cacheKey)
		return nil, errors.New("Invalid URL")
	}

	return cachedURL, nil
}

func (s *URLService) fetchURLFromDB(ctx context.Context, filter bson.M) (*models.URL, error) {
	var url models.URL
	if err := s.Collection.FindOne(ctx, filter).Decode(&url); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("URL not found")
		}
		return nil, err
	}

	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("Invalid URL")
	}

	if err := s.cacheURL(ctx, &url); err != nil {
		log.Printf("Failed to cache the URL: %v", err)
	}

	return &url, nil
}

func (s *URLService) processAnalytics(ctx context.Context, ipAddress string, userAgent string, url *models.URL) {
	analytics := models.URLAnalytics{
		ID:        primitive.NewObjectID(),
		URLID:     url.ID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}

	if _, err := s.AnalyticsCollection.InsertOne(ctx, analytics); err != nil {
		log.Printf("Failed to log URL analytics: %v", err)
	}

	update := bson.D{
		{Key: "$inc", Value: bson.D{{Key: "click_count", Value: 1}}},
		{Key: "$set", Value: bson.D{{Key: "updated_at", Value: time.Now()}}},
	}

	if _, err := s.Collection.UpdateByID(ctx, url.ID, update); err != nil {
		log.Printf("Failed to update URL click count: %v", err)
	}
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
