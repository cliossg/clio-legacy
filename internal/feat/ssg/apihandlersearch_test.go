package ssg

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hermesgen/hm"
)

func TestAPIHandlerSearchContent(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		setupRepo      func(*mockServiceRepo)
		wantStatusCode int
	}{
		{
			name:        "searches content successfully with query",
			queryParams: "?search=test&page=1",
			setupRepo: func(m *mockServiceRepo) {
				// Mock will return empty results by default
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:        "searches content without query",
			queryParams: "",
			setupRepo: func(m *mockServiceRepo) {
				// Mock will return empty results by default
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:        "handles invalid page number gracefully",
			queryParams: "?search=test&page=invalid",
			setupRepo: func(m *mockServiceRepo) {
				// Page defaults to 1 when invalid
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:        "handles page number zero",
			queryParams: "?page=0",
			setupRepo: func(m *mockServiceRepo) {
				// Page defaults to 1 when zero or negative
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:        "fails when service returns error",
			queryParams: "?search=test",
			setupRepo: func(m *mockServiceRepo) {
				m.getContentErr = fmt.Errorf("db error")
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

			req := httptest.NewRequest(http.MethodGet, "/ssg/search"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			handler.SearchContent(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("SearchContent() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}
