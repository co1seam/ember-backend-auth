package http

import (
	"fmt"
	"github.com/co1seam/ember-backend-auth/internal/core/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func (h *Handler) sendOtp(ctx *fiber.Ctx) error {
	reqCtx := ctx.UserContext()

	var req models.SendOtpRequest

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	id, err := h.service.Authorization.Create(reqCtx, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"id": id})
}

func (h *Handler) signUp(ctx *fiber.Ctx) error {
	reqCtx := ctx.UserContext()

	var req models.SignUpRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	id, err := h.service.Authorization.Create(reqCtx, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"id": id})
}

func (h *Handler) signIn(ctx *fiber.Ctx) error {
	reqCtx := ctx.UserContext()

	var req models.SignInRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	userID, err := h.service.Authorization.Read(reqCtx, req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.createTokens(ctx, userID.(int)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"user": userID.(int)})
}

func (h *Handler) signOut(ctx *fiber.Ctx) error {

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (h *Handler) createTokens(ctx *fiber.Ctx, userID int) error {
	refreshCookie, err := h.createJWT(time.Hour*72, jwt.MapClaims{
		"sub":  fmt.Sprint(userID),
		"type": "refresh",
		"aud":  "admin",
	})
	if err != nil {
		return err
	}

	accessToken, err := h.createJWT(time.Minute*15, jwt.MapClaims{
		"sub":  fmt.Sprint(userID),
		"type": "access",
		"aud":  "admin",
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshCookie,
		Expires:  time.Now().Add(72 * time.Hour),
		HTTPOnly: true,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HTTPOnly: true,
	})

	return nil
}
