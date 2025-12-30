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

func TestAPIHandlerCreateParam(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupRepo      func(*mockServiceRepo)
		wantStatusCode int
	}{
		{
			name: "creates param successfully",
			requestBody: map[string]string{
				"name":  "test-param",
				"value": "test-value",
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
				"name":  "test-param",
				"value": "test-value",
			},
			setupRepo: func(m *mockServiceRepo) {
				m.createParamErr = fmt.Errorf("db error")
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

			req := httptest.NewRequest(http.MethodPost, "/ssg/params", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateParam(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("CreateParam() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerGetParam(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		paramID        string
		wantStatusCode int
	}{
		{
			name: "gets param successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.params[id] = Param{ID: id, Name: "test-param"}
				return id
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with invalid UUID",
			paramID:        "invalid-uuid",
			setupRepo:      func(m *mockServiceRepo) uuid.UUID { return uuid.Nil },
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when param not found",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				return uuid.New()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			paramID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			idStr := tt.paramID
			if idStr == "" {
				idStr = paramID.String()
			}

			req := httptest.NewRequest(http.MethodGet, "/ssg/params/"+idStr, nil)
			req.SetPathValue("id", idStr)
			w := httptest.NewRecorder()

			handler.GetParam(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetParam() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerGetParamByName(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo)
		paramName      string
		wantStatusCode int
	}{
		{
			name: "gets param by name successfully",
			setupRepo: func(m *mockServiceRepo) {
				m.paramsByName["test-param"] = Param{Name: "test-param"}
			},
			paramName:      "test-param",
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when param not found",
			setupRepo: func(m *mockServiceRepo) {
				m.getParamErr = fmt.Errorf("not found")
			},
			paramName:      "missing",
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

			req := httptest.NewRequest(http.MethodGet, "/ssg/params/name/"+tt.paramName, nil)
			req.SetPathValue("name", tt.paramName)
			w := httptest.NewRecorder()

			handler.GetParamByName(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetParamByName() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerGetParamByRefKey(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo)
		refKey         string
		wantStatusCode int
	}{
		{
			name: "gets param by ref key successfully",
			setupRepo: func(m *mockServiceRepo) {
				m.paramsByRef["test-ref"] = Param{RefKey: "test-ref"}
			},
			refKey:         "test-ref",
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when param not found",
			setupRepo: func(m *mockServiceRepo) {
				m.getParamErr = fmt.Errorf("not found")
			},
			refKey:         "missing",
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

			req := httptest.NewRequest(http.MethodGet, "/ssg/params/ref/"+tt.refKey, nil)
			req.SetPathValue("ref_key", tt.refKey)
			w := httptest.NewRecorder()

			handler.GetParamByRefKey(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetParamByRefKey() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerListParams(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo)
		wantStatusCode int
	}{
		{
			name: "lists params successfully",
			setupRepo: func(m *mockServiceRepo) {
				m.params[uuid.New()] = Param{Name: "param1"}
				m.params[uuid.New()] = Param{Name: "param2"}
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "returns empty list when no params",
			setupRepo:      func(m *mockServiceRepo) {},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) {
				m.getParamErr = fmt.Errorf("db error")
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

			req := httptest.NewRequest(http.MethodGet, "/ssg/params", nil)
			w := httptest.NewRecorder()

			handler.ListParams(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ListParams() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerUpdateParam(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		requestBody    interface{}
		paramID        string
		wantStatusCode int
	}{
		{
			name: "updates non-system param successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.params[id] = Param{ID: id, Name: "test", System: 0}
				return id
			},
			requestBody: map[string]interface{}{
				"name":  "updated-name",
				"value": "updated-value",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:    "fails with invalid UUID",
			paramID: "invalid-uuid",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				return uuid.Nil
			},
			requestBody: map[string]interface{}{
				"name": "test",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails with invalid JSON body",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.params[id] = Param{ID: id, Name: "test", System: 0}
				return id
			},
			requestBody:    "invalid json",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when cannot get existing param",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.getParamErr = fmt.Errorf("param not found")
				return id
			},
			requestBody: map[string]interface{}{
				"name": "test",
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "fails when trying to change system param name",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.params[id] = Param{ID: id, Name: "system-param", RefKey: "sys.key", Description: "Desc", System: 1}
				return id
			},
			requestBody: map[string]interface{}{
				"name":   "different-name",
				"refKey": "sys.key",
				"value":  "value",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when trying to change system param refkey",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.params[id] = Param{ID: id, Name: "system-param", RefKey: "sys.key", Description: "Desc", System: 1}
				return id
			},
			requestBody: map[string]interface{}{
				"name":    "system-param",
				"ref_key": "different.key",
				"value":   "value",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when trying to change system param description",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.params[id] = Param{ID: id, Name: "system-param", RefKey: "sys.key", Description: "Original", System: 1}
				return id
			},
			requestBody: map[string]interface{}{
				"name":        "system-param",
				"ref_key":     "sys.key",
				"description": "Different description",
				"value":       "value",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "updates system param value successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.params[id] = Param{ID: id, Name: "system-param", RefKey: "sys.key", Description: "Desc", System: 1}
				return id
			},
			requestBody: map[string]interface{}{
				"name":        "system-param",
				"ref_key":     "sys.key",
				"description": "Desc",
				"value":       "new-value",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.params[id] = Param{ID: id, Name: "test", System: 0}
				m.updateParamErr = fmt.Errorf("db error")
				return id
			},
			requestBody: map[string]interface{}{
				"name": "test",
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			paramID := tt.setupRepo(repo)
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

			idStr := tt.paramID
			if idStr == "" {
				idStr = paramID.String()
			}

			req := httptest.NewRequest(http.MethodPut, "/ssg/params/"+idStr, bytes.NewReader(body))
			req.SetPathValue("id", idStr)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.UpdateParam(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("UpdateParam() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerDeleteParam(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		wantStatusCode int
	}{
		{
			name: "deletes non-system param successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.params[id] = Param{ID: id, Name: "test", System: 0}
				return id
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when trying to delete system param",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.params[id] = Param{ID: id, Name: "system-param", System: 1}
				return id
			},
			wantStatusCode: http.StatusForbidden,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.params[id] = Param{ID: id, Name: "test", System: 0}
				m.deleteParamErr = fmt.Errorf("db error")
				return id
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			paramID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			req := httptest.NewRequest(http.MethodDelete, "/ssg/params/"+paramID.String(), nil)
			req.SetPathValue("id", paramID.String())
			w := httptest.NewRecorder()

			handler.DeleteParam(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("DeleteParam() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}
