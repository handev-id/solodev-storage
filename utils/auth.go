package utils

import (
	"crypto/subtle"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// BearerAuth validates Authorization: Bearer <token>.
func BearerAuth(secretToken string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if secretToken == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "upload secret key is not configured",
			})
		}

		token, ok := parseBearerToken(c.Get(fiber.HeaderAuthorization))
		if !ok || subtle.ConstantTimeCompare([]byte(token), []byte(secretToken)) != 1 {
			c.Set("WWW-Authenticate", "Bearer")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or missing bearer token",
			})
		}

		return c.Next()
	}
}

func parseBearerToken(headerValue string) (string, bool) {
	value := strings.TrimSpace(headerValue)
	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 {
		return "", false
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", false
	}

	return token, true
}
