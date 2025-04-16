package services

import (
	"github.com/co1seam/ember-backend-auth/internal/adapters/repository"
	"github.com/co1seam/ember-backend-auth/internal/core/models"
	"github.com/co1seam/ember-backend-auth/internal/ports"
)

type Service struct {
	Authorization ports.IAuthService
}

func NewService(repos *repository.Repository, opts *models.Options) *Service {
	return &Service{
		Authorization: NewAuthorization(repos.Authorization, opts),
	}
}
