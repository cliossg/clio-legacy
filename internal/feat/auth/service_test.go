package auth_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/fake"
	"github.com/hermesgen/clio/internal/feat/auth"
	"github.com/hermesgen/hm"
)

func newTestService(repo auth.Repo) *auth.BaseService {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}
	return auth.NewService(repo, params)
}

func TestBaseServiceGetUserByID(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(*fake.AuthRepo)
		userID    uuid.UUID
		wantErr   bool
	}{
		{
			name: "gets user by ID successfully",
			setupRepo: func(f *fake.AuthRepo) {
				userID := uuid.New()
				f.CreateUser(context.Background(), &auth.User{ID: userID, Username: "testuser"})
			},
			userID:  uuid.New(),
			wantErr: false,
		},
		{
			name: "fails when repo returns error",
			setupRepo: func(f *fake.AuthRepo) {
				f.GetUserFn = func(ctx context.Context, id uuid.UUID) (auth.User, error) {
					return auth.User{}, fmt.Errorf("db error")
				}
			},
			userID:  uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := fake.NewAuthRepo()
			if tt.name == "gets user by ID successfully" {
				tt.userID = uuid.New()
				repo.CreateUser(context.Background(), &auth.User{ID: tt.userID, Username: "testuser"})
			} else {
				tt.setupRepo(repo)
			}
			svc := newTestService(repo)

			got, err := svc.GetUserByID(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != tt.userID {
				t.Errorf("GetUserByID() got ID = %v, want %v", got.ID, tt.userID)
			}
		})
	}
}

func TestBaseServiceGetUsers(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(*fake.AuthRepo)
		wantLen   int
		wantErr   bool
	}{
		{
			name: "gets all users successfully",
			setupRepo: func(f *fake.AuthRepo) {
				f.CreateUser(context.Background(), &auth.User{ID: uuid.New(), Username: "user1"})
				f.CreateUser(context.Background(), &auth.User{ID: uuid.New(), Username: "user2"})
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "fails when repo returns error",
			setupRepo: func(f *fake.AuthRepo) {
				f.GetUsersFn = func(ctx context.Context) ([]auth.User, error) {
					return nil, fmt.Errorf("db error")
				}
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := fake.NewAuthRepo()
			tt.setupRepo(repo)
			svc := newTestService(repo)

			got, err := svc.GetUsers(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.wantLen {
				t.Errorf("GetUsers() got %d users, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestBaseServiceGetUser(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(*fake.AuthRepo, uuid.UUID)
		wantErr   bool
	}{
		{
			name: "gets user successfully",
			setupRepo: func(f *fake.AuthRepo, id uuid.UUID) {
				f.CreateUser(context.Background(), &auth.User{ID: id, Username: "testuser"})
			},
			wantErr: false,
		},
		{
			name: "returns error when user not found",
			setupRepo: func(f *fake.AuthRepo, id uuid.UUID) {
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := fake.NewAuthRepo()
			userID := uuid.New()
			tt.setupRepo(repo, userID)
			svc := newTestService(repo)

			got, err := svc.GetUser(context.Background(), userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.name == "gets user successfully" && got.ID != userID {
				t.Errorf("GetUser() got ID = %v, want %v", got.ID, userID)
			}
		})
	}
}

func TestBaseServiceCreateUser(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(*fake.AuthRepo)
		user      *auth.User
		wantErr   bool
	}{
		{
			name: "creates user successfully",
			setupRepo: func(f *fake.AuthRepo) {
			},
			user: &auth.User{
				ID:       uuid.New(),
				Username: "newuser",
			},
			wantErr: false,
		},
		{
			name: "fails when repo returns error",
			setupRepo: func(f *fake.AuthRepo) {
				f.CreateUserFn = func(ctx context.Context, user *auth.User) error {
					return fmt.Errorf("db error")
				}
			},
			user: &auth.User{
				ID:       uuid.New(),
				Username: "newuser",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := fake.NewAuthRepo()
			tt.setupRepo(repo)
			svc := newTestService(repo)

			err := svc.CreateUser(context.Background(), tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBaseServiceUpdateUser(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(*fake.AuthRepo, *auth.User)
		wantErr   bool
	}{
		{
			name: "updates user successfully",
			setupRepo: func(f *fake.AuthRepo, user *auth.User) {
				f.CreateUser(context.Background(), user)
			},
			wantErr: false,
		},
		{
			name: "fails when repo returns error",
			setupRepo: func(f *fake.AuthRepo, user *auth.User) {
				f.CreateUser(context.Background(), user)
				f.UpdateUserFn = func(ctx context.Context, user *auth.User) error {
					return fmt.Errorf("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := fake.NewAuthRepo()
			user := &auth.User{
				ID:       uuid.New(),
				Username: "testuser",
			}
			tt.setupRepo(repo, user)
			svc := newTestService(repo)

			user.Username = "updated"
			err := svc.UpdateUser(context.Background(), user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBaseServiceDeleteUser(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(*fake.AuthRepo, uuid.UUID)
		wantErr   bool
	}{
		{
			name: "deletes user successfully",
			setupRepo: func(f *fake.AuthRepo, id uuid.UUID) {
				f.CreateUser(context.Background(), &auth.User{ID: id})
			},
			wantErr: false,
		},
		{
			name: "fails when repo returns error",
			setupRepo: func(f *fake.AuthRepo, id uuid.UUID) {
				f.CreateUser(context.Background(), &auth.User{ID: id})
				f.DeleteUserFn = func(ctx context.Context, id uuid.UUID) error {
					return fmt.Errorf("db error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := fake.NewAuthRepo()
			userID := uuid.New()
			tt.setupRepo(repo, userID)
			svc := newTestService(repo)

			err := svc.DeleteUser(context.Background(), userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
