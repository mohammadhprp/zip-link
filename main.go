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

	// Connect to MongoDB
	db := configs.ConnectMongoDB()
	defer configs.MongoClient.Disconnect(nil)

	app.Use(middlewares.LoggerMiddleware())

	urlService := services.NewURLService(db)
	urlHandler := handlers.NewURLHandler(urlService)

	routeHandler := routes.NewRouteHandler(app, urlHandler)
	routeHandler.Setup()

	log.Fatal(app.Listen(":3000"))
}
