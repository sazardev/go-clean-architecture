package middleware

import (
	"strings"

	"go-clean-architecture/internal/infrastructure/auth/jwt"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware(tokenService *jwt.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Bearer token is required",
			})
		}

		// Extract the token
		token := jwt.ExtractTokenFromBearer(authHeader)
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token format",
			})
		}

		// Validate the token
		claims, err := tokenService.ValidateToken(token)
		if err != nil {
			if err == jwt.ErrExpiredToken {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Token has expired",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Set user information in context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_roles", claims.Roles)
		c.Locals("user_permissions", claims.Permissions)
		c.Locals("user_claims", claims)

		return c.Next()
	}
}

// OptionalAuthMiddleware validates JWT tokens but doesn't require them
func OptionalAuthMiddleware(tokenService *jwt.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next() // No auth header, continue without authentication
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Next() // Invalid format, continue without authentication
		}

		// Extract the token
		token := jwt.ExtractTokenFromBearer(authHeader)
		if token == "" {
			return c.Next() // Invalid token format, continue without authentication
		}

		// Validate the token
		claims, err := tokenService.ValidateToken(token)
		if err != nil {
			return c.Next() // Invalid token, continue without authentication
		}

		// Set user information in context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_roles", claims.Roles)
		c.Locals("user_permissions", claims.Permissions)
		c.Locals("user_claims", claims)

		return c.Next()
	}
}

// RefreshTokenMiddleware handles token refresh
func RefreshTokenMiddleware(tokenService *jwt.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from request body or header
		var request struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if request.RefreshToken == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Refresh token is required",
			})
		}

		// Validate the refresh token (allowing expired tokens for refresh)
		claims, err := tokenService.ValidateToken(request.RefreshToken)
		if err != nil && err != jwt.ErrExpiredToken {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid refresh token",
			})
		}

		// Generate new token
		newToken, err := tokenService.RefreshToken(claims)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to refresh token",
			})
		}

		return c.JSON(fiber.Map{
			"access_token": newToken,
			"token_type":   "Bearer",
		})
	}
}
