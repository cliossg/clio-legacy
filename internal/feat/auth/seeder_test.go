package auth_test

import (
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
