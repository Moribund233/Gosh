package router

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gosh/internal/config"
	"gosh/internal/database"
	"gosh/internal/model"
	"gosh/internal/testutil"
)

func setupRouter(t *testing.T) *gin.Engine {
	t.Helper()
	config.AppConfig = &config.Config{
		Server: config.ServerConfig{Mode: "test"},
		JWT:    config.JWTConfig{Secret: "test-secret", ExpireHour: 72},
		Upload: config.UploadConfig{Dir: "/tmp/test-uploads", MaxSize: 10},
	}
	database.Init(config.DatabaseConfig{Driver: "sqlite", Path: ":memory:"})
	database.DB.AutoMigrate(
		&model.User{},
		&model.Address{},
		&model.Favorite{},
		&model.BrowseHistory{},
		&model.MerchantApplication{},
		&model.Category{},
		&model.Product{},
		&model.ProductSKU{},
		&model.ProductReview{},
		&model.SearchHistory{},
		&model.HotSearch{},
		&model.Banner{},
		&model.BrandStory{},
	)
	log, _ := zap.NewDevelopment()
	return New(log)
}

func TestHealthCheck(t *testing.T) {
	r := setupRouter(t)
	w := testutil.PerformRequest(r, "GET", "/health", "")
	assert.Equal(t, http.StatusOK, w.Code)

	var resp testutil.TestResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "ok", resp.Message)
}

func TestUserRegisterRoute(t *testing.T) {
	r := setupRouter(t)
	w := testutil.PerformRequest(r, "POST", "/api/v1/user/register", `{"phone":"13800138000","password":"pass123","nickname":"test"}`)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestUserLoginRoute(t *testing.T) {
	r := setupRouter(t)
	w1 := testutil.PerformRequest(r, "POST", "/api/v1/user/register", `{"phone":"13800138000","password":"pass123","nickname":"test"}`)
	require.Equal(t, http.StatusCreated, w1.Code)

	w2 := testutil.PerformRequest(r, "POST", "/api/v1/user/login", `{"phone":"13800138000","password":"pass123"}`)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func getAuthToken(t *testing.T, r *gin.Engine) string {
	t.Helper()
	w1 := testutil.PerformRequest(r, "POST", "/api/v1/user/register", `{"phone":"13800138000","password":"pass123","nickname":"test"}`)
	require.Equal(t, http.StatusCreated, w1.Code)

	w2 := testutil.PerformRequest(r, "POST", "/api/v1/user/login", `{"phone":"13800138000","password":"pass123"}`)
	require.Equal(t, http.StatusOK, w2.Code)

	var loginResp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &loginResp)
	return loginResp.Data.Token
}

func TestUserProfileRoute_Unauthenticated(t *testing.T) {
	r := setupRouter(t)
	w := testutil.PerformRequest(r, "GET", "/api/v1/user/profile", "")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserProfileRoute_Authenticated(t *testing.T) {
	r := setupRouter(t)
	token := getAuthToken(t, r)

	w := testutil.PerformRequest(r, "GET", "/api/v1/user/profile", "", map[string]string{"Authorization": "Bearer " + token})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBannerRoute_Public(t *testing.T) {
	r := setupRouter(t)
	w := testutil.PerformRequest(r, "GET", "/api/v1/banners", "")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBrandStoryRoute_Public(t *testing.T) {
	r := setupRouter(t)
	w := testutil.PerformRequest(r, "GET", "/api/v1/brand-story", "")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUploadRoute_Unauthenticated(t *testing.T) {
	r := setupRouter(t)
	w := testutil.PerformRequest(r, "POST", "/api/v1/upload", "")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestNotFoundRoute(t *testing.T) {
	r := setupRouter(t)
	w := testutil.PerformRequest(r, "GET", "/api/v1/nonexistent", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}
