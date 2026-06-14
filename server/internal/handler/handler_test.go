package handler

import (
	"net/http"
	"testing"

	"gosh/internal/testutil"
)

func TestRegisterRoute(t *testing.T) {
	r := setupTestEngine(t)
	w := testutil.PerformRequest(r, "POST", "/api/v1/user/register", `{"phone":"13800138000","password":"pass123","nickname":"test"}`)
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Errorf("unexpected status: %d", w.Code)
	}
}

func TestLoginRoute(t *testing.T) {
	r := setupTestEngine(t)

	w1 := testutil.PerformRequest(r, "POST", "/api/v1/user/register", `{"phone":"13800138000","password":"pass123","nickname":"test"}`)
	if w1.Code != http.StatusOK && w1.Code != http.StatusCreated {
		t.Fatalf("register failed: %d", w1.Code)
	}

	w2 := testutil.PerformRequest(r, "POST", "/api/v1/user/login", `{"phone":"13800138000","password":"pass123"}`)
	if w2.Code != http.StatusOK {
		t.Errorf("unexpected status: %d", w2.Code)
	}
}
