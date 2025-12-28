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

func TestAPIHandlerGetAllSections(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo)
		wantStatusCode int
	}{
		{
			name: "gets all sections successfully",
			setupRepo: func(m *mockServiceRepo) {
				m.sections[uuid.New()] = Section{ID: uuid.New(), Name: "Section 1"}
				m.sections[uuid.New()] = Section{ID: uuid.New(), Name: "Section 2"}
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "returns empty list when no sections",
			setupRepo:      func(m *mockServiceRepo) {},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) {
				m.getSectionErr = fmt.Errorf("db error")
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

			req := httptest.NewRequest(http.MethodGet, "/ssg/sections", nil)
			w := httptest.NewRecorder()

			handler.GetAllSections(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetAllSections() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerGetSection(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		sectionID      string
		wantStatusCode int
	}{
		{
			name: "gets section successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.sections[id] = Section{ID: id, Name: "Test Section"}
				return id
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with invalid UUID",
			sectionID:      "invalid-uuid",
			setupRepo:      func(m *mockServiceRepo) uuid.UUID { return uuid.Nil },
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when section not found",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				return uuid.New()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			sectionID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			idStr := tt.sectionID
			if idStr == "" {
				idStr = sectionID.String()
			}

			req := httptest.NewRequest(http.MethodGet, "/ssg/sections/"+idStr, nil)
			req.SetPathValue("id", idStr)
			w := httptest.NewRecorder()

			handler.GetSection(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetSection() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerCreateSection(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupRepo      func(*mockServiceRepo)
		wantStatusCode int
	}{
		{
			name: "creates section successfully",
			requestBody: map[string]interface{}{
				"name":        "Test Section",
				"description": "Test description",
				"path":        "test-section",
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
				"name": "Test Section",
			},
			setupRepo: func(m *mockServiceRepo) {
				m.createSectionErr = fmt.Errorf("db error")
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

			req := httptest.NewRequest(http.MethodPost, "/ssg/sections", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateSection(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("CreateSection() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerUpdateSection(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		requestBody    interface{}
		wantStatusCode int
	}{
		{
			name: "updates section successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.sections[id] = Section{ID: id, Name: "Old Name"}
				return id
			},
			requestBody: map[string]string{
				"name": "New Name",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				m.updateSectionErr = fmt.Errorf("db error")
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
			sectionID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPut, "/ssg/sections/"+sectionID.String(), bytes.NewReader(body))
			req.SetPathValue("id", sectionID.String())
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.UpdateSection(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("UpdateSection() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerDeleteSection(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		wantStatusCode int
	}{
		{
			name: "deletes section successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.sections[id] = Section{ID: id, Name: "Test"}
				return id
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				m.deleteSectionErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			sectionID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			req := httptest.NewRequest(http.MethodDelete, "/ssg/sections/"+sectionID.String(), nil)
			req.SetPathValue("id", sectionID.String())
			w := httptest.NewRecorder()

			handler.DeleteSection(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("DeleteSection() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

// Note: TestAPIHandlerUploadSectionImage and TestAPIHandlerDeleteSectionImage are skipped
// because they require ImageManager which is a concrete struct and difficult to mock.
// These handlers will be tested via integration tests or when ImageManager is refactored to use an interface.
