package middleware

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gosh/internal/config"
	"gosh/internal/testutil"
	"gosh/pkg/auth"
)

func TestCORS(t *testing.T) {
	r := gin.New()
	r.Use(CORS())
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := testutil.PerformRequest(r, "OPTIONS", "/test", "")
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_GET(t *testing.T) {
	r := gin.New()
	r.Use(CORS())
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := testutil.PerformRequest(r, "GET", "/test", "")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRecovery(t *testing.T) {
	log, _ := zap.NewDevelopment()
	r := gin.New()
	r.Use(Recovery(log))
	r.GET("/panic", func(c *gin.Context) { panic("test panic") })

	w := testutil.PerformRequest(r, "GET", "/panic", "")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAuth_NoToken(t *testing.T) {
	r := gin.New()
	r.Use(Auth())
	r.GET("/protected", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := testutil.PerformRequest(r, "GET", "/protected", "")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuth_EmptyBearer(t *testing.T) {
	r := gin.New()
	r.Use(Auth())
	r.GET("/protected", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := testutil.PerformRequest(r, "GET", "/protected", "", map[string]string{"Authorization": "Bearer "})
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuth_InvalidScheme(t *testing.T) {
	r := gin.New()
	r.Use(Auth())
	r.GET("/protected", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := testutil.PerformRequest(r, "GET", "/protected", "", map[string]string{"Authorization": "Basic token"})
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuth_ValidToken(t *testing.T) {
	setupConfig()
	token, _, err := auth.Sign(1, "user", nil)
	assert.NoError(t, err)

	r := gin.New()
	r.Use(Auth())
	r.GET("/protected", func(c *gin.Context) {
		uid, _ := c.Get("user_id")
		role, _ := c.Get("role")
		assert.Equal(t, uint(1), uid)
		assert.Equal(t, "user", role)
		c.Status(http.StatusOK)
	})

	w := testutil.PerformRequest(r, "GET", "/protected", "", map[string]string{"Authorization": "Bearer " + token})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuth_ExpiredToken(t *testing.T) {
	config.AppConfig = &config.Config{
		JWT: config.JWTConfig{Secret: "test-secret", ExpireHour: -1},
	}
	token, _, err := auth.Sign(1, "user", nil)
	assert.NoError(t, err)

	config.AppConfig.JWT.ExpireHour = 72

	r := gin.New()
	r.Use(Auth())
	r.GET("/protected", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := testutil.PerformRequest(r, "GET", "/protected", "", map[string]string{"Authorization": "Bearer " + token})
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequireRole_Allowed(t *testing.T) {
	setupConfig()
	token, _, err := auth.Sign(1, "admin", nil)
	assert.NoError(t, err)

	r := gin.New()
	r.Use(Auth())
	r.Use(RequireRole("admin"))
	r.GET("/admin", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := testutil.PerformRequest(r, "GET", "/admin", "", map[string]string{"Authorization": "Bearer " + token})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_Forbidden(t *testing.T) {
	setupConfig()
	token, _, err := auth.Sign(1, "user", nil)
	assert.NoError(t, err)

	r := gin.New()
	r.Use(Auth())
	r.Use(RequireRole("admin"))
	r.GET("/admin", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := testutil.PerformRequest(r, "GET", "/admin", "", map[string]string{"Authorization": "Bearer " + token})
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func setupConfig() {
	config.AppConfig = &config.Config{
		Server: config.ServerConfig{Mode: "test"},
		JWT:    config.JWTConfig{Secret: "test-secret", ExpireHour: 72},
	}
}
