package signal

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// WithContextFunc 返回给定ctx 的 副本, 该ctx 会在接收到操作系统的 中断信号 的时候 关闭channel
// 回调函数是在 cancel channel 之前调用的
func WithContextFunc(ctx context.Context, f func()) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
		defer signal.Stop(c)

		select {
		case <-ctx.Done():
			log.Printf("%v ctx.Done() received", time.Now().Local().String())
		case <-c:
			f()
			cancel()
		}
	}()

	return ctx
}
