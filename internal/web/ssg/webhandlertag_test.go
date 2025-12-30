package ssg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

func newTestWebHandlerWithMockAPI(getResp interface{}, getErr error, postResp interface{}, postErr error, putErr error, deleteErr error) (*WebHandler, *httptest.Server) {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			// Handle specific paths needed by renderContentForm (always succeed)
			if strings.HasPrefix(r.URL.Path, "/ssg/sections") {
				json.NewEncoder(w).Encode(map[string]interface{}{"sections": []interface{}{}})
				return
			}
			if strings.HasPrefix(r.URL.Path, "/auth/users") {
				json.NewEncoder(w).Encode(map[string]interface{}{"users": []interface{}{}})
				return
			}
			if strings.HasPrefix(r.URL.Path, "/ssg/tags") && !strings.Contains(r.URL.Path, "search") {
				json.NewEncoder(w).Encode(map[string]interface{}{"tags": []interface{}{}})
				return
			}

			// Apply getErr to main requests
			if getErr != nil {
				http.Error(w, getErr.Error(), http.StatusInternalServerError)
				return
			}

			if getResp != nil {
				json.NewEncoder(w).Encode(getResp)
			} else {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{})
			}
		case http.MethodPost:
			if postErr != nil {
				http.Error(w, postErr.Error(), http.StatusInternalServerError)
				return
			}
			if postResp != nil {
				json.NewEncoder(w).Encode(postResp)
			} else {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{})
			}
		case http.MethodPut:
			if putErr != nil {
				http.Error(w, putErr.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{})
		case http.MethodDelete:
			if deleteErr != nil {
				http.Error(w, deleteErr.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{})
		}
	}))

	tm := hm.NewTemplateManager(testAssetsFS, params)
	flash := hm.NewFlashManager(params)
	sessMgr := &mockSessionManager{siteSlug: "test-site"}

	// Create test repo and param manager
	repo := newTestRepo(nil)
	paramMgr := feat.NewParamManager(repo, params)

	var smInterface interface {
		SetUserSession(w http.ResponseWriter, userID uuid.UUID, siteSlug string) error
		GetUserSession(r *http.Request) (userID uuid.UUID, siteSlug string, err error)
		SetSiteSlug(w http.ResponseWriter, r *http.Request, siteSlug string) error
	} = sessMgr

	handler := NewWebHandler(tm, flash, paramMgr, nil, smInterface, params)
	handler.apiClient = hm.NewAPIClient("test-api", func() string { return "" }, mockServer.URL, params)

	return handler, mockServer
}

func TestWebHandlerListTags(t *testing.T) {
	tests := []struct {
		name           string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name: "lists tags successfully",
			getResp: map[string]interface{}{
				"tags": []feat.Tag{
					{Name: "Tag 1"},
					{Name: "Tag 2"},
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails when API returns error",
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/list-tags", nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.ListTags(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ListTags() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerShowTag(t *testing.T) {
	tagID := uuid.New()
	tests := []struct {
		name           string
		queryID        string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name:    "shows tag successfully",
			queryID: tagID.String(),
			getResp: map[string]interface{}{
				"tag": feat.Tag{
					Name: "Test Tag",
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with missing ID",
			queryID:        "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "fails when API returns error",
			queryID:        tagID.String(),
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/show-tag?id="+tt.queryID, nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.ShowTag(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ShowTag() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerCreateTag(t *testing.T) {
	tests := []struct {
		name           string
		formData       url.Values
		postResp       interface{}
		postErr        error
		wantStatusCode int
	}{
		{
			name: "creates tag successfully",
			formData: url.Values{
				"name": []string{"New Tag"},
			},
			postResp: map[string]interface{}{
				"tag": map[string]interface{}{
					"id":   uuid.New().String(),
					"name": "New Tag",
				},
			},
			wantStatusCode: http.StatusSeeOther,
		},
		{
			name: "fails with invalid form data",
			formData: url.Values{
				"name": []string{""},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when API returns error",
			formData: url.Values{
				"name": []string{"New Tag"},
			},
			postErr:        fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(nil, nil, tt.postResp, tt.postErr, nil, nil)
			defer server.Close()

			body := strings.NewReader(tt.formData.Encode())
			req := httptest.NewRequest(http.MethodPost, "/ssg/create-tag", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.CreateTag(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("CreateTag() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerUpdateTag(t *testing.T) {
	tagID := uuid.New()
	tests := []struct {
		name           string
		formData       url.Values
		putErr         error
		wantStatusCode int
	}{
		{
			name: "updates tag successfully",
			formData: url.Values{
				"id":   []string{tagID.String()},
				"name": []string{"Updated Tag"},
			},
			wantStatusCode: http.StatusSeeOther,
		},
		{
			name: "fails with invalid form data",
			formData: url.Values{
				"id":   []string{tagID.String()},
				"name": []string{""},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when API returns error",
			formData: url.Values{
				"id":   []string{tagID.String()},
				"name": []string{"Updated Tag"},
			},
			putErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(nil, nil, nil, nil, tt.putErr, nil)
			defer server.Close()

			body := strings.NewReader(tt.formData.Encode())
			req := httptest.NewRequest(http.MethodPost, "/ssg/update-tag", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.UpdateTag(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("UpdateTag() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerDeleteTag(t *testing.T) {
	tagID := uuid.New()
	tests := []struct {
		name           string
		formData       url.Values
		deleteErr      error
		wantStatusCode int
	}{
		{
			name: "deletes tag successfully",
			formData: url.Values{
				"id": []string{tagID.String()},
			},
			wantStatusCode: http.StatusSeeOther,
		},
		{
			name:           "fails with missing ID",
			formData:       url.Values{},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when API returns error",
			formData: url.Values{
				"id": []string{tagID.String()},
			},
			deleteErr:      fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(nil, nil, nil, nil, nil, tt.deleteErr)
			defer server.Close()

			body := strings.NewReader(tt.formData.Encode())
			req := httptest.NewRequest(http.MethodPost, "/ssg/delete-tag", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.DeleteTag(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("DeleteTag() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerNewTag(t *testing.T) {
	handler, server := newTestWebHandlerWithMockAPI(nil, nil, nil, nil, nil, nil)
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, "/ssg/new-tag", nil)
	ctx := feat.NewContextWithSite("test-site", uuid.New())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.NewTag(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("NewTag() status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestWebHandlerEditTag(t *testing.T) {
	tagID := uuid.New()
	tests := []struct {
		name           string
		queryID        string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name:    "shows edit form successfully",
			queryID: tagID.String(),
			getResp: map[string]interface{}{
				"tag": map[string]interface{}{
					"id":   tagID.String(),
					"name": "Test Tag",
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with missing ID",
			queryID:        "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "fails when API returns error",
			queryID:        tagID.String(),
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/edit-tag?id="+tt.queryID, nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.EditTag(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("EditTag() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}
