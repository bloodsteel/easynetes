package db

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Connect to a database and with a ping
// return a DB object
func Connect(driver, dsn string, openConns, idleConns int, lifeTime time.Duration) (*sqlx.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	// 数据库最大连接数
	db.SetMaxOpenConns(openConns)
	// 最大空闲连接数
	db.SetMaxIdleConns(idleConns)
	// 连接最长存活时间, 超过这个时间连接将不再复用
	db.SetConnMaxLifetime(lifeTime)
	if err := pingDatabase(db); err != nil {
		return nil, err
	}

	return sqlx.NewDb(db, driver), nil
}

// ping database to ensure a connection can be established
func pingDatabase(db *sql.DB) (err error) {
	for i := 0; i < 5; i++ {
		err = db.Ping()
		if err == nil {
			return
		}
		time.Sleep(time.Millisecond * 500)
	}
	return
}
