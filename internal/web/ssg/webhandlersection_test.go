package ssg

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
)

func TestWebHandlerListSections(t *testing.T) {
	tests := []struct {
		name           string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name: "lists sections successfully",
			getResp: map[string]interface{}{
				"sections": []feat.Section{
					{Name: "Section 1"},
					{Name: "Section 2"},
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

			req := httptest.NewRequest(http.MethodGet, "/ssg/list-sections", nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.ListSections(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ListSections() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerCreateSection(t *testing.T) {
	layoutID := uuid.New()
	tests := []struct {
		name           string
		formData       url.Values
		postResp       interface{}
		postErr        error
		wantStatusCode int
	}{
		{
			name: "creates section successfully",
			formData: url.Values{
				"name":      []string{"New Section"},
				"path":      []string{"/new"},
				"layout_id": []string{layoutID.String()},
			},
			postResp: map[string]interface{}{
				"section": map[string]interface{}{
					"id":        uuid.New().String(),
					"name":      "New Section",
					"path":      "/new",
					"layout_id": layoutID.String(),
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
				"name":      []string{"New Section"},
				"path":      []string{"/new"},
				"layout_id": []string{layoutID.String()},
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
			req := httptest.NewRequest(http.MethodPost, "/ssg/create-section", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.CreateSection(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("CreateSection() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerEditSection(t *testing.T) {
	sectionID := uuid.New()
	layoutID := uuid.New()
	tests := []struct {
		name           string
		queryID        string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name:    "shows edit form successfully",
			queryID: sectionID.String(),
			getResp: map[string]interface{}{
				"section": map[string]interface{}{
					"id":        sectionID.String(),
					"name":      "Test Section",
					"path":      "/test",
					"layout_id": layoutID.String(),
				},
				"layouts": []feat.Layout{},
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
			queryID:        sectionID.String(),
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/edit-section?id="+tt.queryID, nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.EditSection(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("EditSection() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerUpdateSection(t *testing.T) {
	sectionID := uuid.New()
	layoutID := uuid.New()
	tests := []struct {
		name           string
		formData       url.Values
		putErr         error
		wantStatusCode int
	}{
		{
			name: "updates section successfully",
			formData: url.Values{
				"id":        []string{sectionID.String()},
				"name":      []string{"Updated Section"},
				"path":      []string{"/updated"},
				"layout_id": []string{layoutID.String()},
			},
			wantStatusCode: http.StatusSeeOther,
		},
		{
			name: "fails with invalid form data",
			formData: url.Values{
				"id":   []string{sectionID.String()},
				"name": []string{""},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when API returns error",
			formData: url.Values{
				"id":        []string{sectionID.String()},
				"name":      []string{"Updated Section"},
				"path":      []string{"/updated"},
				"layout_id": []string{layoutID.String()},
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
			req := httptest.NewRequest(http.MethodPost, "/ssg/update-section", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.UpdateSection(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("UpdateSection() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerDeleteSection(t *testing.T) {
	sectionID := uuid.New()
	tests := []struct {
		name           string
		formData       url.Values
		deleteErr      error
		wantStatusCode int
	}{
		{
			name: "deletes section successfully",
			formData: url.Values{
				"id": []string{sectionID.String()},
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
				"id": []string{sectionID.String()},
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
			req := httptest.NewRequest(http.MethodPost, "/ssg/delete-section", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.DeleteSection(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("DeleteSection() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerNewSection(t *testing.T) {
	handler, server := newTestWebHandlerWithMockAPI(nil, nil, nil, nil, nil, nil)
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, "/ssg/new-section", nil)
	ctx := feat.NewContextWithSite("test-site", uuid.New())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.NewSection(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("NewSection() status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestWebHandlerShowSection(t *testing.T) {
	sectionID := uuid.New()
	layoutID := uuid.New()
	tests := []struct {
		name           string
		queryID        string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name:    "shows section successfully",
			queryID: sectionID.String(),
			getResp: map[string]interface{}{
				"section": map[string]interface{}{
					"id":        sectionID.String(),
					"name":      "Test Section",
					"path":      "/test",
					"layout_id": layoutID.String(),
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
			queryID:        sectionID.String(),
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/show-section?id="+tt.queryID, nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.ShowSection(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ShowSection() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}
