package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mohammadhprp/zip-link/internal/handlers"
)

type RouteHandler struct {
	app        *fiber.App
	urlHandler *handlers.URLHandler
}

func NewRouteHandler(
	app *fiber.App,
	urlHandler *handlers.URLHandler,
) *RouteHandler {
	return &RouteHandler{
		app:        app,
		urlHandler: urlHandler,
	}
}

func (h *RouteHandler) Setup() {
	// Routes
	h.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber! 🚀")
	})

	// API group
	api := h.app.Group("/api")

	api.Get("/up", handlers.HealthCheck)

	urls := api.Group("/urls")
	urls.Post("/", h.urlHandler.Create)

	h.app.Get("/:code", h.urlHandler.Get)
}
