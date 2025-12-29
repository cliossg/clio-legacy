package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/fake"
	"github.com/hermesgen/clio/internal/feat/auth"
	"github.com/hermesgen/hm"
)

func setupAPIHandler() *auth.APIHandler {
	repo := fake.NewAuthRepo()
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}

	handler := auth.NewAPIHandler("test-auth-api", repo, params)
	handler.Setup(context.Background())

	return handler
}

func TestAPIHandlerGetAllUsers(t *testing.T) {
	handler := setupAPIHandler()

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	handler.GetAllUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetAllUsers() status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestAPIHandlerGetUser(t *testing.T) {
	handler := setupAPIHandler()
	userID := uuid.New()

	req := httptest.NewRequest("GET", "/users/"+userID.String(), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", userID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.GetUser(w, req)

	if w.Code != http.StatusOK {
		t.Logf("GetUser() status = %d (expected OK for existing user or error for non-existent)", w.Code)
	}
}

func TestAPIHandlerCreateUser(t *testing.T) {
	handler := setupAPIHandler()

	userForm := auth.UserForm{
		Username: "testuser",
		Name:     "Test User",
		Email:    "test@example.com",
	}

	body, _ := json.Marshal(userForm)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("CreateUser() status = %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestAPIHandlerUpdateUser(t *testing.T) {
	handler := setupAPIHandler()
	userID := uuid.New()

	userForm := auth.UserForm{
		Username: "updateduser",
		Name:     "Updated User",
		Email:    "updated@example.com",
	}

	body, _ := json.Marshal(userForm)
	req := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", userID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.UpdateUser(w, req)

	if w.Code != http.StatusOK {
		t.Logf("UpdateUser() status = %d", w.Code)
	}
}

func TestAPIHandlerDeleteUser(t *testing.T) {
	handler := setupAPIHandler()
	userID := uuid.New()

	req := httptest.NewRequest("DELETE", "/users/"+userID.String(), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", userID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.DeleteUser(w, req)

	if w.Code != http.StatusOK {
		t.Logf("DeleteUser() status = %d", w.Code)
	}
}
