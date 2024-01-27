package utils

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/bloodsteel/easynetes/internal/middleware"
	"go.uber.org/zap"
)

// JSONResult 封装response的数据结构
/*
status 可以取值的范围是: success、fail、error
code 是非常规的http_code; 通常每个code都对应一个message信息; code方便前端和后端的代码编写和错误追踪等
message 通常用来装载人类可读的信息; 比如 "token没有找到;请重新登录认证"
*/
type (
	JSONResult struct {
		RequestID string      `json:"request_id"`
		Status    string      `json:"status"`
		Code      int64       `json:"code" `
		Msg       string      `json:"message"`
		Data      interface{} `json:"data"`
	}
)

func renderJSON(writer http.ResponseWriter, request *http.Request, status, code string, payload interface{}) {
	ctx := request.Context()
	requestID := middleware.GetRequestIDFromCtx(ctx)
	writer.Header().Set("Content-Type", "application/json")

	codes := strings.Split(code, "-")
	httpCode, _ := strconv.Atoi(codes[0])
	statusCode, _ := strconv.ParseInt(codes[1], 10, 64)
	writer.WriteHeader(httpCode)

	enc := json.NewEncoder(writer)
	_ = enc.Encode(&JSONResult{
		RequestID: requestID,
		Status:    status,
		Msg:       Msg[code],
		Code:      statusCode,
		Data:      payload,
	})
}

// RenderSuccess 返回成功数据
func RenderSuccess(writer http.ResponseWriter, request *http.Request, payload interface{}) {
	ctx := request.Context()
	requestID := middleware.GetRequestIDFromCtx(ctx)
	zap.L().Named("default").WithOptions(zap.AddCallerSkip(0)).Info("",
		zap.String("request_id", requestID),
		zap.Any("payload", payload),
	)
	renderJSON(writer, request, StatusSuccess, SCodeOK, payload)
}

// RenderFail 返回失败数据
func RenderFail(writer http.ResponseWriter, request *http.Request, code string) {
	ctx := request.Context()
	requestID := middleware.GetRequestIDFromCtx(ctx)
	zap.L().Named("default").WithOptions(zap.AddCallerSkip(0)).Warn("",
		zap.String("request_id", requestID),
		zap.String("code", code),
	)
	renderJSON(writer, request, StatusFail, code, nil)
}

// RenderError 返回错误数据
func RenderError(writer http.ResponseWriter, request *http.Request, code string, err error) {
	ctx := request.Context()
	requestID := middleware.GetRequestIDFromCtx(ctx)
	if err == nil {
		RenderFail(writer, request, SCodeUnknow)
		return
	}
	if err == sql.ErrNoRows {
		RenderFail(writer, request, SCodeNotFoundWithDao)
		return
	}
	zap.L().Named("default").WithOptions(zap.AddCallerSkip(0)).Error("",
		zap.String("request_id", requestID),
		zap.NamedError("error", err),
	)
	renderJSON(writer, request, StatusError, code, err.Error())
}
