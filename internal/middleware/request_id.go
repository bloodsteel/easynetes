package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// ReqIDCtxKey Context keys for request id
var (
	ReqIDCtxKey = &contextKey{"xRequestID"}
)

// RequestIDHeader 表示应该从哪个HTTP Header中获取Request ID
var RequestIDHeader = "X-Request-Id"

// RequestID 返回一个中间件, 用来在每个request 的 context 中注入 request ID
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(write http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		// Get Request ID from HTTP Header
		requestID := request.Header.Clone().Get(RequestIDHeader)
		if requestID == "" {
			uuidStr := uuid.New()
			requestID = uuidStr.String()
		}
		ctx = context.WithValue(ctx, ReqIDCtxKey, requestID)
		request = request.WithContext(ctx)
		next.ServeHTTP(write, request)
	})
}

// GetRequestIDFromCtx 从 http request ctx 中获取 request id, 方便调用
func GetRequestIDFromCtx(ctx context.Context) string {
	requestID, _ := ctx.Value(ReqIDCtxKey).(string)
	return requestID
}
