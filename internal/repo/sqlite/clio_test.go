package sqlite

import (
	"context"
	"testing"

	"github.com/hermesgen/clio/internal/feat/ssg"
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
		setup   func() *ClioRepo
		wantErr bool
	}{
		{
			name: "returns early when database already set",
			setup: func() *ClioRepo {
				repo, _ := setupTestSsgRepo(t)
				return repo
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setup()
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
