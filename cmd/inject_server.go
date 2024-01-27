package cmd

import (
	"net/http"

	"github.com/bloodsteel/easynetes/internal/handler/api"
	"github.com/bloodsteel/easynetes/internal/handler/health"
	"github.com/bloodsteel/easynetes/pkg/config"
	"github.com/bloodsteel/easynetes/pkg/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/wire"
)

var (
	serverSet = wire.NewSet(
		api.ProvideAPI,
		health.ProvideHealth,
		ProvideRouter,
		ProvideServer,
	)
)

// ProvideRouter is a Wire provider
func ProvideRouter(api *api.Server, health *health.Server) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Mount("/api", api.Handler())
	// mount sub router
	router.Mount("/health", health.Handler())
	return router
}

// ProvideServer is a Wire provider
// returns an http server
func ProvideServer(handler http.Handler, cfg *config.Config) *server.Server {
	return server.NewServer(
		handler,
		cfg.Server.BindAddress,
		cfg.Server.InsecurePort,
		cfg.Server.SecurePort,
		cfg.Server.SecureCert,
		cfg.Server.SecureKey,
		cfg.Server.SecureHost,
		cfg.Server.ReadTimeout,
		cfg.Server.WriteTimeout,
		cfg.Server.TerminationTimeout,
	)
}
