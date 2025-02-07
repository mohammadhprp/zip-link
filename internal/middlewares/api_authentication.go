package middlewares

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammadhprp/zip-link/internal/services"
)

func APIAuthenticationMiddleware(service *services.APIKetService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key, ok := c.GetReqHeaders()["X-Api-Kay"]

		if !ok {
			return fmt.Errorf("unauthorized")
		}

		apiKeyInfo, err := service.GetByKey(c.Context(), key[0])

		if err != nil {
			return fmt.Errorf("Invalid API key")
		}

		if time.Now().After(apiKeyInfo.ExpiresAt) {
			return fmt.Errorf("API key expired")
		}

		if apiKeyInfo.RequestCount >= apiKeyInfo.Limit {
			return fmt.Errorf("API request limit reached")
		}
		service.IncreaseRequestCount(c.Context(), apiKeyInfo)

		return c.Next()
	}
}
