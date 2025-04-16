package http

import (
	"github.com/co1seam/tuneflow-backend-auth/internal/core/services"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service services.Service
}

func NewHandler(service services.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Router() *fiber.App {
	instance := fiber.New()

	api := instance.Group("/api")
	v1 := api.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.Post("/send-otp", h.sendOtp)
			auth.Post("/sign-up", h.signUp)
			auth.Post("/sign-in", h.signIn)
			auth.Post("/sign-out", h.signOut)
		}
	}

	return instance.Handler
}
