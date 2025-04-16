package service

import "github.com/co1seam/tuneflow-backend-auth/internal/ports"

type Service struct {
	ports.IAuthService
}

func NewService() *Service {
	return &Service{}
}
