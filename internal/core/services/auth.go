package services

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/co1seam/ember-backend-auth/internal/core/models"
	"github.com/co1seam/ember-backend-auth/internal/ports"
)

const salt = "4d665e8dbe585764403bdc28bf9848ca"

type Authorization struct {
	repo ports.IAuthRepo
	opts *models.Options
}

func NewAuthorization(repo ports.IAuthRepo, opts *models.Options) *Authorization {
	return &Authorization{repo: repo, opts: opts}
}

func (a *Authorization) Create(ctx context.Context, entity ...interface{}) (interface{}, error) {
	user := entity[0].(models.SignUpRequest)
	user.Password = a.generateHash(user.Password)

	return a.repo.Create(ctx, user)
}

func (a *Authorization) Read(ctx context.Context, entity ...interface{}) (interface{}, error) {
	userModel := entity[0].(models.SignInRequest)
	userModel.Password = a.generateHash(userModel.Password)

	id, err := a.repo.Read(ctx, userModel)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (a *Authorization) Update(ctx context.Context, entity ...interface{}) (interface{}, error) {
	return nil, nil
}

func (a *Authorization) Delete(ctx context.Context, filter ...interface{}) (interface{}, error) {
	return nil, nil
}

func (a *Authorization) generateHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
