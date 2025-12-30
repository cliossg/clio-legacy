package ssg

import (
	"testing"

	"github.com/hermesgen/hm"
)

func TestNewAPIRouter(t *testing.T) {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}
	repo := newMockServiceRepo()
	svc := newTestService(repo)
	handler := NewAPIHandler("test-handler", svc, nil, params)
	mw := []hm.Middleware{}

	router := NewAPIRouter(handler, mw, params)

	if router == nil {
		t.Fatal("NewAPIRouter() returned nil")
	}
}
