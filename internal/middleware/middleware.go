package middleware

import (
	"fmt"
	"io"
)

var (
	dunno = []byte("???")
)

// 使用context.WithValue将一个value存入http request ctx中的key
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "easynetes/middleware request context value " + k.name
}

// buffer Write
func bW(w io.Writer, s string, args ...interface{}) {
	fmt.Fprintf(w, s, args...)
}
