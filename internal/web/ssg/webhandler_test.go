package ssg

import (
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

func TestWebHandlerServeStaticImage(t *testing.T) {
	tests := []struct {
		name           string
		siteSlug       string
		path           string
		wantStatusCode int
	}{
		{
			name:           "returns error when site not in context",
			path:           "/static/images/test.png",
			wantStatusCode: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := hm.NewConfig()
			params := hm.XParams{Cfg: cfg}
			tm := hm.NewTemplateManager(testAssetsFS, params)
			flash := hm.NewFlashManager(params)
			handler := NewWebHandler(tm, flash, nil, nil, nil, params)

			req := httptest.NewRequest("GET", tt.path, nil)
			if tt.siteSlug != "" {
				ctx := feat.NewContextWithSite(tt.siteSlug, uuid.Nil)
				req = req.WithContext(ctx)
			}
			w := httptest.NewRecorder()

			handler.ServeStaticImage(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ServeStaticImage() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}
