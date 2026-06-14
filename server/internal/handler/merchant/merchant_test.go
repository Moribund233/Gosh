package merchant

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gosh/internal/config"
	"gosh/internal/database"
	"gosh/internal/model"
	"gosh/internal/testutil"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	config.AppConfig = &config.Config{
		Server: config.ServerConfig{Mode: "test"},
		JWT:    config.JWTConfig{Secret: "test-secret", ExpireHour: 72},
	}
	cfg := config.DatabaseConfig{Driver: "sqlite", Path: ":memory:"}
	err := database.Init(cfg)
	require.NoError(t, err)
	err = database.DB.AutoMigrate(&model.MerchantApplication{}, &model.User{})
	require.NoError(t, err)
}

func TestApplyMerchant_Success(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	body := `{"shop_name":"暖食集旗舰店","shop_desc":"优质食品","contact_name":"李四","contact_phone":"13800138000"}`
	c, w := testutil.NewGinContext("POST", "/api/v1/merchant/apply", body)
	c.Set("user_id", uint(1))
	h.Apply(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "暖食集旗舰店", resp.Data["shop_name"])
	assert.Equal(t, "pending", resp.Data["status"])
}

func TestApplyMerchant_Duplicate(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	body := `{"shop_name":"暖食集旗舰店","contact_name":"李四","contact_phone":"13800138000"}`
	c1, _ := testutil.NewGinContext("POST", "/api/v1/merchant/apply", body)
	c1.Set("user_id", uint(1))
	h.Apply(c1)

	c2, w2 := testutil.NewGinContext("POST", "/api/v1/merchant/apply", body)
	c2.Set("user_id", uint(1))
	h.Apply(c2)

	assert.Equal(t, http.StatusConflict, w2.Code)
}

func TestMyApplication(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	body := `{"shop_name":"暖食集旗舰店","contact_name":"李四","contact_phone":"13800138000"}`
	c1, _ := testutil.NewGinContext("POST", "/api/v1/merchant/apply", body)
	c1.Set("user_id", uint(1))
	h.Apply(c1)

	c2, w2 := testutil.NewGinContext("GET", "/api/v1/merchant/application", "")
	c2.Set("user_id", uint(1))
	h.MyApplication(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "pending", resp.Data["status"])
}

func TestMyApplication_NotFound(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	c, w := testutil.NewGinContext("GET", "/api/v1/merchant/application", "")
	c.Set("user_id", uint(99))
	h.MyApplication(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestReviewApprove(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	// Create the user that will become a merchant
	database.DB.Create(&model.User{Phone: "13800138000", Password: "hash", Nickname: "test_user", Role: model.RoleUser, Status: model.StatusActive})

	body := `{"shop_name":"暖食集旗舰店","contact_name":"李四","contact_phone":"13800138000"}`
	c1, w1 := testutil.NewGinContext("POST", "/api/v1/merchant/apply", body)
	c1.Set("user_id", uint(1))
	h.Apply(c1)

	var createResp struct {
		Data struct {
			ID float64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(w1.Body.Bytes(), &createResp)

	reviewBody := fmt.Sprintf(`{"application_id":%.0f,"action":"approve"}`, createResp.Data.ID)
	c2, w2 := testutil.NewGinContext("POST", "/api/v1/admin/merchant/review", reviewBody)
	h.Review(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "approved", resp.Data["status"])
}

func TestReviewReject(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	body := `{"shop_name":"暖食集旗舰店","contact_name":"李四","contact_phone":"13800138000"}`
	c1, w1 := testutil.NewGinContext("POST", "/api/v1/merchant/apply", body)
	c1.Set("user_id", uint(1))
	h.Apply(c1)

	var createResp struct {
		Data struct {
			ID float64 `json:"id"`
		} `json:"data"`
	}
	json.Unmarshal(w1.Body.Bytes(), &createResp)

	reviewBody := fmt.Sprintf(`{"application_id":%.0f,"action":"reject","remark":"资质不足"}`, createResp.Data.ID)
	c2, w2 := testutil.NewGinContext("POST", "/api/v1/admin/merchant/review", reviewBody)
	h.Review(c2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "rejected", resp.Data["status"])
	assert.Equal(t, "资质不足", resp.Data["remark"])
}

func TestReview_NotFound(t *testing.T) {
	setupTestDB(t)
	h := NewHandler()

	reviewBody := `{"application_id":999,"action":"approve"}`
	c, w := testutil.NewGinContext("POST", "/api/v1/admin/merchant/review", reviewBody)
	h.Review(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
