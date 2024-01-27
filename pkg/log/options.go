package log

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// 默认值
const (
	DefaultScopeName          = "default"
	OverrideScopeName         = "all"
	defaultOutputLevel        = InfoLevel
	defaultStackTraceLevel    = NoneLevel
	defaultOutputPath         = "stdout"
	defaultErrorOutputPath    = "stderr"
	defaultRotationMaxAge     = 30
	defaultRotationMaxSize    = 100 * 1024 * 1024
	defaultRotationMaxBackups = 1000
)

// Level 支持的所有日志级别的枚举
type Level int

const (
	// NoneLevel 禁用日志记录
	NoneLevel Level = iota
	// FatalLevel 启用 fatal 级别日志记录
	FatalLevel
	// ErrorLevel 启用 error 级别日志记录
	ErrorLevel
	// WarnLevel 启用 warn 级别日志记录
	WarnLevel
	// InfoLevel 启用 info 级别日志记录
	InfoLevel
	// DebugLevel 启用 debug 级别日志记录
	DebugLevel
)

var levelToString = map[Level]string{
	DebugLevel: "debug",
	InfoLevel:  "info",
	WarnLevel:  "warn",
	ErrorLevel: "error",
	FatalLevel: "fatal",
	NoneLevel:  "none",
}

var stringToLevel = map[string]Level{
	"debug": DebugLevel,
	"info":  InfoLevel,
	"warn":  WarnLevel,
	"error": ErrorLevel,
	"fatal": FatalLevel,
	"none":  NoneLevel,
}

// Options 定义了日志库支持的命令行选项的集合
type Options struct {
	// OutputPaths 表示要将日志数据写入的文件系统路径的列表. 默认为stdout;
	OutputPaths []string

	// ErrorOutputPaths 表示要将错误日志写入的文件系统路径的列表
	ErrorOutputPaths []string

	// RotateOutputPath 表示轮转日志文件的路径. 默认不轮转. 该文件会根据轮转参数的定义随着时间的推移自动轮转
	RotateOutputPath string

	// RotationMaxSize 表示日志文件轮转前的最大大小, 默认为100兆
	RotationMaxSize int

	// RotationMaxAge 旧日志保存的最大天数, 默认删除超过30天的日志文件
	RotationMaxAge int

	// RotationMaxBackups 旧日志保留的最大文件数, 默认最多保留1000个
	RotationMaxBackups int

	// JSONEncoding 控制日志是否被格式化为JSON
	JSONEncoding bool

	// LogGrpc 表示是否捕获Grpc日志, 默认为true
	LogGrpc bool

	outputLevels     string
	logCallers       string
	stackTraceLevels string
}

// DefaultOptions 使用默认值创建一个*Options实例
func DefaultOptions() *Options {
	return &Options{
		OutputPaths:        []string{defaultOutputPath},
		ErrorOutputPaths:   []string{defaultErrorOutputPath},
		RotationMaxSize:    defaultRotationMaxSize,
		RotationMaxAge:     defaultRotationMaxAge,
		RotationMaxBackups: defaultRotationMaxBackups,
		outputLevels:       DefaultScopeName + ":" + levelToString[defaultOutputLevel],
		stackTraceLevels:   DefaultScopeName + ":" + levelToString[defaultStackTraceLevel],
		LogGrpc:            false,
	}
}

// SetOutputLevel 为给定的scope设置给定的日志输出级别, 目的是重新定义命令行中给定的scope的日志级别;
func (o *Options) SetOutputLevel(scope string, level Level) {
	sl := scope + ":" + levelToString[level]
	levels := strings.Split(o.outputLevels, ",")

	if scope == DefaultScopeName {
		// 如果任何一个记录没有scope前缀, 这表示这是default scope
		for i, ol := range levels {
			if !strings.Contains(ol, ":") {
				levels[i] = sl
				o.outputLevels = strings.Join(levels, ",")
				return
			}
		}
	}

	prefix := scope + ":"
	for i, ol := range levels {
		if strings.HasPrefix(ol, prefix) {
			levels[i] = sl
			o.outputLevels = strings.Join(levels, ",")
			return
		}
	}

	levels = append(levels, sl)
	o.outputLevels = strings.Join(levels, ",")
}

// SetStackTraceLevel 目的是用于为给定的scope设置日志堆栈跟踪(stack tracing)级别
func (o *Options) SetStackTraceLevel(scope string, level Level) {
	sl := scope + ":" + levelToString[level]
	levels := strings.Split(o.stackTraceLevels, ",")

	if scope == DefaultScopeName {
		for i, ol := range levels {
			if !strings.Contains(ol, ":") {
				levels[i] = sl
				o.stackTraceLevels = strings.Join(levels, ",")
				return
			}
		}
	}

	prefix := scope + ":"
	for i, ol := range levels {
		if strings.HasPrefix(ol, prefix) {
			levels[i] = sl
			o.stackTraceLevels = strings.Join(levels, ",")
			return
		}
	}

	levels = append(levels, sl)
	o.stackTraceLevels = strings.Join(levels, ",")
}

// SetLogCallers 设置是否为给定的scope开启logCaller(在日志中输出源码位置)
func (o *Options) SetLogCallers(scope string, include bool) {
	scopes := strings.Split(o.logCallers, ",")

	// remove any occurrence of the scope
	for i, s := range scopes {
		if s == scope {
			scopes[i] = ""
		}
	}

	if include {
		// find a free slot if there is one
		for i, s := range scopes {
			if s == "" {
				scopes[i] = scope
				o.logCallers = strings.Join(scopes, ",")
				return
			}
		}

		scopes = append(scopes, scope)
	}

	o.logCallers = strings.Join(scopes, ",")
}

// GetOutputLevel 返回给定的scope的日志输出级别
func (o *Options) GetOutputLevel(scope string) (Level, error) {
	levels := strings.Split(o.outputLevels, ",")

	if scope == DefaultScopeName {
		for _, ol := range levels {
			if !strings.Contains(ol, ":") {
				_, l, err := convertScopedLevel(ol)
				return l, err
			}
		}
	}

	prefix := scope + ":"
	for _, ol := range levels {
		if strings.HasPrefix(ol, prefix) {
			_, l, err := convertScopedLevel(ol)
			return l, err
		}
	}

	return NoneLevel, fmt.Errorf("no level defined for scope '%s'", scope)
}

// GetStackTraceLevel 返回给定的scope的日志堆栈跟踪级别
func (o *Options) GetStackTraceLevel(scope string) (Level, error) {
	levels := strings.Split(o.stackTraceLevels, ",")

	if scope == DefaultScopeName {
		// see if we have an entry without a scope prefix (which represents the default scope)
		for _, ol := range levels {
			if !strings.Contains(ol, ":") {
				_, l, err := convertScopedLevel(ol)
				return l, err
			}
		}
	}

	prefix := scope + ":"
	for _, ol := range levels {
		if strings.HasPrefix(ol, prefix) {
			_, l, err := convertScopedLevel(ol)
			return l, err
		}
	}

	return NoneLevel, fmt.Errorf("no level defined for scope '%s'", scope)
}

// GetLogCallers 返回是否为给定的scope开启logCaller
func (o *Options) GetLogCallers(scope string) bool {
	scopes := strings.Split(o.logCallers, ",")

	for _, s := range scopes {
		if s == scope {
			return true
		}
	}

	return false
}

func convertScopedLevel(sl string) (string, Level, error) {
	var s string
	var l string

	pieces := strings.Split(sl, ":")
	if len(pieces) == 1 {
		s = DefaultScopeName
		l = pieces[0]
	} else if len(pieces) == 2 {
		s = pieces[0]
		l = pieces[1]
	} else {
		return "", NoneLevel, fmt.Errorf("invalid output level format '%s'", sl)
	}

	level, ok := stringToLevel[l]
	if !ok {
		return "", NoneLevel, fmt.Errorf("invalid output level '%s'", sl)
	}

	return s, level, nil
}

// AttachCobraFlags 为给定的Cobra command附加一组Cobra flags, 方便通过命令行控制日志选项
func (o *Options) AttachCobraFlags(cmd *cobra.Command) {
	o.AttachFlags(
		cmd.PersistentFlags().StringArrayVar,
		cmd.PersistentFlags().StringVar,
		cmd.PersistentFlags().IntVar,
		cmd.PersistentFlags().BoolVar)
}

// AttachFlags 添加flag
func (o *Options) AttachFlags(
	stringArrayVar func(p *[]string, name string, value []string, usage string),
	stringVar func(p *string, name string, value string, usage string),
	intVar func(p *int, name string, value int, usage string),
	boolVar func(p *bool, name string, value bool, usage string)) {
	stringArrayVar(&o.OutputPaths, "log_target", o.OutputPaths,
		"The set of paths where to output the log. This can be any path as well as the special values stdout and stderr")

	stringVar(&o.RotateOutputPath, "log_rotate", o.RotateOutputPath,
		"The path for the optional rotating log file")

	intVar(&o.RotationMaxAge, "log_rotate_max_age", o.RotationMaxAge,
		"The maximum age in days of a log file beyond which the file is rotated (0 indicates no limit)")

	intVar(&o.RotationMaxSize, "log_rotate_max_size", o.RotationMaxSize,
		"The maximum size in megabytes of a log file beyond which the file is rotated")

	intVar(&o.RotationMaxBackups, "log_rotate_max_backups", o.RotationMaxBackups,
		"The maximum number of log file backups to keep before older files are deleted (0 indicates no limit)")

	boolVar(&o.JSONEncoding, "log_as_json", o.JSONEncoding,
		"Whether to format output as JSON or in plain console-friendly format")

	levelListString := fmt.Sprintf("[%s, %s, %s, %s, %s, %s]",
		levelToString[DebugLevel],
		levelToString[InfoLevel],
		levelToString[WarnLevel],
		levelToString[ErrorLevel],
		levelToString[FatalLevel],
		levelToString[NoneLevel])

	allScopes := Scopes()
	if len(allScopes) > 1 {
		keys := make([]string, 0, len(allScopes))
		for name := range allScopes {
			keys = append(keys, name)
		}
		keys = append(keys, OverrideScopeName)
		sort.Strings(keys)
		s := strings.Join(keys, ", ")

		stringVar(&o.outputLevels, "log_output_level", o.outputLevels,
			fmt.Sprintf("Comma-separated minimum per-scope logging level of messages to output, in the form of "+
				"<scope>:<level>,<scope>:<level>,... where scope can be one of [%s] and level can be one of %s",
				s, levelListString))

		stringVar(&o.stackTraceLevels, "log_stacktrace_level", o.stackTraceLevels,
			fmt.Sprintf("Comma-separated minimum per-scope logging level at which stack traces are captured, in the form of "+
				"<scope>:<level>,<scope:level>,... where scope can be one of [%s] and level can be one of %s",
				s, levelListString))

		stringVar(&o.logCallers, "log_caller", o.logCallers,
			fmt.Sprintf("Comma-separated list of scopes for which to include caller information, scopes can be any of [%s]", s))
	} else {
		stringVar(&o.outputLevels, "log_output_level", o.outputLevels,
			fmt.Sprintf("The minimum logging level of messages to output,  can be one of %s",
				levelListString))

		stringVar(&o.stackTraceLevels, "log_stacktrace_level", o.stackTraceLevels,
			fmt.Sprintf("The minimum logging level at which stack traces are captured, can be one of %s",
				levelListString))

		stringVar(&o.logCallers, "log_caller", o.logCallers,
			"Comma-separated list of scopes for which to include called information, scopes can be any of [default]")
	}
}
