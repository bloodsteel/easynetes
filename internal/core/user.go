package core

import (
	"context"
	"time"
)

type (
	// User 表示一个用户对象
	User struct {
		ID         int64      `db:"id" json:"id"`
		UserName   string     `db:"user_name" json:"user_name"`
		UserEmail  string     `db:"user_email" json:"user_email"`
		UserPwd    string     `db:"user_pwd" json:"-"`
		UserPhone  string     `db:"user_phone" json:"user_phone"`
		IsAdmin    bool       `db:"is_admin" json:"is_admin"`
		LastLogin  *time.Time `db:"last_login" json:"last_login"`
		CreateTime time.Time  `db:"create_time" json:"create_time"`
		UpdateTime time.Time  `db:"update_time" json:"update_time"`
		Remark     string     `db:"remark" json:"remark"`
	}

	// UserDao 定义了一组从数据库操作用户表的一系列操作
	UserDao interface {
		// Get 根据ID从数据库中获取用户对象
		Get(context.Context, int64) (*User, error)
		// List 从数据库中获取一组用户对象
		List(context.Context) ([]*User, error)
		// Create 在数据库中创建一个用户对象
		Create(context.Context, *User) (int64, error)
		// Update 更新数据库中已经存在的一个用户
		Update(context.Context, *User) (*User, error)
		// Delete 从数据库中删除一个已经存在的用户
		Delete(context.Context, int64) error
	}
)
