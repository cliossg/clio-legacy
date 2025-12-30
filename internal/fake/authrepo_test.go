package fake_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/fake"
	"github.com/hermesgen/clio/internal/feat/auth"
)

func TestAuthRepoCreateUser(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.AuthRepo)
		user        *auth.User
		expectedErr error
		expectCalls int
	}{
		{
			name:        "creates user successfully",
			setupFake:   func(f *fake.AuthRepo) {},
			user:        &auth.User{ID: uuid.New(), Username: "testuser"},
			expectedErr: nil,
			expectCalls: 1,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.AuthRepo) {
				f.CreateUserFn = func(ctx context.Context, user *auth.User) error {
					return errors.New("db error")
				}
			},
			user:        &auth.User{ID: uuid.New(), Username: "testuser"},
			expectedErr: errors.New("db error"),
			expectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewAuthRepo()
			tt.setupFake(f)

			err := f.CreateUser(context.Background(), tt.user)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(f.CreateUserCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.CreateUserCalls))
			}
		})
	}
}

func TestAuthRepoGetUser(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name         string
		setupFake    func(f *fake.AuthRepo)
		id           uuid.UUID
		expectedUser auth.User
		expectedErr  error
		expectCalls  int
	}{
		{
			name: "gets existing user",
			setupFake: func(f *fake.AuthRepo) {
				f.CreateUser(context.Background(), &auth.User{ID: userID, Username: "testuser"})
			},
			id:           userID,
			expectedUser: auth.User{ID: userID, Username: "testuser"},
			expectedErr:  nil,
			expectCalls:  1,
		},
		{
			name:         "returns empty user when not found",
			setupFake:    func(f *fake.AuthRepo) {},
			id:           uuid.New(),
			expectedUser: auth.User{},
			expectedErr:  sql.ErrNoRows,
			expectCalls:  1,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.AuthRepo) {
				f.GetUserFn = func(ctx context.Context, id uuid.UUID) (auth.User, error) {
					return auth.User{}, errors.New("db error")
				}
			},
			id:           userID,
			expectedUser: auth.User{},
			expectedErr:  errors.New("db error"),
			expectCalls:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewAuthRepo()
			tt.setupFake(f)

			user, err := f.GetUser(context.Background(), tt.id)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if user.ID != tt.expectedUser.ID || user.Username != tt.expectedUser.Username {
				t.Errorf("expected user %+v, got %+v", tt.expectedUser, user)
			}

			if len(f.GetUserCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.GetUserCalls))
			}
		})
	}
}

func TestAuthRepoGetUsers(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.AuthRepo)
		expectedLen int
		expectedErr error
		expectCalls int
	}{
		{
			name: "returns all users",
			setupFake: func(f *fake.AuthRepo) {
				f.CreateUser(context.Background(), &auth.User{ID: uuid.New(), Username: "user1"})
				f.CreateUser(context.Background(), &auth.User{ID: uuid.New(), Username: "user2"})
			},
			expectedLen: 2,
			expectedErr: nil,
			expectCalls: 1,
		},
		{
			name:        "returns empty slice when no users",
			setupFake:   func(f *fake.AuthRepo) {},
			expectedLen: 0,
			expectedErr: nil,
			expectCalls: 1,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.AuthRepo) {
				f.GetUsersFn = func(ctx context.Context) ([]auth.User, error) {
					return nil, errors.New("db error")
				}
			},
			expectedLen: 0,
			expectedErr: errors.New("db error"),
			expectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewAuthRepo()
			tt.setupFake(f)

			users, err := f.GetUsers(context.Background())

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(users) != tt.expectedLen {
				t.Errorf("expected %d users, got %d", tt.expectedLen, len(users))
			}

			if len(f.GetUsersCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.GetUsersCalls))
			}
		})
	}
}

func TestAuthRepoUpdateUser(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.AuthRepo)
		user        *auth.User
		expectedErr error
		expectCalls int
	}{
		{
			name: "updates user successfully",
			setupFake: func(f *fake.AuthRepo) {
				f.CreateUser(context.Background(), &auth.User{ID: userID, Username: "oldname"})
			},
			user:        &auth.User{ID: userID, Username: "newname"},
			expectedErr: nil,
			expectCalls: 1,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.AuthRepo) {
				f.UpdateUserFn = func(ctx context.Context, user *auth.User) error {
					return errors.New("db error")
				}
			},
			user:        &auth.User{ID: userID, Username: "newname"},
			expectedErr: errors.New("db error"),
			expectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewAuthRepo()
			tt.setupFake(f)

			err := f.UpdateUser(context.Background(), tt.user)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(f.UpdateUserCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.UpdateUserCalls))
			}
		})
	}
}

func TestAuthRepoDeleteUser(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.AuthRepo)
		id          uuid.UUID
		expectedErr error
		expectCalls int
	}{
		{
			name: "deletes user successfully",
			setupFake: func(f *fake.AuthRepo) {
				f.CreateUser(context.Background(), &auth.User{ID: userID, Username: "testuser"})
			},
			id:          userID,
			expectedErr: nil,
			expectCalls: 1,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.AuthRepo) {
				f.DeleteUserFn = func(ctx context.Context, id uuid.UUID) error {
					return errors.New("db error")
				}
			},
			id:          userID,
			expectedErr: errors.New("db error"),
			expectCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewAuthRepo()
			tt.setupFake(f)

			err := f.DeleteUser(context.Background(), tt.id)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(f.DeleteUserCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.DeleteUserCalls))
			}
		})
	}
}

func TestAuthRepoGetUserByUsername(t *testing.T) {
	tests := []struct {
		name         string
		setupFake    func(f *fake.AuthRepo)
		username     string
		expectedUser auth.User
		expectedErr  error
		expectCalls  int
	}{
		{
			name: "gets user by username",
			setupFake: func(f *fake.AuthRepo) {
				f.CreateUser(context.Background(), &auth.User{ID: uuid.New(), Username: "testuser"})
			},
			username:     "testuser",
			expectedUser: auth.User{Username: "testuser"},
			expectedErr:  nil,
			expectCalls:  1,
		},
		{
			name:         "returns empty user when not found",
			setupFake:    func(f *fake.AuthRepo) {},
			username:     "nonexistent",
			expectedUser: auth.User{},
			expectedErr:  nil,
			expectCalls:  1,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.AuthRepo) {
				f.GetUserByUsernameFn = func(ctx context.Context, username string) (auth.User, error) {
					return auth.User{}, errors.New("db error")
				}
			},
			username:     "testuser",
			expectedUser: auth.User{},
			expectedErr:  errors.New("db error"),
			expectCalls:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewAuthRepo()
			tt.setupFake(f)

			user, err := f.GetUserByUsername(context.Background(), tt.username)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if user.Username != tt.expectedUser.Username {
				t.Errorf("expected username %v, got %v", tt.expectedUser.Username, user.Username)
			}

			if len(f.GetUserByUsernameCalls) != tt.expectCalls {
				t.Errorf("expected %d calls, got %d", tt.expectCalls, len(f.GetUserByUsernameCalls))
			}
		})
	}
}
