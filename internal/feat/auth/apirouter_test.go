package auth_test

import (
	"testing"

	"github.com/hermesgen/clio/internal/fake"
	"github.com/hermesgen/clio/internal/feat/auth"
	"github.com/hermesgen/hm"
)

func TestNewAPIRouter(t *testing.T) {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}
	repo := fake.NewAuthRepo()
	handler := auth.NewAPIHandler("test-handler", repo, params)
	mw := []hm.Middleware{}

	router := auth.NewAPIRouter(handler, mw, params)

	if router == nil {
		t.Fatal("NewAPIRouter() returned nil")
	}
}
