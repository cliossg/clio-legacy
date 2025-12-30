package core

import (
	"net/http/httptest"
	"testing"

	"github.com/hermesgen/hm"
)

func TestNewMultiSitePreviewHandler(t *testing.T) {
	cfg := hm.NewConfig()
	log := hm.NewLogger("error")
	params := hm.XParams{Cfg: cfg, Log: log}

	handler := NewMultiSitePreviewHandler(params)

	if handler == nil {
		t.Fatal("NewMultiSitePreviewHandler() returned nil")
	}
}

func TestMultiSitePreviewHandlerServeHTTP(t *testing.T) {
	tests := []struct {
		name           string
		host           string
		path           string
		wantStatusCode int
	}{
		{
			name:           "rejects invalid host",
			host:           "example.com",
			path:           "/",
			wantStatusCode: 400,
		},
		{
			name:           "handles localhost default",
			host:           "localhost",
			path:           "/",
			wantStatusCode: 404,
		},
		{
			name:           "handles subdomain with port",
			host:           "test.localhost:8080",
			path:           "/",
			wantStatusCode: 404,
		},
		{
			name:           "handles path traversal attempt",
			host:           "test.localhost",
			path:           "/../../../etc/passwd",
			wantStatusCode: 403,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := hm.NewConfig()
			log := hm.NewLogger("error")
			params := hm.XParams{Cfg: cfg, Log: log}
			handler := NewMultiSitePreviewHandler(params)

			req := httptest.NewRequest("GET", tt.path, nil)
			req.Host = tt.host
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ServeHTTP() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}
