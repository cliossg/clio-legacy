package auth_test

import (
	"context"
	"embed"
	"testing"

	"github.com/hermesgen/clio/internal/fake"
	"github.com/hermesgen/clio/internal/feat/auth"
	"github.com/hermesgen/hm"
)

//go:embed testdata
var testAssetsFS embed.FS

func TestNewSeeder(t *testing.T) {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}
	repo := fake.NewAuthRepo()

	seeder := auth.NewSeeder(testAssetsFS, "test", repo, params)

	if seeder == nil {
		t.Fatal("NewSeeder() returned nil")
	}
}

func TestSeederSetup(t *testing.T) {
	cfg := hm.NewConfig()
	cfg.Set(hm.Key.DBSQLiteDSN, "file::memory:?cache=shared")
	log := hm.NewLogger("error")
	params := hm.XParams{Cfg: cfg, Log: log}
	repo := fake.NewAuthRepo()

	seeder := auth.NewSeeder(testAssetsFS, "test", repo, params)
	err := seeder.Setup(context.Background())

	if err != nil {
		t.Errorf("Setup() error = %v", err)
	}
}
