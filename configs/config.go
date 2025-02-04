package configs

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func FiberConfig() fiber.Config {
	appName := os.Getenv("APP_NAME")

	return fiber.Config{
		AppName:               appName,
		IdleTimeout:           5 * time.Second,
		DisableStartupMessage: false,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	}
}

