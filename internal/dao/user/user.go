package user

import (
	"context"

	"github.com/bloodsteel/easynetes/internal/core"
	"github.com/jmoiron/sqlx"
)

func ProvideUserDao(db *sqlx.DB) core.UserDao {
	return &userDao{db: db}
}

type userDao struct {
	db *sqlx.DB
}

var _ core.UserDao = &userDao{}

func (user *userDao) Get(ctx context.Context, in int64) (*core.User, error) {
	return nil, nil
}

func (user *userDao) List(ctx context.Context) ([]*core.User, error) {
	return nil, nil
}

func (user *userDao) Create(ctx context.Context, in *core.User) (int64, error) {
	return 1, nil
}

func (user *userDao) Update(ctx context.Context, in *core.User) (*core.User, error) {
	return nil, nil
}

func (user *userDao) Delete(ctx context.Context, in int64) error {
	return nil
}
