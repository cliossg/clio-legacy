package ssg

import (
	"testing"

	"github.com/hermesgen/hm"
)

func TestNewWebRouter(t *testing.T) {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}
	tm := hm.NewTemplateManager(testAssetsFS, params)
	flash := hm.NewFlashManager(params)
	handler := NewWebHandler(tm, flash, nil, nil, nil, params)
	mw := []hm.Middleware{}

	router := NewWebRouter(handler, mw, params)

	if router == nil {
		t.Fatal("NewWebRouter() returned nil")
	}
}
