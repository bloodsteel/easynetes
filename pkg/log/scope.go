package log

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

//
// Scope 支持使用WithLabels进行结构化日志记录:
//
//   s := RegisterScope("MyScope", "Description", 0)
//   s = s.WithLabels("foo", "bar", "baz", 123, "qux", 0.123)
//   s.Info("Hello")                      // <time>   info   MyScope   Hello  foo=bar baz=123 qux=0.123
//

// Scope 将日志控制限制在命名范围级别
type Scope struct {
	// 不可变, 在创建的时候指定
	name        string
	nameToEmit  string
	description string
	callerSkip  int

	// 通过专门的方法进行配置, 可以动态调整
	outputLevel     atomic.Value
	stackTraceLevel atomic.Value
	logCallers      atomic.Value

	// 存放 labels 数据 - key slice 用来保证顺序
	labelKeys []string
	labels    map[string]interface{}
}

var (
	scopes = make(map[string]*Scope) // 用来存放所有注册的scope
	lock   sync.RWMutex

	// 用来存放 记录日志 使用的回调函数
	defaultHandlers   []scopeHandlerCallbackFunc
	defaultHandlersMu sync.RWMutex
)

// scopeHandlerCallbackFunc 表示一种函数类型, 用来在调用Fatal* Error* Warn& Info* 和 Debug* 的时候都使用该回调函数进行日志记录
type scopeHandlerCallbackFunc func(
	level Level,
	scope *Scope,
	msg string)

// registerDefaultHandler 注册一个scope handler, 默认所有的 scope 都会调用该回调函数
func registerDefaultHandler(callback scopeHandlerCallbackFunc) {
	defaultHandlersMu.Lock()
	defer defaultHandlersMu.Unlock()
	defaultHandlers = append(defaultHandlers, callback)
}

// RegisterScope 注册一个新的 日志记录scope. 是并发安全的, 如果 名称一样 则已经存在的 scope 会被返回
// Scope 的名称不允许包含 `: , .`
func RegisterScope(name string, description string, callerSkip int) *Scope {
	if strings.ContainsAny(name, ":,.") {
		panic(fmt.Sprintf("scope name %s is invalid, it cannot contain colons, commas, or periods", name))
	}

	lock.Lock()
	defer lock.Unlock()

	s, ok := scopes[name]
	if !ok {
		s = &Scope{
			name:        name,
			description: description,
			callerSkip:  callerSkip,
		}
		// 每个scope对象的默认设置
		s.SetOutputLevel(InfoLevel)
		s.SetStackTraceLevel(NoneLevel)
		s.SetLogCallers(false)

		s.nameToEmit = name

		scopes[name] = s
	}

	s.labels = make(map[string]interface{})

	return s
}

// FindScope 返回已经注册的scope, 如果没有找到, 返回 nil
func FindScope(scope string) *Scope {
	lock.RLock()
	defer lock.RUnlock()

	s := scopes[scope]
	return s
}

// Scopes 返回当前定义的所有scopes集合的快照
func Scopes() map[string]*Scope {
	lock.RLock()
	defer lock.RUnlock()

	s := make(map[string]*Scope, len(scopes))
	for k, v := range scopes {
		s[k] = v
	}

	return s
}

// Fatal 使用 fmt.Sprint 在 fatal level 构造并记录日志
func (s *Scope) Fatal(args ...interface{}) {
	if s.GetOutputLevel() >= FatalLevel {
		s.callHandlers(FatalLevel, s, fmt.Sprint(args...))
	}
}

// Fatalf 使用 fmt.Sprintf  在 fatal level 构造并记录日志
func (s *Scope) Fatalf(args ...interface{}) {
	if s.GetOutputLevel() >= FatalLevel {
		msg := fmt.Sprint(args[0])
		s.callHandlers(FatalLevel, s, fmt.Sprintf(msg, args[1:]...))
	}
}

// FatalEnabled 返回 是否 启用 该 scope的 fatal level
func (s *Scope) FatalEnabled() bool {
	return s.GetOutputLevel() >= FatalLevel
}

// Error 使用 fmt.Sprint 在 error level 构造并记录日志
func (s *Scope) Error(args ...interface{}) {
	if s.GetOutputLevel() >= ErrorLevel {
		s.callHandlers(ErrorLevel, s, fmt.Sprint(args...))
	}
}

// Errorf 使用 fmt.Sprintf  在 error level 构造并记录日志
func (s *Scope) Errorf(args ...interface{}) {
	if s.GetOutputLevel() >= ErrorLevel {
		msg := fmt.Sprint(args[0])
		s.callHandlers(ErrorLevel, s, fmt.Sprintf(msg, args[1:]...))
	}
}

// ErrorEnabled 返回 是否 启用 该 scope的 error level
func (s *Scope) ErrorEnabled() bool {
	return s.GetOutputLevel() >= ErrorLevel
}

// Warn 使用 fmt.Sprint 在 warn level 构造并记录日志
func (s *Scope) Warn(args ...interface{}) {
	if s.GetOutputLevel() >= WarnLevel {
		s.callHandlers(WarnLevel, s, fmt.Sprint(args...))
	}
}

// Warnf 使用 fmt.Sprintf  在 warn level 构造并记录日志
func (s *Scope) Warnf(args ...interface{}) {
	if s.GetOutputLevel() >= WarnLevel {
		msg := fmt.Sprint(args[0])
		s.callHandlers(WarnLevel, s, fmt.Sprintf(msg, args[1:]...))
	}
}

// WarnEnabled 返回 是否 启用 该 scope的 warn level
func (s *Scope) WarnEnabled() bool {
	return s.GetOutputLevel() >= WarnLevel
}

// Info 使用 fmt.Sprint 在 info level 构造并记录日志
func (s *Scope) Info(args ...interface{}) {
	if s.GetOutputLevel() >= InfoLevel {
		s.callHandlers(InfoLevel, s, fmt.Sprint(args...))
	}
}

// Infof 使用 fmt.Sprintf  在 info level 构造并记录日志
func (s *Scope) Infof(args ...interface{}) {
	if s.GetOutputLevel() >= InfoLevel {
		msg := fmt.Sprint(args[0])
		s.callHandlers(InfoLevel, s, fmt.Sprintf(msg, args[1:]...))
	}
}

// InfoEnabled 返回 是否 启用 该 scope的 info level
func (s *Scope) InfoEnabled() bool {
	return s.GetOutputLevel() >= InfoLevel
}

// Debug 使用 fmt.Sprint 在 debug level 构造并记录日志
func (s *Scope) Debug(args ...interface{}) {
	if s.GetOutputLevel() >= DebugLevel {
		s.callHandlers(DebugLevel, s, fmt.Sprint(args...))
	}
}

// Debugf 使用 fmt.Sprintf  在 debug level 构造并记录日志
func (s *Scope) Debugf(args ...interface{}) {
	if s.GetOutputLevel() >= DebugLevel {
		msg := fmt.Sprint(args[0])
		s.callHandlers(DebugLevel, s, fmt.Sprintf(msg, args[1:]...))
	}
}

// DebugEnabled 返回 是否 启用 该 scope的 debug level
func (s *Scope) DebugEnabled() bool {
	return s.GetOutputLevel() >= DebugLevel
}

// Name 返回这个 scope 的 name
func (s *Scope) Name() string {
	return s.name
}

// Description 返回这个scope 的 description
func (s *Scope) Description() string {
	return s.description
}

// SetOutputLevel 调整 该 scope 的 日志 输出级别
func (s *Scope) SetOutputLevel(l Level) {
	s.outputLevel.Store(l)
}

// GetOutputLevel 返回 该 scope 的 日志 输出级别
func (s *Scope) GetOutputLevel() Level {
	return s.outputLevel.Load().(Level)
}

// SetStackTraceLevel 调整 该 scope 的 日志 stacktrace级别
func (s *Scope) SetStackTraceLevel(l Level) {
	s.stackTraceLevel.Store(l)
}

// GetStackTraceLevel 返回 该 scope 的 日志 stacktrace级别
func (s *Scope) GetStackTraceLevel() Level {
	return s.stackTraceLevel.Load().(Level)
}

// SetLogCallers 调整 该 scope 的 日志 logCaller
func (s *Scope) SetLogCallers(logCallers bool) {
	s.logCallers.Store(logCallers)
}

// GetLogCallers 返回 该 scope 的 日志 logCaller
func (s *Scope) GetLogCallers() bool {
	return s.logCallers.Load().(bool)
}

// copy 创建 s 的 拷贝并返回指向它的指针
func (s *Scope) copy() *Scope {
	out := *s
	out.labels = copyStringInterfaceMap(s.labels)
	return &out
}

// WithLabels 添加 key-value 到 scope 的 lables 中. key 必须是字符串. 返回 添加了 labels 的 s 的 副本
// e.g. newScope := oldScope.WithLabels("foo", "bar", "baz", 123, "qux", 0.123)
func (s *Scope) WithLabels(kvlist ...interface{}) *Scope {
	out := s.copy()
	if len(kvlist)%2 != 0 {
		out.labels["WithLabels error"] = fmt.Sprintf("even number of parameters required, got %d", len(kvlist))
		return out
	}

	for i := 0; i < len(kvlist); i += 2 {
		keyi := kvlist[i]
		key, ok := keyi.(string)
		if !ok {
			out.labels["WithLabels error"] = fmt.Sprintf("label name %v must be a string, got %T ", keyi, keyi)
			return out
		}
		out.labels[key] = kvlist[i+1]
		out.labelKeys = append(out.labelKeys, key)
	}
	return out
}

// callHandlers 调用注册到当前s对象的所有的handler
func (s *Scope) callHandlers(
	severity Level,
	scope *Scope,
	msg string) {
	defaultHandlersMu.RLock()
	defer defaultHandlersMu.RUnlock()
	for _, h := range defaultHandlers {
		h(severity, scope, msg)
	}
}

func copyStringInterfaceMap(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
