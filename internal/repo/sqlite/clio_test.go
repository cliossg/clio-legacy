package sqlite

import (
	"context"
	"testing"

	"github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

func TestClioRepoGetDB(t *testing.T) {
	repo, _ := setupTestSsgRepo(t)
	defer repo.db.Close()

	db := repo.GetDB()
	if db == nil {
		t.Error("GetDB() returned nil")
	}
}

func TestClioRepoStop(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *ClioRepo
		wantErr bool
	}{
		{
			name: "stops repo with open database",
			setup: func() *ClioRepo {
				repo, _ := setupTestSsgRepo(t)
				return repo
			},
			wantErr: false,
		},
		{
			name: "stops repo with nil database",
			setup: func() *ClioRepo {
				repo, _ := setupTestSsgRepo(t)
				repo.db = nil
				return repo
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setup()
			err := repo.Stop(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClioRepoBeginTx(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()

	ctx := ssg.NewContextWithSite("test-site", siteID)

	ctxWithTx, tx, err := repo.BeginTx(ctx)
	if err != nil {
		t.Errorf("BeginTx() error = %v", err)
		return
	}

	if tx == nil {
		t.Error("BeginTx() returned nil transaction")
	}

	if ctxWithTx == nil {
		t.Error("BeginTx() returned nil context")
	}

	// Clean up transaction
	tx.Rollback()
}

func TestClioRepoSetup(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*testing.T) *ClioRepo
		wantErr bool
	}{
		{
			name: "returns early when database already set",
			setup: func(t *testing.T) *ClioRepo {
				repo, _ := setupTestSsgRepo(t)
				return repo
			},
			wantErr: false,
		},
		{
			name: "sets up database with valid DSN",
			setup: func(t *testing.T) *ClioRepo {
				tmpDB := t.TempDir() + "/test_setup.db"
				cfg := hm.NewConfig()
				cfg.Set(hm.Key.DBSQLiteDSN, tmpDB)
				log := hm.NewLogger("debug")
				params := hm.XParams{Cfg: cfg, Log: log}
				qm := hm.NewQueryManager(testAssetsFS, "sqlite", params)
				return NewClioRepo(qm, params)
			},
			wantErr: false,
		},
		{
			name: "returns error when DSN not in config",
			setup: func(t *testing.T) *ClioRepo {
				cfg := hm.NewConfig()
				log := hm.NewLogger("debug")
				params := hm.XParams{Cfg: cfg, Log: log}
				qm := hm.NewQueryManager(testAssetsFS, "sqlite", params)
				return NewClioRepo(qm, params)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setup(t)
			defer func() {
				if repo.db != nil {
					repo.db.Close()
				}
			}()

			err := repo.Setup(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Setup() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && repo.db == nil {
				t.Error("Setup() database is nil after setup")
			}
		})
	}
}
