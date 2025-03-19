package middleware

import (
	"auth_project/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

// Protected protect routes
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.LoadEnvVariable("SECRET"))},
		ErrorHandler: jwtError,
	})
}

func AdminOnly(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token) // Get token from context
	claims := user.Claims.(jwt.MapClaims) // Extract claims
	role := claims["role"].(string)       // Get user role

	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Access denied. Admins only."})
	}

	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
