package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
)

const defaultLimit int64 = 50

// OffsetCtxKey Context keys for list offset
// LimitCtxKey Context keys for list limit
var (
	OffsetCtxKey = &contextKey{"OffsetKey"}
	LimitCtxKey  = &contextKey{"LimitKey"}
	PageCtxKey   = &contextKey{"PageKey"}
)

type (
	// Response 自定义 Response, 用于实现 http.ResponseWriter 接口
	Response struct {
		http.ResponseWriter
		buf *bytes.Buffer // 用于存放返回的数据
	}

	// ResponseData 真正返回的数据
	ResponseData struct {
		RequestID string      `json:"request_id"`
		Status    string      `json:"status"`
		Code      int64       `json:"code"`
		Msg       string      `json:"message"`
		Data      interface{} `json:"data"`
		Page      int64       `json:"page"`
		Next      int64       `json:"next"`
		Prev      int64       `json:"prev"`
		PageSize  int64       `json:"page_size"`
		Count     int64       `json:"count"`
	}
)

// Write 重写了 http.ResponseWriter 的 Write 方法
// 将用户返回的数据 写入到 buf 中
func (resp *Response) Write(p []byte) (int, error) {
	return resp.buf.Write(p)
}

// Paginate 分页功能
func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(write http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		var (
			page     = request.URL.Query().Get("page")
			pageSize = request.URL.Query().Get("page_size")
		)
		// hold page num
		pageCtxValue := 1
		// 计算 offset & limit
		offset, _ := strconv.Atoi(page)
		limit, _ := strconv.Atoi(pageSize)
		if limit < 1 || limit > 200 {
			limit = int(defaultLimit)
		}
		write.Header().Set("Page-Size", strconv.Itoa(limit))
		switch offset {
		case 0, 1:
			offset = 0
			write.Header().Set("Page", "1")
			ctx = context.WithValue(ctx, PageCtxKey, pageCtxValue)
		default:
			offset = (offset - 1) * limit
			write.Header().Set("Page", page)
			pageCtxValue, _ = strconv.Atoi(page)
			ctx = context.WithValue(ctx, PageCtxKey, pageCtxValue)
		}
		ctx = context.WithValue(ctx, LimitCtxKey, limit)
		ctx = context.WithValue(ctx, OffsetCtxKey, offset)

		// 封装自定义的Response, 用于捕获返回的数据, 添加分页信息
		wWriter := &Response{
			ResponseWriter: write,
			buf:            bytes.NewBufferString(""),
		}
		next.ServeHTTP(wWriter, request.WithContext(ctx))

		// paginate process
		nextN, prevN := getNextAndPrev(pageCtxValue, limit, int(getHeaderCount(wWriter)))
		respData := &ResponseData{
			PageSize: int64(limit),
			Page:     int64(pageCtxValue),
			Next:     nextN,
			Prev:     prevN,
			Count:    getHeaderCount(wWriter),
		}
		// 用于从 buf 中获取 用户返回的数据 并写入到 真正的 http.ResponseWriter
		_ = json.NewDecoder(wWriter.buf).Decode(respData)
		_ = json.NewEncoder(wWriter.ResponseWriter).Encode(respData)
	})
}

// GetPaginateFromCtx 从http request ctx中获取limit & offset, 方便调用
func GetPaginateFromCtx(ctx context.Context) (limit, offset, page int64) {
	if ctx == nil {
		return defaultLimit, 0, 1
	}
	if v, ok := ctx.Value(LimitCtxKey).(int); ok {
		limit = int64(v)
	} else {
		limit = defaultLimit
	}
	if v, ok := ctx.Value(OffsetCtxKey).(int); ok {
		offset = int64(v)
	} else {
		offset = 0
	}
	if v, ok := ctx.Value(PageCtxKey).(int); ok {
		page = int64(v)
	} else {
		page = 1
	}
	return
}

// getHeaderCount 从Header头获取Count
func getHeaderCount(response *Response) int64 {
	count, _ := strconv.Atoi(response.ResponseWriter.Header().Get("Count"))
	return int64(count)
}

// getNextAndPrev 获取上一页和下一页
func getNextAndPrev(page, pageSize, count int) (next, prev int64) {
	pageCount := int(math.Ceil(float64(count) / float64(pageSize)))
	if page > pageCount {
		next = int64(0)
		prev = int64(pageCount)
	} else {
		if page == pageCount {
			next = int64(0)
		} else {
			next = int64(page + 1)
		}
		prev = int64(page - 1)
	}
	return
}
