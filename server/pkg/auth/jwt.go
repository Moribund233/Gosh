package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gosh/internal/config"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Role     string `json:"role"`
	TenantID *uint  `json:"tenant_id,omitempty"`
	jwt.RegisteredClaims
}

func Sign(userID uint, role string, tenantID *uint) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(config.AppConfig.JWT.ExpireHour) * time.Hour)
	claims := &Claims{
		UserID:   userID,
		Role:     role,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.AppConfig.JWT.Secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

func Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
