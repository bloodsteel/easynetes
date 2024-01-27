package log

import (
	"fmt"

	"github.com/go-logr/logr"
)

// 在klog中级别超过3则认为是debug , 0~3 认为是 info 级别, 并且 没有 warn 级别
// 同时klog 的命令行 参数 -v 一定要大于3 才可以
const debugLevelThreshold = 3

// 实现 logr.LogSink 接口
type zapLogger struct {
	l      *Scope // 原本表示日志记录器, 在这里scope就表示日志记录器
	lvlSet bool
}

// Init 方法在使用 logr.New 方法的时候被调用, 用来设置 callerSkip, 我们这里的scope不需要这样设置, 这里什么都不做
func (zl *zapLogger) Init(ri logr.RuntimeInfo) {
	zl.l.SetLogCallers(false) // 注意: 截获之后 logCaller 会错乱, 所以不开启
	return
}

func (zl *zapLogger) Enabled(lvl int) bool {
	if zl.lvlSet && lvl > debugLevelThreshold {
		return zl.l.DebugEnabled()
	}
	return zl.l.InfoEnabled()
}

func (zl *zapLogger) Info(lvl int, msg string, keysAndVals ...interface{}) {
	if zl.lvlSet && lvl > debugLevelThreshold {
		zl.l.WithLabels(keysAndVals...).Debug(trimNewline(msg))
	} else {
		zl.l.WithLabels(keysAndVals...).Info(trimNewline(msg))
	}
}

func (zl *zapLogger) Error(err error, msg string, keysAndVals ...interface{}) {
	if zl.l.ErrorEnabled() {
		if err == nil {
			zl.l.WithLabels(keysAndVals...).Error(trimNewline(msg))
		} else {
			zl.l.WithLabels(keysAndVals...).Error(fmt.Sprintf("%v: %s", err.Error(), trimNewline(msg)))
		}
	}
}

func (zl *zapLogger) WithValues(keysAndValues ...interface{}) logr.LogSink {
	newLogger := *zl
	newLogger.l.WithLabels(keysAndValues...)
	return &newLogger
}

// 我们的scope name不允许改变, 这里就什么都不做了
func (zl *zapLogger) WithName(name string) logr.LogSink {
	return zl
}

// 这个是为了实现 logr.CallDepthLogSink 接口, 由于我们的scope是不允许在运行时改变callerSkip, 这里就什么都不做了
func (zl *zapLogger) WithCallDepth(depth int) logr.LogSink {
	return zl
}

// NewLogrAdapter 创建一个新的 logr.Logger 便于使用给定的zap日志记录器进行 日志记录
func NewLogrAdapter(l *Scope) logr.Logger {
	zl := &zapLogger{
		l:      l,
		lvlSet: true,
	}
	return logr.New(zl)
}

// 将日志中的换行符去掉, zap编码器中自带换行符了
func trimNewline(msg string) string {
	if len(msg) == 0 {
		return msg
	}
	lc := len(msg) - 1
	if msg[lc] == '\n' {
		return msg[:lc]
	}
	return msg
}
