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

func TestWebHandlerListLayouts(t *testing.T) {
	tests := []struct {
		name           string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name: "lists layouts successfully",
			getResp: map[string]interface{}{
				"layouts": []feat.Layout{
					{Name: "Layout 1"},
					{Name: "Layout 2"},
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

			req := httptest.NewRequest(http.MethodGet, "/ssg/list-layouts", nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.ListLayouts(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ListLayouts() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerCreateLayout(t *testing.T) {
	tests := []struct {
		name           string
		formData       url.Values
		postResp       interface{}
		postErr        error
		wantStatusCode int
	}{
		{
			name: "creates layout successfully",
			formData: url.Values{
				"name": []string{"New Layout"},
				"code": []string{"<html></html>"},
			},
			postResp: map[string]interface{}{
				"layout": map[string]interface{}{
					"id":   uuid.New().String(),
					"name": "New Layout",
					"code": "<html></html>",
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
				"name": []string{"New Layout"},
				"code": []string{"<html></html>"},
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
			req := httptest.NewRequest(http.MethodPost, "/ssg/create-layout", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.CreateLayout(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("CreateLayout() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerEditLayout(t *testing.T) {
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
			queryID: layoutID.String(),
			getResp: map[string]interface{}{
				"layout": map[string]interface{}{
					"id":   layoutID.String(),
					"name": "Test Layout",
					"code": "<html></html>",
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
			queryID:        layoutID.String(),
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/edit-layout?id="+tt.queryID, nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.EditLayout(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("EditLayout() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerUpdateLayout(t *testing.T) {
	layoutID := uuid.New()
	tests := []struct {
		name           string
		formData       url.Values
		putErr         error
		wantStatusCode int
	}{
		{
			name: "updates layout successfully",
			formData: url.Values{
				"id":   []string{layoutID.String()},
				"name": []string{"Updated Layout"},
				"code": []string{"<html>updated</html>"},
			},
			wantStatusCode: http.StatusSeeOther,
		},
		{
			name: "fails with invalid form data",
			formData: url.Values{
				"id":   []string{layoutID.String()},
				"name": []string{""},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when API returns error",
			formData: url.Values{
				"id":   []string{layoutID.String()},
				"name": []string{"Updated Layout"},
				"code": []string{"<html>updated</html>"},
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
			req := httptest.NewRequest(http.MethodPost, "/ssg/update-layout", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.UpdateLayout(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("UpdateLayout() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerDeleteLayout(t *testing.T) {
	layoutID := uuid.New()
	tests := []struct {
		name           string
		formData       url.Values
		deleteErr      error
		wantStatusCode int
	}{
		{
			name: "deletes layout successfully",
			formData: url.Values{
				"id": []string{layoutID.String()},
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
				"id": []string{layoutID.String()},
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
			req := httptest.NewRequest(http.MethodPost, "/ssg/delete-layout", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.DeleteLayout(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("DeleteLayout() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerNewLayout(t *testing.T) {
	handler, server := newTestWebHandlerWithMockAPI(nil, nil, nil, nil, nil, nil)
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, "/ssg/new-layout", nil)
	ctx := feat.NewContextWithSite("test-site", uuid.New())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.NewLayout(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("NewLayout() status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestWebHandlerShowLayout(t *testing.T) {
	layoutID := uuid.New()
	tests := []struct {
		name           string
		queryID        string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name:    "shows layout successfully",
			queryID: layoutID.String(),
			getResp: map[string]interface{}{
				"layout": map[string]interface{}{
					"id":   layoutID.String(),
					"name": "Test Layout",
					"code": "<html></html>",
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
			queryID:        layoutID.String(),
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/show-layout?id="+tt.queryID, nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.ShowLayout(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ShowLayout() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}
