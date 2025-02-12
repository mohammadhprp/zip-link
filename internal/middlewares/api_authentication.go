package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadhprp/zip-link/internal/services"
)

func APIAuthenticationMiddleware(service *services.APIKeyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-Api-Key")

		if key == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		apiKeyInfo, err := service.GetByKey(c.Context(), key)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}

		if time.Now().After(apiKeyInfo.ExpiresAt) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API key expired",
			})
		}

		if apiKeyInfo.RequestCount >= apiKeyInfo.Limit {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "API request limit reached",
			})
		}

		if err := service.IncreaseRequestCount(c.Context(), apiKeyInfo); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update request count",
			})
		}

		return c.Next()
	}
}
