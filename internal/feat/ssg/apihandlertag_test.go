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

func TestAPIHandlerCreateTag(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupRepo      func(*mockServiceRepo)
		wantStatusCode int
	}{
		{
			name: "creates tag successfully",
			requestBody: map[string]string{
				"name": "test-tag",
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
				"name": "test-tag",
			},
			setupRepo: func(m *mockServiceRepo) {
				m.createTagErr = fmt.Errorf("db error")
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

			req := httptest.NewRequest(http.MethodPost, "/ssg/tags", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateTag(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("CreateTag() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerGetTag(t *testing.T) {
	tests := []struct {
		name           string
		tagID          string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		wantStatusCode int
	}{
		{
			name: "gets tag successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.tags[id] = Tag{ID: id, Name: "test-tag"}
				return id
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with invalid UUID",
			tagID:          "invalid-uuid",
			setupRepo:      func(m *mockServiceRepo) uuid.UUID { return uuid.Nil },
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:  "fails when tag not found",
			tagID: uuid.New().String(),
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				return uuid.New()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tagID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			idStr := tt.tagID
			if idStr == "" {
				idStr = tagID.String()
			}

			req := httptest.NewRequest(http.MethodGet, "/ssg/tags/"+idStr, nil)
			req.SetPathValue("id", idStr)
			w := httptest.NewRecorder()

			handler.GetTag(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetTag() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerGetTagByName(t *testing.T) {
	tests := []struct {
		name           string
		tagName        string
		setupRepo      func(*mockServiceRepo)
		wantStatusCode int
	}{
		{
			name:    "gets tag by name successfully",
			tagName: "test-tag",
			setupRepo: func(m *mockServiceRepo) {
				id := uuid.New()
				m.tagsByName["test-tag"] = Tag{ID: id, Name: "test-tag"}
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:    "fails when tag not found",
			tagName: "nonexistent",
			setupRepo: func(m *mockServiceRepo) {
				m.getTagErr = fmt.Errorf("tag not found")
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "fails with empty name",
			tagName:        "",
			setupRepo:      func(m *mockServiceRepo) {},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			req := httptest.NewRequest(http.MethodGet, "/ssg/tags/by-name/"+tt.tagName, nil)
			req.SetPathValue("name", tt.tagName)
			w := httptest.NewRecorder()

			handler.GetTagByName(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetTagByName() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerGetAllTags(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo)
		wantStatusCode int
	}{
		{
			name: "gets all tags successfully",
			setupRepo: func(m *mockServiceRepo) {
				m.tags[uuid.New()] = Tag{ID: uuid.New(), Name: "tag1"}
				m.tags[uuid.New()] = Tag{ID: uuid.New(), Name: "tag2"}
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "returns empty list when no tags",
			setupRepo:      func(m *mockServiceRepo) {},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) {
				m.getTagErr = fmt.Errorf("db error")
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

			req := httptest.NewRequest(http.MethodGet, "/ssg/tags", nil)
			w := httptest.NewRecorder()

			handler.GetAllTags(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetAllTags() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerUpdateTag(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		requestBody    interface{}
		wantStatusCode int
	}{
		{
			name: "updates tag successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.tags[id] = Tag{ID: id, Name: "old-name"}
				return id
			},
			requestBody: func(id uuid.UUID) map[string]interface{} {
				return map[string]interface{}{
					"id":   id.String(),
					"name": "new-name",
				}
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				m.updateTagErr = fmt.Errorf("db error")
				return uuid.New()
			},
			requestBody: map[string]string{
				"name": "test",
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tagID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			var bodyData interface{}
			if fn, ok := tt.requestBody.(func(uuid.UUID) map[string]interface{}); ok {
				bodyData = fn(tagID)
			} else {
				bodyData = tt.requestBody
			}

			body, err := json.Marshal(bodyData)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPut, "/ssg/tags/"+tagID.String(), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.SetPathValue("id", tagID.String())
			w := httptest.NewRecorder()

			handler.UpdateTag(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("UpdateTag() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerDeleteTag(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		tagID          string
		wantStatusCode int
	}{
		{
			name: "deletes tag successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.tags[id] = Tag{ID: id, Name: "test-tag"}
				return id
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with invalid UUID",
			tagID:          "invalid-uuid",
			setupRepo:      func(m *mockServiceRepo) uuid.UUID { return uuid.Nil },
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				m.deleteTagErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tagID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			idStr := tt.tagID
			if idStr == "" {
				idStr = tagID.String()
			}

			req := httptest.NewRequest(http.MethodDelete, "/ssg/tags/"+idStr, nil)
			req.SetPathValue("id", idStr)
			w := httptest.NewRecorder()

			handler.DeleteTag(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("DeleteTag() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}
