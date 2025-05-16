package rpc

import (
	authv1 "github.com/co1seam/ember-backend-api-contracts/gen/go/auth"
	"github.com/co1seam/ember-backend-auth/internal/core/models"
	"github.com/co1seam/ember-backend-auth/internal/core/services"
)

type Handler struct {
	Authorization authv1.AuthServer
	opts          *models.Options
}

func NewHandler(service *services.Service, opts *models.Options) *Handler {
	return &Handler{
		Authorization: NewAuthorization(service.Authorization, opts),
		opts:          opts,
	}
}
