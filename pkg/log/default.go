package log

func registerDefaultScope() *Scope {
	return RegisterScope(DefaultScopeName, "Unscoped logging messages.", 0)
}

var defaultScope = registerDefaultScope()

// Fatal 使用 fmt.Sprint 在 fatal level 构造并记录日志
func Fatal(fields ...interface{}) {
	defaultScope.Fatal(fields...)
}

// Fatalf 使用 fmt.Sprintf  在 fatal level 构造并记录日志
func Fatalf(args ...interface{}) {
	defaultScope.Fatalf(args...)
}

// FatalEnabled 返回 是否 启用 该 scope的 fatal level
func FatalEnabled() bool {
	return defaultScope.FatalEnabled()
}

// Error 使用 fmt.Sprint 在 error level 构造并记录日志
func Error(fields ...interface{}) {
	defaultScope.Error(fields...)
}

// Errorf 使用 fmt.Sprintf  在 error level 构造并记录日志
func Errorf(args ...interface{}) {
	defaultScope.Errorf(args...)
}

// ErrorEnabled 返回 是否 启用 该 scope的 error level
func ErrorEnabled() bool {
	return defaultScope.ErrorEnabled()
}

// Warn 使用 fmt.Sprint 在 warn level 构造并记录日志
func Warn(fields ...interface{}) {
	defaultScope.Warn(fields...)
}

// Warnf 使用 fmt.Sprintf  在 warn level 构造并记录日志
func Warnf(args ...interface{}) {
	defaultScope.Warnf(args...)
}

// WarnEnabled 返回 是否 启用 该 scope的 warn level
func WarnEnabled() bool {
	return defaultScope.WarnEnabled()
}

// Info 使用 fmt.Sprint 在 info level 构造并记录日志
func Info(fields ...interface{}) {
	defaultScope.Info(fields...)
}

// Infof 使用 fmt.Sprintf  在 info level 构造并记录日志
func Infof(args ...interface{}) {
	defaultScope.Infof(args...)
}

// InfoEnabled 返回 是否 启用 该 scope的 info level
func InfoEnabled() bool {
	return defaultScope.InfoEnabled()
}

// Debug 使用 fmt.Sprint 在 debug level 构造并记录日志
func Debug(fields ...interface{}) {
	defaultScope.Debug(fields...)
}

// Debugf 使用 fmt.Sprintf  在 debug level 构造并记录日志
func Debugf(args ...interface{}) {
	defaultScope.Debugf(args...)
}

// DebugEnabled 返回 是否 启用 该 scope的 debug level
func DebugEnabled() bool {
	return defaultScope.DebugEnabled()
}

// WithLabels 添加 key-value 到 scope 的 lables 中. key 必须是字符串. 返回 添加了 labels 的 s 的 副本
func WithLabels(kvlist ...interface{}) *Scope {
	return defaultScope.WithLabels(kvlist...)
}
