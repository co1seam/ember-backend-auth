package ports

import (
	"context"
)

type (
	IAuthRepo interface {
		CRUD
	}

	IAuthService interface {
		CRUD
		SendOTP(ctx context.Context, email string) error
		VerifyOTP(ctx context.Context, otp string) (string, error)
	}
)
