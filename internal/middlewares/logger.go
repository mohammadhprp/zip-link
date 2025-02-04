package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func LoggerMiddleware() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "[${time}] ${ip} - ${method} ${path} ${status} - ${latency}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Local",
	})
}
