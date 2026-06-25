package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gosh/internal/config"
)

func init() {
	config.AppConfig = &config.Config{
		JWT: config.JWTConfig{Secret: "test-secret", ExpireHour: 72},
	}
}

func TestSignAndParse_Success(t *testing.T) {
	token, expiresAt, err := Sign(uint(1), "user", nil)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	assert.True(t, expiresAt.After(time.Now()))

	claims, err := Parse(token)
	require.NoError(t, err)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "user", claims.Role)
	assert.Nil(t, claims.TenantID)
}

func TestSignAndParse_WithTenantID(t *testing.T) {
	tenantID := uint(42)
	token, _, err := Sign(uint(2), "merchant", &tenantID)
	require.NoError(t, err)

	claims, err := Parse(token)
	require.NoError(t, err)
	assert.Equal(t, uint(2), claims.UserID)
	assert.Equal(t, "merchant", claims.Role)
	require.NotNil(t, claims.TenantID)
	assert.Equal(t, uint(42), *claims.TenantID)
}

func TestParse_InvalidToken(t *testing.T) {
	_, err := Parse("invalid.jwt.token")
	assert.Error(t, err)
}

func TestParse_ExpiredToken(t *testing.T) {
	claims := &Claims{
		UserID: 1,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("test-secret"))
	require.NoError(t, err)

	_, err = Parse(signed)
	assert.Error(t, err)
}

func TestParse_WrongSecret(t *testing.T) {
	claims := &Claims{
		UserID: 1,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("wrong-secret"))
	require.NoError(t, err)

	_, err = Parse(signed)
	assert.Error(t, err)
}
