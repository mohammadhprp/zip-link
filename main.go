package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadhprp/zip-link/configs"
	"github.com/mohammadhprp/zip-link/internal/handlers"
	"github.com/mohammadhprp/zip-link/internal/middlewares"
	"github.com/mohammadhprp/zip-link/internal/routes"
	"github.com/mohammadhprp/zip-link/internal/services"
)

func main() {
	app := fiber.New(configs.FiberConfig())

	db := configs.ConnectMongoDB()
	defer configs.MongoClient.Disconnect(nil)

	cache := configs.ConnectRedisDB()

	app.Use(middlewares.LoggerMiddleware())

	cacheService := services.NewCacheService(cache)

	apiKeyService := services.NewAPIKeyService(db)

	urlService := services.NewURLService(db, cacheService)
	urlHandler := handlers.NewURLHandler(urlService)

	routeHandler := routes.NewRouteHandler(app, urlHandler, apiKeyService)
	routeHandler.Setup()

	log.Fatal(app.Listen(":3000"))
}
