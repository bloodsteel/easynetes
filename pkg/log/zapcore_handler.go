package log

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	toLevel = map[zapcore.Level]Level{
		zapcore.FatalLevel: FatalLevel,
		zapcore.ErrorLevel: ErrorLevel,
		zapcore.WarnLevel:  WarnLevel,
		zapcore.InfoLevel:  InfoLevel,
		zapcore.DebugLevel: DebugLevel,
	}
	toZapLevel = map[Level]zapcore.Level{
		FatalLevel: zapcore.FatalLevel,
		ErrorLevel: zapcore.ErrorLevel,
		WarnLevel:  zapcore.WarnLevel,
		InfoLevel:  zapcore.InfoLevel,
		DebugLevel: zapcore.DebugLevel,
	}
)

func init() {
	registerDefaultHandler(ZapLogHandlerCallbackFunc)
}

// ZapLogHandlerCallbackFunc 是一个handler函数, 用来模拟日志输出; 最终会被Configure 中的对象给截获
func ZapLogHandlerCallbackFunc(
	level Level,
	scope *Scope,
	msg string) {
	var fields []zapcore.Field
	// 是否是JSON
	if useJSON.Load().(bool) {
		for _, k := range scope.labelKeys {
			v := scope.labels[k]
			fields = append(fields, zap.Field{
				Key:       k,
				Type:      zapcore.ReflectType,
				Interface: v,
			})
		}
	} else {
		sb := &strings.Builder{}
		sb.WriteString(msg)
		if len(scope.labelKeys) > 0 {
			sb.WriteString("\t")
		}
		space := false
		for _, k := range scope.labelKeys {
			if space {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("%s=%v", k, scope.labels[k]))
			space = true
		}
		msg = sb.String()
	}
	emit(scope, toZapLevel[level], msg, fields)
}

// callerSkipOffset 用来确定调用方函数位置, emit -> ZapLogHandlerCallbackFunc -> registerDefaultHandler -> Info*
const callerSkipOffset = 4

func dumpStack(level zapcore.Level, scope *Scope) bool {
	thresh := toLevel[level]
	if scope != defaultScope { // 有疑问
		thresh = ErrorLevel
		switch level {
		case zapcore.FatalLevel:
			thresh = FatalLevel
		}
	}
	return scope.GetStackTraceLevel() >= thresh
}

func emit(scope *Scope, level zapcore.Level, msg string, fields []zapcore.Field) {
	e := zapcore.Entry{
		Message:    msg,
		Level:      level,
		Time:       time.Now(),
		LoggerName: scope.nameToEmit,
	}

	if scope.GetLogCallers() {
		e.Caller = zapcore.NewEntryCaller(runtime.Caller(scope.callerSkip + callerSkipOffset))
	}

	if dumpStack(level, scope) {
		e.Stack = zap.Stack("").String
	}

	pt := funcs.Load().(patchTable)
	if pt.write != nil {
		if err := pt.write(e, fields); err != nil {
			_, _ = fmt.Fprintf(pt.errorSink, "%v log write error: %v\n", time.Now(), err)
			_ = pt.errorSink.Sync()
		}
	}
}
