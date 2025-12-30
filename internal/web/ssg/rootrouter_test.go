package ssg

import (
	"testing"

	"github.com/hermesgen/hm"
)

func TestNewRootRouter(t *testing.T) {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}
	tm := hm.NewTemplateManager(testAssetsFS, params)
	flash := hm.NewFlashManager(params)
	handler := NewWebHandler(tm, flash, nil, nil, nil, params)

	rootRouter := NewRootRouter(handler, params)

	if rootRouter == nil {
		t.Fatal("NewRootRouter() returned nil")
	}
}

func TestRootRouterSetupRoutes(t *testing.T) {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}
	tm := hm.NewTemplateManager(testAssetsFS, params)
	flash := hm.NewFlashManager(params)
	handler := NewWebHandler(tm, flash, nil, nil, nil, params)
	rootRouter := NewRootRouter(handler, params)
	router := hm.NewWebRouter("test-router", params)

	rootRouter.SetupRoutes(router)

	if router == nil {
		t.Error("SetupRoutes() should not modify router to nil")
	}
}
