package core

import (
	"context"
	"embed"
	"testing"

	"github.com/hermesgen/hm"
)

//go:embed testdata
var testAssetsFS embed.FS

func TestNewApp(t *testing.T) {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}

	app := NewApp("test-app", "1.0.0", testAssetsFS, params)

	if app == nil {
		t.Fatal("NewApp() returned nil")
	}
}

func TestAppSetup(t *testing.T) {
	cfg := hm.NewConfig()
	cfg.Set(hm.Key.DBSQLiteDSN, "file::memory:?cache=shared")
	log := hm.NewLogger("error")
	params := hm.XParams{Cfg: cfg, Log: log}

	app := NewApp("test-app", "1.0.0", testAssetsFS, params)
	err := app.Setup(context.Background())

	if err != nil {
		t.Errorf("Setup() error = %v", err)
	}
}
