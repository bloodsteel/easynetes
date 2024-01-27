package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// Logger 请求日志 TODO: 使用自定义的Logger重写
func Logger(next http.Handler) http.Handler {
	return middleware.Logger(next)
}
