package cmd

import (
	"fmt"
	"net/url"

	"github.com/bloodsteel/easynetes/internal/dao/host"
	"github.com/bloodsteel/easynetes/internal/dao/user"
	"github.com/bloodsteel/easynetes/pkg/config"
	"github.com/bloodsteel/easynetes/pkg/db"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
)

var daoSet = wire.NewSet(
	// db.NewRepository,
	provideDatabase,
	user.ProvideUserDao,
	host.ProvideHostDao,
)

// provideDatabase is a Wire provider
// returns a database connection pool
func provideDatabase(cfg *config.Config) (*sqlx.DB, error) {
	var dataSourceName string
	var tz = url.QueryEscape("Asia/Shanghai")
	var dbDriver = "mysql"
	dataSourceName = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?timeout=120s&charset=utf8mb4,utf8&parseTime=true&loc=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DatabaseName,
		tz,
	)
	return db.Connect(
		dbDriver,
		dataSourceName,
		cfg.Database.MaxOpenConns,
		cfg.Database.MaxIdleConns,
		cfg.Database.ConnMaxLifetime,
	)
}
