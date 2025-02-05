package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadhprp/zip-link/internal/requests"
	"github.com/mohammadhprp/zip-link/internal/services"
	"github.com/mohammadhprp/zip-link/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type URLHandler struct {
	Service *services.URLService
}

func NewURLHandler(service *services.URLService) *URLHandler {
	return &URLHandler{Service: service}
}

func (h *URLHandler) Create(c *fiber.Ctx) error {
	var payload requests.StoreURLRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := payload.Validate(); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	url, err := h.Service.Create(c.Context(), payload)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	appUrl := utils.GetAppURL()
	zipedUrl := fmt.Sprintf("%s/%s", appUrl, url.ShortCode)

	response := fiber.Map{
		"expires_at":   url.ExpiresAt,
		"original_url": url.OriginalURL,
		"ziped_url":    zipedUrl,
	}

	return c.Status(http.StatusCreated).JSON(response)
}

func (h *URLHandler) Get(c *fiber.Ctx) error {
	code := c.Params("code")

	filters := bson.M{
		"short_code": code,
	}

	url, err := h.Service.Get(c.Context(), filters)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.Service.SetAnalytics(c, *url); err != nil {
		return errors.New("Invalid URL")
	}

	return c.Redirect(url.OriginalURL, http.StatusMovedPermanently)
}
