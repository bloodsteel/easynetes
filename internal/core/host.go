package core

import (
	"context"
	"time"
)

type (
	// HostInstance 主机实例(宿主机、云主机、虚拟机)
	HostInstance struct {
		ID            int64     `db:"id" json:"id"`
		InstanceID    string    `db:"instance_id" json:"instance_id"`
		HostName      string    `db:"host_name" json:"host_name"`
		CPUCores      int8      `db:"cpu_cores" json:"cpu_cores"`
		CPUSockets    int8      `db:"cpu_sockets" json:"cpu_sockets"`
		MemSize       int       `db:"mem_size" json:"mem_size"`
		OSName        string    `db:"os_name" json:"os_name"`
		KernelVersion string    `db:"kernel_version" json:"kernel_version"`
		ConnPort      int       `db:"conn_port" json:"conn_port"`
		HostStatus    int       `db:"host_status" json:"host_status"`
		HostType      int       `db:"host_type" json:"host_type"`
		CreateTime    time.Time `db:"create_time" json:"create_time"`
		UpdateTime    time.Time `db:"update_time" json:"update_time"`
		Remark        string    `db:"remark" json:"remark"`
	}
	// HostInstanceDao 定义了一组从数据库操作主机实例的一系列操作
	HostInstanceDao interface {
		// Get 根据ID从数据库中获取主机实例对象
		Get(context.Context, int64) (*HostInstance, error)
		// List 从数据库中获取一组主机实例对象
		List(context.Context, map[string]interface{}) ([]*HostInstance, error)
		// Create 在数据库中创建一个主机实例对象
		Create(context.Context, *HostInstance) (int64, error)
		// Update 更新数据库中已经存在的一个主机实例
		Update(context.Context, *HostInstance) (*HostInstance, error)
		// Delete 从数据库中删除一个已经存在的主机实例
		Delete(context.Context, int64) error
	}
)
