package http

import (
	"github.com/co1seam/ember-backend-auth/internal/core/models"
	"github.com/co1seam/ember-backend-auth/internal/core/services"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *services.Service
	opts    *models.Options
}

func NewHandler(service *services.Service, opts *models.Options) *Handler {
	return &Handler{
		service: service,
		opts:    opts,
	}
}

func (h *Handler) Router(instance *fiber.App) *fiber.App {
	api := instance.Group("/api")
	v1 := api.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			signUp := auth.Group("/sign-up")
			{
				signUp.Post("/send-otp", h.sendOtp)
				signUp.Post("/verify", h.signUp)

			}
			auth.Post("/sign-in", h.signIn)
			auth.Post("/sign-out", h.signOut)
		}

		protected := v1.Group("")
		{
			protected.Get("/test", h.test)
		}
	}

	return instance
}

func (h *Handler) test(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"jwt": "successful",
	})
}
