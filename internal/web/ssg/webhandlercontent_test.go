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

func TestWebHandlerListContent(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name:        "lists content successfully",
			queryParams: "",
			getResp: map[string]interface{}{
				"contents": []feat.Content{
					{Heading: "Content 1"},
					{Heading: "Content 2"},
				},
				"page":        1,
				"total_pages": 1,
				"total_count": 2,
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:        "lists content with pagination",
			queryParams: "?page=2",
			getResp: map[string]interface{}{
				"contents": []feat.Content{
					{Heading: "Content 3"},
				},
				"page":        2,
				"total_pages": 3,
				"total_count": 51,
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:        "lists content with search",
			queryParams: "?search=test",
			getResp: map[string]interface{}{
				"contents": []feat.Content{
					{Heading: "Test Content"},
				},
				"page":        1,
				"total_pages": 1,
				"total_count": 1,
				"search":      "test",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails when API returns error",
			queryParams:    "",
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/list-content"+tt.queryParams, nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.ListContent(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ListContent() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerCreateContent(t *testing.T) {
	sectionID := uuid.New()
	userID := uuid.New()
	tests := []struct {
		name           string
		formData       url.Values
		postResp       interface{}
		postErr        error
		wantStatusCode int
	}{
		{
			name: "creates content successfully",
			formData: url.Values{
				"heading":    []string{"Test Heading"},
				"body":       []string{"Test Body"},
				"section_id": []string{sectionID.String()},
				"user_id":    []string{userID.String()},
			},
			postResp: map[string]interface{}{
				"content": map[string]interface{}{
					"id":         uuid.New().String(),
					"heading":    "Test Heading",
					"body":       "Test Body",
					"section_id": sectionID.String(),
					"user_id":    userID.String(),
				},
			},
			wantStatusCode: http.StatusSeeOther,
		},
		{
			name: "fails with invalid form data",
			formData: url.Values{
				"heading": []string{""},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when API returns error",
			formData: url.Values{
				"heading":    []string{"Test Heading"},
				"body":       []string{"Test Body"},
				"section_id": []string{sectionID.String()},
				"user_id":    []string{userID.String()},
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
			req := httptest.NewRequest(http.MethodPost, "/ssg/create-content", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.CreateContent(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("CreateContent() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerEditContent(t *testing.T) {
	contentID := uuid.New()
	sectionID := uuid.New()
	userID := uuid.New()
	tests := []struct {
		name           string
		queryID        string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name:    "shows edit form successfully",
			queryID: contentID.String(),
			getResp: map[string]interface{}{
				"content": map[string]interface{}{
					"id":         contentID.String(),
					"heading":    "Test Content",
					"body":       "Test Body",
					"section_id": sectionID.String(),
					"user_id":    userID.String(),
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
			queryID:        contentID.String(),
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/edit-content?id="+tt.queryID, nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.EditContent(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("EditContent() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerUpdateContent(t *testing.T) {
	contentID := uuid.New()
	sectionID := uuid.New()
	userID := uuid.New()
	tests := []struct {
		name           string
		formData       url.Values
		putErr         error
		wantStatusCode int
	}{
		{
			name: "updates content successfully",
			formData: url.Values{
				"id":         []string{contentID.String()},
				"heading":    []string{"Updated Heading"},
				"body":       []string{"Updated Body"},
				"section_id": []string{sectionID.String()},
				"user_id":    []string{userID.String()},
			},
			wantStatusCode: http.StatusSeeOther,
		},
		{
			name: "fails with invalid form data",
			formData: url.Values{
				"id":      []string{contentID.String()},
				"heading": []string{""},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when API returns error",
			formData: url.Values{
				"id":         []string{contentID.String()},
				"heading":    []string{"Updated Heading"},
				"body":       []string{"Updated Body"},
				"section_id": []string{sectionID.String()},
				"user_id":    []string{userID.String()},
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
			req := httptest.NewRequest(http.MethodPost, "/ssg/update-content", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.UpdateContent(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("UpdateContent() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerDeleteContent(t *testing.T) {
	contentID := uuid.New()
	tests := []struct {
		name           string
		formData       url.Values
		deleteErr      error
		wantStatusCode int
	}{
		{
			name: "deletes content successfully",
			formData: url.Values{
				"id": []string{contentID.String()},
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
				"id": []string{contentID.String()},
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
			req := httptest.NewRequest(http.MethodPost, "/ssg/delete-content", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.DeleteContent(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("DeleteContent() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerNewContent(t *testing.T) {
	handler, server := newTestWebHandlerWithMockAPI(nil, nil, nil, nil, nil, nil)
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, "/ssg/new-content", nil)
	ctx := feat.NewContextWithSite("test-site", uuid.New())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.NewContent(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("NewContent() status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestWebHandlerShowContent(t *testing.T) {
	contentID := uuid.New()
	sectionID := uuid.New()
	userID := uuid.New()
	tests := []struct {
		name           string
		queryID        string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name:    "shows content successfully",
			queryID: contentID.String(),
			getResp: map[string]interface{}{
				"content": map[string]interface{}{
					"id":         contentID.String(),
					"heading":    "Test Content",
					"body":       "Test Body",
					"section_id": sectionID.String(),
					"user_id":    userID.String(),
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
			queryID:        contentID.String(),
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/show-content?id="+tt.queryID, nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.ShowContent(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ShowContent() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerSearchContent(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name:        "searches content successfully",
			queryParams: "?search=test",
			getResp: map[string]interface{}{
				"contents": []feat.Content{
					{Heading: "Test Content"},
				},
				"page":        1,
				"total_pages": 1,
				"total_count": 1,
				"search":      "test",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:        "searches with pagination",
			queryParams: "?search=test&page=2",
			getResp: map[string]interface{}{
				"contents": []feat.Content{
					{Heading: "Test Content 2"},
				},
				"page":        2,
				"total_pages": 3,
				"total_count": 60,
				"search":      "test",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails when API returns error",
			queryParams:    "?search=test",
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/search-content"+tt.queryParams, nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.SearchContent(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("SearchContent() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerGenerateHTML(t *testing.T) {
	tests := []struct {
		name           string
		postErr        error
		wantStatusCode int
	}{
		{
			name:           "generates HTML successfully",
			wantStatusCode: http.StatusSeeOther,
		},
		{
			name:           "fails when API returns error",
			postErr:        fmt.Errorf("api error"),
			wantStatusCode: http.StatusSeeOther,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(nil, nil, nil, tt.postErr, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodPost, "/ssg/generate-html", nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.GenerateHTML(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GenerateHTML() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestGeneratePageNumbers(t *testing.T) {
	tests := []struct {
		name        string
		currentPage int
		totalPages  int
		want        []int
	}{
		{
			name:        "total pages less than or equal to 7",
			currentPage: 3,
			totalPages:  5,
			want:        []int{1, 2, 3, 4, 5},
		},
		{
			name:        "current page in first 4 pages",
			currentPage: 2,
			totalPages:  10,
			want:        []int{1, 2, 3, 4, 5, -1, 10},
		},
		{
			name:        "current page in last 4 pages",
			currentPage: 9,
			totalPages:  10,
			want:        []int{1, -1, 6, 7, 8, 9, 10},
		},
		{
			name:        "current page in middle",
			currentPage: 5,
			totalPages:  10,
			want:        []int{1, -1, 4, 5, 6, -1, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generatePageNumbers(tt.currentPage, tt.totalPages)
			if len(got) != len(tt.want) {
				t.Errorf("generatePageNumbers() length = %d, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("generatePageNumbers()[%d] = %d, want %d", i, got[i], tt.want[i])
				}
			}
		})
	}
}
