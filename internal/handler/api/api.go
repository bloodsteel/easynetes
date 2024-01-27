package api

import (
	"net/http"

	"github.com/bloodsteel/easynetes/internal/core"
	"github.com/bloodsteel/easynetes/internal/handler/api/host"
	"github.com/bloodsteel/easynetes/internal/handler/api/user"
	"github.com/bloodsteel/easynetes/internal/middleware"

	"github.com/go-chi/chi/v5"
)

// ProvideAPI is a Wire provider
// returns a api server include router handlers
func ProvideAPI(
	userDao core.UserDao,
	hostDao core.HostInstanceDao,
) *Server {
	return &Server{
		userDao: userDao,
		hostDao: hostDao,
	}
}

// Server payload
type Server struct {
	userDao core.UserDao
	hostDao core.HostInstanceDao
}

// Handler http router for api
func (s Server) Handler() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// 用户管理相关的APIs
	router.Route("/users", func(r chi.Router) {
		r.With(middleware.Paginate).Get("/", user.ListUsers(s.userDao))
		r.Post("/", user.CreateUser(s.userDao))

		// 针对单用户的路由
		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", user.HandlerUser(s.userDao))
			r.Put("/", user.HandlerUser(s.userDao))
			r.Delete("/", user.HandlerUser(s.userDao))
		})
	})

	// 资产相关的APIs
	router.Route("/cmdb", func(r chi.Router) {
		// 主机数据路由
		r.Route("/host", func(r chi.Router) {
			r.With(middleware.Paginate).Get("/", host.ListHosts(s.hostDao))
			r.Post("/", host.CreateHost(s.hostDao))

			r.Route("/{hostID}", func(r chi.Router) {
				r.Get("/", host.HandlerHost(s.hostDao))
				r.Put("/", host.HandlerHost(s.hostDao))
				r.Delete("/", host.HandlerHost(s.hostDao))
			})
		})
		// 可用区数据路由
		r.Route("/azone", func(r chi.Router) {})
	})

	return router
}
