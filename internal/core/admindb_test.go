package core_test

import (
	"context"
	"embed"
	"testing"

	"github.com/hermesgen/clio/internal/core"
	"github.com/hermesgen/hm"
)

//go:embed *.go
var testCoreFS embed.FS

func TestNewAdminDBManager(t *testing.T) {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}
	migrator := hm.NewMigrator(testCoreFS, "sqlite3", params)

	manager := core.NewAdminDBManager(testCoreFS, "sqlite3", migrator, params)

	if manager == nil {
		t.Fatal("NewAdminDBManager() returned nil")
	}
}

func TestAdminDBManagerGetDB(t *testing.T) {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}
	migrator := hm.NewMigrator(testCoreFS, "sqlite3", params)

	manager := core.NewAdminDBManager(testCoreFS, "sqlite3", migrator, params)

	db := manager.GetDB()
	if db != nil {
		t.Error("GetDB() should return nil before setup")
	}
}

func TestAdminDBManagerStop(t *testing.T) {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}
	migrator := hm.NewMigrator(testCoreFS, "sqlite3", params)

	manager := core.NewAdminDBManager(testCoreFS, "sqlite3", migrator, params)

	err := manager.Stop(context.Background())
	if err != nil {
		t.Errorf("Stop() error = %v when db is nil", err)
	}
}
