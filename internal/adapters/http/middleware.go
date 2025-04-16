package http

import (
	"fmt"
	fiberjwt "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
)

func (h *Handler) authMiddleware(requiredClaims ...string) func(ctx *fiber.Ctx) error {
	return fiberjwt.New(fiberjwt.Config{
		SuccessHandler: func(ctx *fiber.Ctx) error {
			user := ctx.Locals("user").(*jwt.Token)
			claims, ok := user.Claims.(jwt.MapClaims)
			if !ok {
				return fiber.ErrUnauthorized
			}

			if err := validateClaims(claims, requiredClaims); err != nil {
				return err
			}

			if tokenType, ok := claims["type"].(string); !ok || tokenType != "access" {
				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid token type",
				})
			}

			if err := setContextClaims(ctx, claims, requiredClaims); err != nil {
				return err
			}

			return ctx.Next()
		},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		},
		SigningKey: fiberjwt.SigningKey{
			JWTAlg: fiberjwt.HS256,
			Key:    h.opts.Config.Token.Secret,
		},
		TokenLookup: "header:Authorization,cookie:access_token",
		AuthScheme:  "Bearer",
	})
}

func validateClaims(claims jwt.MapClaims, required []string) error {
	for _, key := range required {
		if _, exists := claims[key]; !exists {
			return fiber.ErrUnauthorized
		}
	}
	return nil
}

func setContextClaims(ctx *fiber.Ctx, claims jwt.MapClaims, keys []string) error {
	for _, key := range keys {
		value := claims[key]

		switch key {
		case "sub":
			strValue, ok := value.(string)
			if !ok {
				return fiber.NewError(fiber.StatusUnauthorized, "invalid sub claim")
			}

			id, err := strconv.ParseUint(strValue, 10, 64)
			if err != nil {
				return fiber.NewError(fiber.StatusUnauthorized, "invalid sub format")
			}
			ctx.Locals(key, uint(id))

		case "type", "aud":
			strValue, ok := value.(string)
			if !ok {
				return fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("invalid %s claim", key))
			}
			ctx.Locals(key, strValue)

		default:
			ctx.Locals(key, value)
		}
	}
	return nil
}
