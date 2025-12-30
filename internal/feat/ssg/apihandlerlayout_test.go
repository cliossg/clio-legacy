package ssg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
)

func TestAPIHandlerGetLayout(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		layoutID       string
		wantStatusCode int
	}{
		{
			name: "gets layout successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.layouts[id] = Layout{ID: id, Name: "Test Layout"}
				return id
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with invalid UUID",
			layoutID:       "invalid-uuid",
			setupRepo:      func(m *mockServiceRepo) uuid.UUID { return uuid.Nil },
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when layout not found",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				return uuid.New()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			layoutID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			idStr := tt.layoutID
			if idStr == "" {
				idStr = layoutID.String()
			}

			req := httptest.NewRequest(http.MethodGet, "/ssg/layouts/"+idStr, nil)
			req.SetPathValue("id", idStr)
			w := httptest.NewRecorder()

			handler.GetLayout(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetLayout() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerGetAllLayouts(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo)
		wantStatusCode int
	}{
		{
			name: "gets all layouts successfully",
			setupRepo: func(m *mockServiceRepo) {
				m.layouts[uuid.New()] = Layout{ID: uuid.New(), Name: "Layout 1"}
				m.layouts[uuid.New()] = Layout{ID: uuid.New(), Name: "Layout 2"}
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "returns empty list when no layouts",
			setupRepo:      func(m *mockServiceRepo) {},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) {
				m.getLayoutErr = fmt.Errorf("db error")
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			req := httptest.NewRequest(http.MethodGet, "/ssg/layouts", nil)
			w := httptest.NewRecorder()

			handler.GetAllLayouts(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetAllLayouts() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerCreateLayout(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupRepo      func(*mockServiceRepo)
		wantStatusCode int
	}{
		{
			name: "creates layout successfully",
			requestBody: map[string]string{
				"name": "Test Layout",
				"code": "test-code",
			},
			setupRepo:      func(m *mockServiceRepo) {},
			wantStatusCode: http.StatusCreated,
		},
		{
			name:           "fails with invalid JSON",
			requestBody:    "invalid json",
			setupRepo:      func(m *mockServiceRepo) {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when service returns error",
			requestBody: map[string]string{
				"name": "Test Layout",
			},
			setupRepo: func(m *mockServiceRepo) {
				m.createLayoutErr = fmt.Errorf("db error")
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				var err error
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatal(err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/ssg/layouts", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateLayout(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("CreateLayout() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerUpdateLayout(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		requestBody    interface{}
		layoutID       string
		wantStatusCode int
	}{
		{
			name: "updates layout successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.layouts[id] = Layout{ID: id, Name: "Old Name"}
				return id
			},
			requestBody: map[string]string{
				"name": "New Name",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with invalid UUID",
			layoutID:       "invalid-uuid",
			setupRepo:      func(m *mockServiceRepo) uuid.UUID { return uuid.Nil },
			requestBody:    map[string]string{"name": "Test"},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "fails with invalid JSON body",
			setupRepo:      func(m *mockServiceRepo) uuid.UUID { return uuid.New() },
			requestBody:    "invalid json",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				m.updateLayoutErr = fmt.Errorf("db error")
				return uuid.New()
			},
			requestBody: map[string]string{
				"name": "Test",
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			layoutID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				var err error
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatal(err)
				}
			}

			idStr := tt.layoutID
			if idStr == "" {
				idStr = layoutID.String()
			}

			req := httptest.NewRequest(http.MethodPut, "/ssg/layouts/"+idStr, bytes.NewReader(body))
			req.SetPathValue("id", idStr)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.UpdateLayout(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("UpdateLayout() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerDeleteLayout(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		wantStatusCode int
	}{
		{
			name: "deletes layout successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.layouts[id] = Layout{ID: id, Name: "Test"}
				return id
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				m.deleteLayoutErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			layoutID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			req := httptest.NewRequest(http.MethodDelete, "/ssg/layouts/"+layoutID.String(), nil)
			req.SetPathValue("id", layoutID.String())
			w := httptest.NewRecorder()

			handler.DeleteLayout(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("DeleteLayout() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}
