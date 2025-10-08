package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/feat/auth"
)

var (
	featAuth = "auth"
	resUser  = "user"
)

func (repo *ClioRepo) GetUsers(ctx context.Context) ([]auth.User, error) {
	query, err := repo.Query().Get(featAuth, resUser, "GetAll")
	if err != nil {
		return nil, err
	}

	var users []auth.User
	err = repo.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *ClioRepo) GetUser(ctx context.Context, id uuid.UUID) (auth.User, error) {
	query, err := repo.Query().Get(featAuth, resUser, "Get")
	if err != nil {
		return auth.User{}, err
	}

	var user auth.User
	err = repo.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return auth.User{}, err
	}

	return user, nil
}

func (repo *ClioRepo) GetUserByUsername(ctx context.Context, username string) (auth.User, error) {
	query, err := repo.Query().Get(featAuth, resUser, "GetByUsername")
	if err != nil {
		return auth.User{}, err
	}

	var user auth.User
	err = repo.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.User{}, errors.New("user not found")
		}
		return auth.User{}, err
	}

	return user, nil
}

func (repo *ClioRepo) CreateUser(ctx context.Context, user *auth.User) error {
	query, err := repo.Query().Get(featAuth, resUser, "Create")
	if err != nil {
		return err
	}

	_, err = repo.db.NamedExecContext(ctx, query, user)
	return err
}

func (repo *ClioRepo) UpdateUser(ctx context.Context, user *auth.User) error {
	query, err := repo.Query().Get(featAuth, resUser, "Update")
	if err != nil {
		return err
	}

	_, err = repo.db.NamedExecContext(ctx, query, user)
	return err
}

func (repo *ClioRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query().Get(featAuth, resUser, "Delete")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, id)
	return err
}
