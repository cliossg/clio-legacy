package core_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/hermesgen/clio/internal/core"
	"github.com/hermesgen/hm"
)

func TestNewAdminFileServer(t *testing.T) {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}

	server := core.NewAdminFileServer(params)

	if server == nil {
		t.Fatal("NewAdminFileServer() returned nil")
	}
}

func TestAdminFileServerSetup(t *testing.T) {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}

	server := core.NewAdminFileServer(params)
	err := server.Setup(context.Background())

	if err != nil {
		t.Errorf("Setup() error = %v", err)
	}
}

func TestAdminFileServerHandler(t *testing.T) {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}

	server := core.NewAdminFileServer(params)
	handler := server.Handler()

	if handler == nil {
		t.Error("Handler() returned nil")
	}

	req := httptest.NewRequest("GET", "/static/images/test.jpg", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != 200 && w.Code != 404 {
		t.Logf("Handler() status = %d (expected 200 or 404 for non-existent file)", w.Code)
	}
}
