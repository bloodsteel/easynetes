package health

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// ProvideHealth is a Wire provider
// returns a health server include router handlers
func ProvideHealth() *Server {
	return &Server{}
}

// Server payload
type Server struct {
}

// Handler http router for health
func (s Server) Handler() http.Handler {
	router := chi.NewRouter()

	router.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(time.Millisecond * 500)
		writer.WriteHeader(200)
		writer.Header().Set("Content-Type", "text/plain")
		_, _ = writer.Write([]byte("ok"))
	})

	return router
}
