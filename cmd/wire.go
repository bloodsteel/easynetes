//

//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/bloodsteel/easynetes/pkg/config"
	"github.com/bloodsteel/easynetes/pkg/server"
	"github.com/google/wire"
)

// InitializeApplication 初始化app, 用来启动http server
func InitializeApplication(c *config.Config) (*server.Server, error) {
	wire.Build(
		daoSet,
		serverSet,
	)
	return &server.Server{}, nil
}
