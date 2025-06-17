package services

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"github.com/co1seam/ember-backend-auth/internal/adapters/repository"
	"github.com/co1seam/ember-backend-auth/internal/core/models"
	"github.com/co1seam/ember-backend-auth/internal/ports"
	"math/big"
	"net/smtp"
	"time"
)

const salt = "4d665e8dbe585764403bdc28bf9848ca"

type Authorization struct {
	repo  ports.IAuthRepo
	cache *repository.Redis
	opts  *models.Options
}

func NewAuthorization(repo ports.IAuthRepo, cache *repository.Redis, opts *models.Options) *Authorization {
	return &Authorization{repo: repo, cache: cache, opts: opts}
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

func (a *Authorization) VerifyOTP(ctx context.Context, otp string) (string, error) {
	key, err := a.cache.Redis.Get(ctx, otp).Result()
	if err != nil {
		return "", err
	}

	return key, nil
}

func (a *Authorization) SendOTP(ctx context.Context, email string) error {
	otp, err := a.generateOTP(6)
	if err != nil {
		return err
	}
	subject := "OTP"
	body := fmt.Sprintf("Вы запросили одноразовый OTP код для регистрации.\nВаш OTP код: %s", otp)
	to := []string{
		email,
	}

	message := []byte(
		fmt.Sprintf("From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n",
			a.opts.Config.SMTP.From,
			email,
			subject,
			body),
	)

	if err := smtp.SendMail(a.opts.Config.SMTP.Host+":"+a.opts.Config.SMTP.Port, nil, a.opts.Config.SMTP.From, to, message); err != nil {
		return err
	}

	if err := a.cache.Redis.Set(ctx, otp, email, 15*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}

func (a *Authorization) generateOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)

	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", fmt.Errorf("OTP generate error: %v", err)
		}
		otp[i] = digits[num.Int64()]
	}

	return string(otp), nil
}

func (a *Authorization) generateHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
