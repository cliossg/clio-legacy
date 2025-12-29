package sqlite

import (
	"context"
	"embed"
	"testing"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/feat/auth"
	"github.com/hermesgen/hm"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed assets
var testAssetsFS embed.FS

func setupTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user (
		id TEXT PRIMARY KEY,
		short_id TEXT NOT NULL DEFAULT '',
		name TEXT NOT NULL DEFAULT '',
		username TEXT NOT NULL DEFAULT '',
		email TEXT NOT NULL DEFAULT '',
		created_by TEXT NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
		updated_by TEXT NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		t.Fatalf("create user table: %v", err)
	}

	return db
}

func setupTestRepo(t *testing.T) *ClioRepo {
	db := setupTestDB(t)
	cfg := hm.NewConfig()
	log := hm.NewLogger("debug")
	params := hm.XParams{Cfg: cfg, Log: log}

	qm := hm.NewQueryManager(testAssetsFS, "sqlite", params)
	ctx := context.Background()
	err := qm.Setup(ctx)
	if err != nil {
		t.Fatalf("setup query manager: %v", err)
	}

	repo := NewClioRepo(qm, params)
	repo.SetDB(db)

	return repo
}

func TestClioRepoCreateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    *auth.User
		wantErr bool
	}{
		{
			name: "creates user successfully",
			user: &auth.User{
				ID:       uuid.New(),
				Username: "testuser",
				Email:    "test@example.com",
				Name:     "Test User",
			},
			wantErr: false,
		},
	}

	repo := setupTestRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateUser(ctx, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClioRepoGetUser(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	user := &auth.User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
	}
	repo.CreateUser(ctx, user)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "gets existing user",
			id:      user.ID,
			wantErr: false,
		},
		{
			name:    "fails with non-existent user",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetUser(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != tt.id {
				t.Errorf("GetUser() got ID = %v, want %v", got.ID, tt.id)
			}
		})
	}
}

func TestClioRepoGetUserByUsername(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	user := &auth.User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
	}
	repo.CreateUser(ctx, user)

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "gets user by username",
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "fails with non-existent username",
			username: "nonexistent",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetUserByUsername(ctx, tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Username != tt.username {
				t.Errorf("GetUserByUsername() got username = %v, want %v", got.Username, tt.username)
			}
		})
	}
}

func TestClioRepoGetUsers(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	user1 := &auth.User{
		ID:           uuid.New(),
		Username:     "user1",
		Name:     "User One",
	}
	user2 := &auth.User{
		ID:           uuid.New(),
		Username:     "user2",
		Name:     "User Two",
	}
	repo.CreateUser(ctx, user1)
	repo.CreateUser(ctx, user2)

	users, err := repo.GetUsers(ctx)
	if err != nil {
		t.Errorf("GetUsers() error = %v", err)
		return
	}

	if len(users) != 2 {
		t.Errorf("GetUsers() got %d users, want 2", len(users))
	}
}

func TestClioRepoUpdateUser(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	user := &auth.User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
	}
	repo.CreateUser(ctx, user)

	tests := []struct {
		name    string
		user    *auth.User
		wantErr bool
	}{
		{
			name: "updates user successfully",
			user: &auth.User{
				ID:           user.ID,
				Username:     "testuser",
				Name:     "Updated Name",
				Email:        "newemail@example.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateUser(ctx, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				updated, _ := repo.GetUser(ctx, tt.user.ID)
				if updated.Email != tt.user.Email {
					t.Errorf("UpdateUser() email = %v, want %v", updated.Email, tt.user.Email)
				}
			}
		})
	}
}

func TestClioRepoDeleteUser(t *testing.T) {
	repo := setupTestRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	user := &auth.User{
		ID:           uuid.New(),
		Username:     "testuser",
		Name:     "Delete Test",
	}
	repo.CreateUser(ctx, user)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "deletes user successfully",
			id:      user.ID,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteUser(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				_, err := repo.GetUser(ctx, tt.id)
				if err == nil {
					t.Error("DeleteUser() user still exists after deletion")
				}
			}
		})
	}
}
