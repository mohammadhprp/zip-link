package utils

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetAppURL() string {
	appURL := os.Getenv("APP_URL")

	if appURL == "" {
		log.Println("Warning: APP_URL is not set")
	}

	return appURL
}

// GetClientIP extracts the real client IP address from the request
func GetClientIP(c *fiber.Ctx) string {
	if ip := c.Get("X-Forwarded-For"); ip != "" {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if ip := c.Get("X-Real-IP"); ip != "" {
		return ip
	}

	ipAddress := c.IP()
	if net.ParseIP(ipAddress) == nil {
		ipAddress = "unknown"
	}

	return ipAddress
}
