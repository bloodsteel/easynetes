package host

import (
	"context"

	"github.com/bloodsteel/easynetes/internal/core"
	"github.com/jmoiron/sqlx"
)

func ProvideHostDao(db *sqlx.DB) core.HostInstanceDao {
	return &hostDao{db: db}
}

type hostDao struct {
	db *sqlx.DB
}

var _ core.HostInstanceDao = &hostDao{}

func (host *hostDao) Get(ctx context.Context, in int64) (*core.HostInstance, error) {
	return nil, nil
}

func (host *hostDao) List(ctx context.Context, in map[string]interface{}) ([]*core.HostInstance, error) {
	return nil, nil
}

func (host *hostDao) Create(ctx context.Context, in *core.HostInstance) (int64, error) {
	return 1, nil
}

func (host *hostDao) Update(ctx context.Context, in *core.HostInstance) (*core.HostInstance, error) {
	return nil, nil
}

func (host *hostDao) Delete(ctx context.Context, in int64) error {
	return nil
}
