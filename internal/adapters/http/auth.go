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

	err := h.service.Authorization.SendOTP(reqCtx, req.Email)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "otp successful send"})
}

func (h *Handler) verifyOtp(ctx *fiber.Ctx) error {
	reqCtx := ctx.UserContext()

	var req models.VerifyOtpRequest

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	email, err := h.service.Authorization.VerifyOTP(reqCtx, req.OTP)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"email": email})
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
