package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
)

type mockAuthRepo struct {
	hm.Core
	users        map[uuid.UUID]User
	getUsersErr  error
	getUserErr   error
	createUserErr error
	updateUserErr error
	deleteUserErr error
}

func newMockAuthRepo() *mockAuthRepo {
	cfg := hm.NewConfig()
	return &mockAuthRepo{
		Core:  hm.NewCore("mock-auth-repo", hm.XParams{Cfg: cfg}),
		users: make(map[uuid.UUID]User),
	}
}

func (m *mockAuthRepo) GetUsers(ctx context.Context) ([]User, error) {
	if m.getUsersErr != nil {
		return nil, m.getUsersErr
	}
	var users []User
	for _, u := range m.users {
		users = append(users, u)
	}
	return users, nil
}

func (m *mockAuthRepo) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	if m.getUserErr != nil {
		return User{}, m.getUserErr
	}
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return User{}, fmt.Errorf("user not found")
}

func (m *mockAuthRepo) GetUserByUsername(ctx context.Context, username string) (User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return User{}, fmt.Errorf("user not found")
}

func (m *mockAuthRepo) CreateUser(ctx context.Context, user *User) error {
	if m.createUserErr != nil {
		return m.createUserErr
	}
	m.users[user.ID] = *user
	return nil
}

func (m *mockAuthRepo) UpdateUser(ctx context.Context, user *User) error {
	if m.updateUserErr != nil {
		return m.updateUserErr
	}
	m.users[user.ID] = *user
	return nil
}

func (m *mockAuthRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if m.deleteUserErr != nil {
		return m.deleteUserErr
	}
	delete(m.users, id)
	return nil
}

func (m *mockAuthRepo) BeginTx(ctx context.Context) (context.Context, hm.Tx, error) {
	return ctx, nil, nil
}

func (m *mockAuthRepo) Query() *hm.QueryManager {
	return nil
}

func (m *mockAuthRepo) GetDB() interface{} {
	return nil
}

func newTestService(repo Repo) *BaseService {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}
	return NewService(repo, params)
}

func TestBaseServiceGetUserByID(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(*mockAuthRepo)
		userID    uuid.UUID
		wantErr   bool
	}{
		{
			name: "gets user by ID successfully",
			setupRepo: func(m *mockAuthRepo) {
				userID := uuid.New()
				m.users[userID] = User{ID: userID, Username: "testuser"}
			},
			userID:  uuid.New(),
			wantErr: false,
		},
		{
			name: "fails when repo returns error",
			setupRepo: func(m *mockAuthRepo) {
				m.getUserErr = fmt.Errorf("db error")
			},
			userID:  uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockAuthRepo()
			if tt.name == "gets user by ID successfully" {
				tt.userID = uuid.New()
				repo.users[tt.userID] = User{ID: tt.userID, Username: "testuser"}
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
		setupRepo func(*mockAuthRepo)
		wantLen   int
		wantErr   bool
	}{
		{
			name: "gets all users successfully",
			setupRepo: func(m *mockAuthRepo) {
				m.users[uuid.New()] = User{Username: "user1"}
				m.users[uuid.New()] = User{Username: "user2"}
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "fails when repo returns error",
			setupRepo: func(m *mockAuthRepo) {
				m.getUsersErr = fmt.Errorf("db error")
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockAuthRepo()
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
		setupRepo func(*mockAuthRepo, uuid.UUID)
		wantErr   bool
	}{
		{
			name: "gets user successfully",
			setupRepo: func(m *mockAuthRepo, id uuid.UUID) {
				m.users[id] = User{ID: id, Username: "testuser"}
			},
			wantErr: false,
		},
		{
			name: "fails when user not found",
			setupRepo: func(m *mockAuthRepo, id uuid.UUID) {
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockAuthRepo()
			userID := uuid.New()
			tt.setupRepo(repo, userID)
			svc := newTestService(repo)

			got, err := svc.GetUser(context.Background(), userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != userID {
				t.Errorf("GetUser() got ID = %v, want %v", got.ID, userID)
			}
		})
	}
}

func TestBaseServiceCreateUser(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(*mockAuthRepo)
		user      *User
		wantErr   bool
	}{
		{
			name: "creates user successfully",
			setupRepo: func(m *mockAuthRepo) {
			},
			user: &User{
				ID:       uuid.New(),
				Username: "newuser",
			},
			wantErr: false,
		},
		{
			name: "fails when repo returns error",
			setupRepo: func(m *mockAuthRepo) {
				m.createUserErr = fmt.Errorf("db error")
			},
			user: &User{
				ID:       uuid.New(),
				Username: "newuser",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockAuthRepo()
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
		setupRepo func(*mockAuthRepo, *User)
		wantErr   bool
	}{
		{
			name: "updates user successfully",
			setupRepo: func(m *mockAuthRepo, user *User) {
				m.users[user.ID] = *user
			},
			wantErr: false,
		},
		{
			name: "fails when repo returns error",
			setupRepo: func(m *mockAuthRepo, user *User) {
				m.users[user.ID] = *user
				m.updateUserErr = fmt.Errorf("db error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockAuthRepo()
			user := &User{
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
		setupRepo func(*mockAuthRepo, uuid.UUID)
		wantErr   bool
	}{
		{
			name: "deletes user successfully",
			setupRepo: func(m *mockAuthRepo, id uuid.UUID) {
				m.users[id] = User{ID: id}
			},
			wantErr: false,
		},
		{
			name: "fails when repo returns error",
			setupRepo: func(m *mockAuthRepo, id uuid.UUID) {
				m.users[id] = User{ID: id}
				m.deleteUserErr = fmt.Errorf("db error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockAuthRepo()
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
