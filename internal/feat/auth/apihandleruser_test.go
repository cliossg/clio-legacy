package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

func setupAPIHandlerWithRepo(repo auth.Repo) *auth.APIHandler {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}

	handler := auth.NewAPIHandler("test-auth-api", repo, params)
	handler.Setup(context.Background())

	return handler
}

func TestAPIHandlerGetAllUsers(t *testing.T) {
	tests := []struct {
		name           string
		wantStatusCode int
	}{
		{
			name:           "gets all users successfully",
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := setupAPIHandler()
			req := httptest.NewRequest("GET", "/users", nil)
			w := httptest.NewRecorder()

			handler.GetAllUsers(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetAllUsers() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerGetUser(t *testing.T) {
	repo := fake.NewAuthRepo()
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}
	svc := auth.NewService(repo, params)

	// Create a user first for success test
	user := auth.NewUser("testuser", "Test User", "test@example.com")
	err := svc.CreateUser(context.Background(), &user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name           string
		userID         string
		wantStatusCode int
	}{
		{
			name:           "gets user successfully",
			userID:         user.ID.String(),
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with invalid UUID",
			userID:         "invalid-uuid",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "fails with non existent user",
			userID:         uuid.New().String(),
			wantStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := setupAPIHandlerWithRepo(repo)
			req := httptest.NewRequest("GET", "/users/"+tt.userID, nil)
			req.SetPathValue("id", tt.userID)
			w := httptest.NewRecorder()

			handler.GetUser(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetUser() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		wantStatusCode int
	}{
		{
			name: "creates user successfully",
			requestBody: auth.UserForm{
				Username: "newuser",
				Name:     "New User",
				Email:    "new@example.com",
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name:           "fails with invalid JSON",
			requestBody:    "invalid json",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "creates user with minimal fields",
			requestBody: auth.UserForm{
				Username: "minimal",
				Name:     "Minimal",
				Email:    "minimal@example.com",
			},
			wantStatusCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := setupAPIHandler()

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("CreateUser() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerUpdateUser(t *testing.T) {
	repo := fake.NewAuthRepo()
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}
	svc := auth.NewService(repo, params)

	// Create a user first
	user := auth.NewUser("existinguser", "Existing User", "existing@example.com")
	err := svc.CreateUser(context.Background(), &user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name           string
		userID         string
		requestBody    interface{}
		wantStatusCode int
	}{
		{
			name:   "updates user successfully",
			userID: user.ID.String(),
			requestBody: auth.UserForm{
				Username: "updateduser",
				Name:     "Updated User",
				Email:    "updated@example.com",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with invalid UUID",
			userID:         "invalid-uuid",
			requestBody:    auth.UserForm{Username: "test", Name: "Test", Email: "test@example.com"},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "fails with invalid JSON",
			userID:         user.ID.String(),
			requestBody:    "invalid json",
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := setupAPIHandlerWithRepo(repo)
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("PUT", "/users/"+tt.userID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.SetPathValue("id", tt.userID)
			w := httptest.NewRecorder()

			handler.UpdateUser(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("UpdateUser() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerDeleteUser(t *testing.T) {
	repo := fake.NewAuthRepo()
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}
	svc := auth.NewService(repo, params)

	// Create a user first
	user := auth.NewUser("deleteuser", "Delete User", "delete@example.com")
	err := svc.CreateUser(context.Background(), &user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name           string
		userID         string
		wantStatusCode int
	}{
		{
			name:           "deletes user successfully",
			userID:         user.ID.String(),
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with invalid UUID",
			userID:         "invalid-uuid",
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := setupAPIHandlerWithRepo(repo)
			req := httptest.NewRequest("DELETE", "/users/"+tt.userID, nil)
			req.SetPathValue("id", tt.userID)
			w := httptest.NewRecorder()

			handler.DeleteUser(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("DeleteUser() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}
