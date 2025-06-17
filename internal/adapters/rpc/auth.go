package rpc

import (
	"context"
	"fmt"
	authv1 "github.com/co1seam/ember-backend-api-contracts/gen/go/auth"
	"github.com/co1seam/ember-backend-auth/config"
	"github.com/co1seam/ember-backend-auth/internal/core/models"
	"github.com/co1seam/ember-backend-auth/internal/ports"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

type Authorization struct {
	authv1.UnimplementedAuthServer
	service ports.IAuthService
	opts    *models.Options
}

func NewAuthorization(service ports.IAuthService, opts *models.Options) *Authorization {
	return &Authorization{
		service: service,
		opts:    opts,
	}
}

func (a *Authorization) SendOTP(ctx context.Context, req *authv1.SendOTPRequest) (*authv1.SendOTPResponse, error) {
	otp := models.SendOtpRequest{
		Email: req.Email,
	}

	err := a.service.SendOTP(ctx, otp.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.SendOTPResponse{Success: true}, nil
}

func (a *Authorization) VerifyOTP(ctx context.Context, req *authv1.VerifyOTPRequest) (*authv1.VerifyOTPResponse, error) {
	user := models.VerifyOtpRequest{
		OTP: req.Otp,
	}

	fmt.Println(req.Otp)
	fmt.Println(user.OTP)

	email, err := a.service.VerifyOTP(ctx, user.OTP)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.VerifyOTPResponse{Email: email}, nil
}

func (a *Authorization) SignUp(ctx context.Context, req *authv1.SignUpRequest) (*authv1.SignUpResponse, error) {
	user := models.SignUpRequest{
		Name:     req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	id, err := a.service.Create(ctx, user)
	if err != nil {
		return &authv1.SignUpResponse{AccessToken: "", RefreshToken: ""}, status.Error(codes.Internal, err.Error())
	}

	tokens, err := createTokens(id.(int), &a.opts.Config.Token)
	if err != nil {
		return &authv1.SignUpResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &authv1.SignUpResponse{AccessToken: tokens[1], RefreshToken: tokens[0]}, nil
}

func (a *Authorization) SignIn(ctx context.Context, req *authv1.SignInRequest) (*authv1.SignInResponse, error) {
	user := models.SignInRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	id, err := a.service.Read(ctx, user)
	if err != nil {
		return &authv1.SignInResponse{}, status.Error(codes.Internal, err.Error())
	}
	tokens, err := createTokens(id.(int), &a.opts.Config.Token)
	if err != nil {
		return &authv1.SignInResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &authv1.SignInResponse{AccessToken: tokens[1], RefreshToken: tokens[0]}, nil
}

func (a *Authorization) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	claims, err := verifyJWT(req.RefreshToken, &a.opts.Config.Token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, status.Error(codes.Internal, "invalid refresh token")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid subject")
	}

	userID, err := strconv.Atoi(sub)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tokens, err := createTokens(userID, &a.opts.Config.Token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.RefreshTokenResponse{AccessToken: tokens[1], RefreshToken: tokens[0]}, nil
}

func (a *Authorization) ValidateToken(_ context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	claims, err := verifyJWT(req.AccessToken, &a.opts.Config.Token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return nil, status.Error(codes.Internal, "invalid token")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, status.Error(codes.Internal, "invalid subject")
	}

	return &authv1.ValidateTokenResponse{Subject: sub}, nil
}

func createTokens(userID int, cfg *config.Token) ([]string, error) {
	refreshCookie, err := createJWT(time.Hour*72, cfg, jwt.MapClaims{
		"sub":  fmt.Sprint(userID),
		"type": "refresh",
		"aud":  "admin",
	})
	if err != nil {
		return []string{}, err
	}

	accessToken, err := createJWT(time.Minute*15, cfg, jwt.MapClaims{
		"sub":  fmt.Sprint(userID),
		"type": "access",
		"aud":  "admin",
	})
	if err != nil {
		return []string{}, err
	}

	return []string{refreshCookie, accessToken}, nil
}
