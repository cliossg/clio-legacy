package fake

import (
	"context"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/feat/auth"
	"github.com/hermesgen/hm"
)

type AuthRepo struct {
	hm.Core

	GetUserByUsernameFn func(ctx context.Context, username string) (auth.User, error)
	GetUsersFn          func(ctx context.Context) ([]auth.User, error)
	GetUserFn           func(ctx context.Context, id uuid.UUID) (auth.User, error)
	CreateUserFn        func(ctx context.Context, user *auth.User) error
	UpdateUserFn        func(ctx context.Context, user *auth.User) error
	DeleteUserFn        func(ctx context.Context, id uuid.UUID) error

	GetUserByUsernameCalls []struct {
		Ctx      context.Context
		Username string
	}
	GetUsersCalls []struct {
		Ctx context.Context
	}
	GetUserCalls []struct {
		Ctx context.Context
		ID  uuid.UUID
	}
	CreateUserCalls []struct {
		Ctx  context.Context
		User *auth.User
	}
	UpdateUserCalls []struct {
		Ctx  context.Context
		User *auth.User
	}
	DeleteUserCalls []struct {
		Ctx context.Context
		ID  uuid.UUID
	}

	users map[uuid.UUID]auth.User
}

func NewAuthRepo() *AuthRepo {
	cfg := hm.NewConfig()
	return &AuthRepo{
		Core:  hm.NewCore("fake-auth-repo", hm.XParams{Cfg: cfg}),
		users: make(map[uuid.UUID]auth.User),
	}
}

func (f *AuthRepo) Query() *hm.QueryManager {
	return nil
}

func (f *AuthRepo) BeginTx(ctx context.Context) (context.Context, hm.Tx, error) {
	return ctx, nil, nil
}

func (f *AuthRepo) GetUserByUsername(ctx context.Context, username string) (auth.User, error) {
	f.GetUserByUsernameCalls = append(f.GetUserByUsernameCalls, struct {
		Ctx      context.Context
		Username string
	}{Ctx: ctx, Username: username})

	if f.GetUserByUsernameFn != nil {
		return f.GetUserByUsernameFn(ctx, username)
	}

	for _, user := range f.users {
		if user.Username == username {
			return user, nil
		}
	}
	return auth.User{}, nil
}

func (f *AuthRepo) GetUsers(ctx context.Context) ([]auth.User, error) {
	f.GetUsersCalls = append(f.GetUsersCalls, struct {
		Ctx context.Context
	}{Ctx: ctx})

	if f.GetUsersFn != nil {
		return f.GetUsersFn(ctx)
	}

	var users []auth.User
	for _, user := range f.users {
		users = append(users, user)
	}
	return users, nil
}

func (f *AuthRepo) GetUser(ctx context.Context, id uuid.UUID) (auth.User, error) {
	f.GetUserCalls = append(f.GetUserCalls, struct {
		Ctx context.Context
		ID  uuid.UUID
	}{Ctx: ctx, ID: id})

	if f.GetUserFn != nil {
		return f.GetUserFn(ctx, id)
	}

	if user, ok := f.users[id]; ok {
		return user, nil
	}
	return auth.User{}, nil
}

func (f *AuthRepo) CreateUser(ctx context.Context, user *auth.User) error {
	f.CreateUserCalls = append(f.CreateUserCalls, struct {
		Ctx  context.Context
		User *auth.User
	}{Ctx: ctx, User: user})

	if f.CreateUserFn != nil {
		return f.CreateUserFn(ctx, user)
	}

	f.users[user.ID] = *user
	return nil
}

func (f *AuthRepo) UpdateUser(ctx context.Context, user *auth.User) error {
	f.UpdateUserCalls = append(f.UpdateUserCalls, struct {
		Ctx  context.Context
		User *auth.User
	}{Ctx: ctx, User: user})

	if f.UpdateUserFn != nil {
		return f.UpdateUserFn(ctx, user)
	}

	f.users[user.ID] = *user
	return nil
}

func (f *AuthRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	f.DeleteUserCalls = append(f.DeleteUserCalls, struct {
		Ctx context.Context
		ID  uuid.UUID
	}{Ctx: ctx, ID: id})

	if f.DeleteUserFn != nil {
		return f.DeleteUserFn(ctx, id)
	}

	delete(f.users, id)
	return nil
}
