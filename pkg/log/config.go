package log

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc/grpclog"
	"gopkg.in/natefinch/lumberjack.v2"
	"k8s.io/klog/v2"
)

const (
	// none 用于禁用日志输出以及禁用堆栈跟踪
	none zapcore.Level = 100
	// GrpcScopeName grpc scope 名称
	GrpcScopeName string = "grpc"
)

var levelToZap = map[Level]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
	FatalLevel: zapcore.FatalLevel,
	NoneLevel:  none,
}

var defaultEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "@time",
	LevelKey:       "@level",
	NameKey:        "@scope",
	CallerKey:      "@caller",
	MessageKey:     "@msg",
	StacktraceKey:  "@stack",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeTime:     formatDate,
}

// 控制 日志的 写入  sync  写入报错等动作, 被Info* Error* 中的回调函数 调用
// 同时也为了方便在单测中 改写 这一行为
type patchTable struct {
	write       func(ent zapcore.Entry, fields []zapcore.Field) error
	sync        func() error
	exitProcess func(code int)
	errorSink   zapcore.WriteSyncer
	close       func() error
}

var (
	// 没搞懂 为什么是 指针
	funcs = &atomic.Value{}
	// 用于控制, 是否是JSON格式输出, 这种方式的查询效率要高于查询结构体内部字段, 而且更加方便
	useJSON atomic.Value
	logGrpc bool
)

func init() {
	// 使用默认的options启动
	_ = Configure(DefaultOptions())
}

func formatDate(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	t = t.UTC()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	micros := t.Nanosecond() / 1000

	buf := make([]byte, 27)

	buf[0] = byte((year/1000)%10) + '0'
	buf[1] = byte((year/100)%10) + '0'
	buf[2] = byte((year/10)%10) + '0'
	buf[3] = byte(year%10) + '0'
	buf[4] = '-'
	buf[5] = byte((month)/10) + '0'
	buf[6] = byte((month)%10) + '0'
	buf[7] = '-'
	buf[8] = byte((day)/10) + '0'
	buf[9] = byte((day)%10) + '0'
	buf[10] = 'T'
	buf[11] = byte((hour)/10) + '0'
	buf[12] = byte((hour)%10) + '0'
	buf[13] = ':'
	buf[14] = byte((minute)/10) + '0'
	buf[15] = byte((minute)%10) + '0'
	buf[16] = ':'
	buf[17] = byte((second)/10) + '0'
	buf[18] = byte((second)%10) + '0'
	buf[19] = '.'
	buf[20] = byte((micros/100000)%10) + '0'
	buf[21] = byte((micros/10000)%10) + '0'
	buf[22] = byte((micros/1000)%10) + '0'
	buf[23] = byte((micros/100)%10) + '0'
	buf[24] = byte((micros/10)%10) + '0'
	buf[25] = byte((micros)%10) + '0'
	buf[26] = 'Z'

	enc.AppendString(string(buf))
}

// prepZap 是后续函数Configure使用的一个实用函数; 用于生成 zapcore.Core 实例
func prepZap(options *Options) (zapcore.Core, zapcore.Core, zapcore.WriteSyncer, error) {
	var enc zapcore.Encoder // 编码器

	encCfg := defaultEncoderConfig

	if options.JSONEncoding {
		enc = zapcore.NewJSONEncoder(encCfg)
		useJSON.Store(true)
	} else {
		enc = zapcore.NewConsoleEncoder(encCfg)
		useJSON.Store(false)
	}
	// 日志轮转
	var rotaterSink zapcore.WriteSyncer
	if options.RotateOutputPath != "" {
		rotaterSink = zapcore.AddSync(&lumberjack.Logger{
			Filename:   options.RotateOutputPath,
			MaxSize:    options.RotationMaxSize,
			MaxBackups: options.RotationMaxBackups,
			MaxAge:     options.RotationMaxAge,
			Compress:   true,
			LocalTime:  true,
		})
	}
	// 错误日志, 注意这里的Open方法
	errSink, closeErrorSink, err := zap.Open(options.ErrorOutputPaths...)
	if err != nil {
		return nil, nil, nil, err
	}
	// 常规输出日志
	var outputSink zapcore.WriteSyncer
	if len(options.OutputPaths) > 0 {
		outputSink, _, err = zap.Open(options.OutputPaths...)
		if err != nil {
			// 关掉上一步打开的 errSink
			closeErrorSink()
			return nil, nil, nil, err
		}
	}
	// 合并 常规日志输出sink和日志轮转sink
	var sink zapcore.WriteSyncer
	if rotaterSink != nil && outputSink != nil {
		sink = zapcore.NewMultiWriteSyncer(outputSink, rotaterSink)
	} else if rotaterSink != nil {
		sink = rotaterSink
	} else {
		sink = outputSink
	}
	// zapcore.LevelEnabler 接口实现
	var enabler zap.LevelEnablerFunc = func(lvl zapcore.Level) bool {
		switch lvl {
		case zapcore.ErrorLevel:
			return defaultScope.ErrorEnabled()
		case zapcore.WarnLevel:
			return defaultScope.WarnEnabled()
		case zapcore.InfoLevel:
			return defaultScope.InfoEnabled()
		}
		return defaultScope.DebugEnabled()
	}
	// 创建 core 对象, 第一个core是xxx, 第二个core是用于捕获日志
	return zapcore.NewCore(enc, sink, zap.NewAtomicLevelAt(zapcore.DebugLevel)),
		zapcore.NewCore(enc, sink, enabler),
		errSink, nil
}

func updateScopes(options *Options) error {
	allScopes := Scopes()

	// 更新所有的scopes的 日志输出级别
	if err := processLevels(allScopes, options.outputLevels, func(s *Scope, l Level) { s.SetOutputLevel(l) }); err != nil {
		return err
	}

	// 更新所有的scopes的 日志堆栈级别
	if err := processLevels(allScopes, options.stackTraceLevels, func(s *Scope, l Level) { s.SetStackTraceLevel(l) }); err != nil {
		return err
	}

	// 更新所有的scopes的 logCaller
	sc := strings.Split(options.logCallers, ",")
	for _, s := range sc {
		if s == "" {
			continue
		}

		if s == OverrideScopeName {
			// 忽略所有其他的内容, 直接使用覆盖值
			for _, scope := range allScopes {
				scope.SetLogCallers(true)
			}

			return nil
		}

		if scope, ok := allScopes[s]; ok {
			scope.SetLogCallers(true)
		} else {
			return fmt.Errorf("unknown scope '%s' specified", s)
		}
	}

	// 如果有必要, 则更新 LogGrpc
	if logGrpc {
		options.LogGrpc = true
	}

	return nil
}

// processLevels 将 传递进来的 args 分解为一组 scope和level, 然后将分解后的结果应用到对应的scope对象上
// 注意: 支持使用 all 进行覆盖
func processLevels(allScopes map[string]*Scope, arg string, setter func(*Scope, Level)) error {
	levels := strings.Split(arg, ",")
	for _, sl := range levels {
		s, l, err := convertScopedLevel(sl)
		if err != nil {
			return err
		}

		if scope, ok := allScopes[s]; ok {
			setter(scope, l)
		} else if s == OverrideScopeName {
			// 直接覆盖
			for _, scope := range allScopes {
				setter(scope, l)
			}
			return nil
		} else if s == GrpcScopeName {
			grpcScope := RegisterScope(GrpcScopeName, "", 3)
			logGrpc = true
			setter(grpcScope, l)
			return nil
		} else {
			return fmt.Errorf("unknown scope '%s' specified", s)
		}
	}

	return nil
}

// Configure 初始化日志子系统; 通常在进程启动的时候调用一次
func Configure(options *Options) error {
	core, captureCore, errSink, err := prepZap(options)
	if err != nil {
		return err
	}

	// 将命令行配置的设置 应用到 scope 对象
	if err = updateScopes(options); err != nil {
		return err
	}

	// 暂时没用到, 如果要打印日志到kafka, 可以使用
	closeFns := make([]func() error, 0)

	// 非常重要! 注意这里的core
	pt := patchTable{
		write: func(ent zapcore.Entry, fields []zapcore.Field) error {
			err := core.Write(ent, fields)
			if ent.Level == zapcore.FatalLevel {
				funcs.Load().(patchTable).exitProcess(1)
			}

			return err
		},
		sync:        core.Sync,
		exitProcess: os.Exit,
		errorSink:   errSink,
		close: func() error {
			core.Sync() // nolint: errcheck
			for _, f := range closeFns {
				if err := f(); err != nil {
					return err
				}
			}
			return nil
		},
	}
	funcs.Store(pt)

	// 定制zap.Logger需要的各种opts
	opts := []zap.Option{
		zap.ErrorOutput(errSink),
		zap.AddCallerSkip(1),
	}
	if defaultScope.GetLogCallers() {
		opts = append(opts, zap.AddCaller()) // 是否启用logCaller
	}
	l := defaultScope.GetStackTraceLevel()
	if l != NoneLevel {
		opts = append(opts, zap.AddStacktrace(levelToZap[l]))
	}

	// 非常重要! 注意这里的captureCore
	captureLogger := zap.New(captureCore, opts...)

	// 捕获全局的zap日志记录 并且强制通过我们的logger
	_ = zap.ReplaceGlobals(captureLogger)

	// 捕获golang标准库"log"的日志输出 并强制通过我们的logger
	_ = zap.RedirectStdLog(captureLogger)

	// 捕获 gRPC 日志
	if options.LogGrpc {
		grpclog.SetLogger(zapgrpc.NewLogger(captureLogger.WithOptions(zap.AddCallerSkip(3))))
	}

	// 捕获 klog (Kubernetes logging) 通过我们的logger , 注意这里是单例模式
	configureKlog.Do(func() {
		klog.SetLogger(NewLogrAdapter(KlogScope))
	})
	if klogVerbose() {
		KlogScope.SetOutputLevel(DebugLevel)
	}

	return nil
}

// Sync 是刷新所有缓存的日志条目, 通常在程序退出的之前使用; 是 pt.sync() 的快捷方法
func Sync() error {
	return funcs.Load().(patchTable).sync()
}

// Close 是 调用 pt.close() 的快捷方法
func Close() error {
	return funcs.Load().(patchTable).close()
}
