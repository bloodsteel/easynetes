package db

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/bloodsteel/easynetes/pkg/config"
)

type Repository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewRepository(db *gorm.DB, logger *slog.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

var (
	gormDB *gorm.DB
	sqlDB  *sql.DB
)

// NewGormDB 获取 gormDB
func NewGormDB() *gorm.DB {
	return gormDB
}

// NewMysqlDB 初始化 mysql 连接
func NewMysqlDB(config *config.DB) (gormDB *gorm.DB, sqlDB *sql.DB, err error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Mysql.User,
		config.Mysql.Password,
		config.Mysql.Host,
		config.Mysql.DatabaseName,
	)
	mysqlConfig := mysql.Config{
		DSN:                       dataSourceName,
		DefaultStringSize:         255,   // string类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用datetime精度, MySQL5.6之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式, MySQL 5.7之前的数据库和MariaDB不支持重命名索引
		DontSupportRenameColumn:   true,  // 用`change`重命名列, MySQL8之前的数据库和MariaDB不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前MySQL版本自动配置
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	gormConfig := gorm.Config{
		// Logger: NewLogger().LogMode(ParseLevel(config.LogLevel)),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.TablePrefix,
			SingularTable: true,
		},
		Logger: newLogger,
	}
	gormDB, err = gorm.Open(mysql.New(mysqlConfig), &gormConfig)
	if err != nil {
		return
	}
	sqlDB, err = gormDB.DB()
	if err != nil {
		return
	}

	// 连接池配置
	sqlDB.SetMaxIdleConns(config.Mysql.MaxIdleConns) //  用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxOpenConns(config.Mysql.MaxOpenConns) //  设置打开数据库连接的最大数量。
	sqlDB.SetConnMaxLifetime(time.Hour)              //  设置了连接可复用的最大时间。

	return
}
