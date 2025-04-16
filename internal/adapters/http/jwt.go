package http

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func (h *Handler) createJWT(ttl time.Duration, extraClaims jwt.MapClaims) (string, error) {
	now := time.Now()
	baseClaims := jwt.MapClaims{
		"iat": jwt.NewNumericDate(now),
		"exp": jwt.NewNumericDate(now.Add(ttl)),
		"iss": "ember.com",
	}

	for k, v := range extraClaims {
		if k == "iat" || k == "exp" || k == "iss" {
			return "", fmt.Errorf("protected claim %q cannot be overwritten", k)
		}
		baseClaims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, baseClaims)
	jwtToken, err := token.SignedString([]byte(h.opts.Config.Token.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %v", err)
	}
	return jwtToken, nil
}
