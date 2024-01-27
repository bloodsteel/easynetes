package log

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/grpclog"
	"k8s.io/klog/v2"
)

const timePattern = "[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9].[0-9][0-9][0-9][0-9][0-9][0-9]Z"

func resetGlobals() {
	scopes = make(map[string]*Scope, 1)
	defaultScope = registerDefaultScope()
}

func testOptions() *Options {
	return DefaultOptions()
}

func captureStdout(f func()) ([]string, error) {
	tf, err := ioutil.TempFile("", "log_test")
	if err != nil {
		return nil, err
	}

	old := os.Stdout
	os.Stdout = tf

	f()

	os.Stdout = old
	path := tf.Name()
	_ = tf.Sync()
	_ = tf.Close()

	content, err := ioutil.ReadFile(path)
	_ = os.Remove(path)

	if err != nil {
		return nil, err
	}

	return strings.Split(string(content), "\n"), nil
}

// 用来测试在 命令行中给定的 all scope的 设置会覆盖其他 scope的 设置
// 在 config.go 中的 processLevels 函数中实现
func TestOverrides(t *testing.T) {
	resetGlobals()
	s := RegisterScope("TestOverrides", "For testing", 0)

	o := DefaultOptions()
	o.outputLevels = "default:debug,all:info"
	if err := Configure(o); err != nil {
		t.Errorf("Expecting success, got %v", err)
	} else if s.GetOutputLevel() != InfoLevel {
		t.Errorf("Expecting InfoLevel, got %v", s.GetOutputLevel())
	} else {
		t.Logf("Test success, got %v", levelToString[s.GetOutputLevel()])
	}

	o = DefaultOptions()
	o.stackTraceLevels = "default:debug,all:info"
	if err := Configure(o); err != nil {
		t.Errorf("Expecting success, got %v", err)
	} else if s.GetStackTraceLevel() != InfoLevel {
		t.Errorf("Expecting InfoLevel, got %v", s.GetStackTraceLevel())
	} else {
		t.Logf("Test success, got %v", levelToString[s.GetStackTraceLevel()])
	}

	o = DefaultOptions()
	o.logCallers = "all"
	if err := Configure(o); err != nil {
		t.Errorf("Expecting success, got %v", err)
	} else if !s.GetLogCallers() {
		t.Error("Expecting true, got false")
	} else {
		t.Logf("Test success, got %v", s.GetLogCallers())
	}
}

// 确保 rotation 设置正确
func TestRotateNoStdout(t *testing.T) {
	dir, _ := ioutil.TempDir("", "TestRotateNoStdout")
	defer os.RemoveAll(dir)

	file := dir + "/rot.log"

	o := DefaultOptions()
	o.OutputPaths = []string{} // 关闭stdout
	o.RotateOutputPath = file
	if err := Configure(o); err != nil {
		t.Fatalf("Unable to configure logging: %v", err)
	}

	defaultScope.Error("HELLO")
	Sync()

	content, err := ioutil.ReadFile(file)
	if err != nil {
		t.Errorf("Got failure '%v', expecting success", err)
	}

	lines := strings.Split(string(content), "\n")
	if !strings.Contains(lines[0], "HELLO") {
		t.Errorf("Expecting for first line of log to contain HELLO, got %s", lines[0])
	} else {
		t.Logf("Test success, got: %v", lines[0])
	}
}

// 测试 rotate 和 stdout 同时设置
func TestRotateAndStdout(t *testing.T) {
	dir, _ := ioutil.TempDir("", "TestRotateAndStdout")
	defer os.RemoveAll(dir)

	file := dir + "/rot.log"

	stdoutLines, _ := captureStdout(func() {
		o := DefaultOptions()     // 默认是stdout 开启的
		o.RotateOutputPath = file // 使用 rotate
		if err := Configure(o); err != nil {
			t.Fatalf("Unable to configure logger: %v", err)
		}

		defaultScope.Error("HELLO")
		Sync()

		content, err := ioutil.ReadFile(file)
		if err != nil {
			t.Errorf("Got failure '%v', expecting success", err)
		}

		rotLines := strings.Split(string(content), "\n")
		if !strings.Contains(rotLines[0], "HELLO") {
			t.Errorf("Expecting for first line of log to contain HELLO, got %s", rotLines[0])
		} else {
			t.Logf("Test success, got: %v", rotLines[0])
		}
	})

	if !strings.Contains(stdoutLines[0], "HELLO") {
		t.Errorf("Expecting for first line of log to contain HELLO, got %s", stdoutLines[0])
	} else {
		t.Logf("Test success, got: %v", stdoutLines[0])
	}
}

// 测试 FindScope 函数
func TestFind(t *testing.T) {
	resetGlobals()
	if z := FindScope("TestFind"); z != nil {
		t.Errorf("Found scope %v, but expected it wouldn't exist", z)
	}

	_ = RegisterScope("TestFind", "", 0)

	if z := FindScope("TestFind"); z == nil {
		t.Error("Did not find scope, expected to find it")
	}
}

// 测试scope name, 反例
func TestBadNames(t *testing.T) {
	resetGlobals()
	badNames := []string{
		"a:b",
		"a,b",
		"a.b",

		":ab",
		",ab",
		".ab",

		"ab:",
		"ab,",
		"ab.",
	}

	for _, name := range badNames {
		tryBadName(t, name)
	}
}

func tryBadName(t *testing.T, name string) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
		t.Errorf("Expecting to panic when using bad scope name %s, but didn't", name)
	}()

	_ = RegisterScope(name, "A poorly named scope", 0) // name 不合规 会 panic
}

// 测试 RegisterScope
func TestMultipleScopesWithSameName(t *testing.T) {
	resetGlobals()
	z1 := RegisterScope("zzzz", "z", 0)
	z2 := RegisterScope("zzzz", "z", 0)

	if z1 != z2 {
		t.Error("Expecting the same scope objects, got different ones")
	}
}

// 测试写入失败的时候 errSink 生效
// 这里并不是一个有效的测试单例, 要想测试成功, 需要捕获 stderr 的数据
func TestBadWriter(t *testing.T) {
	o := testOptions()
	if err := Configure(o); err != nil {
		t.Errorf("Got err '%v', expecting success", err)
	}

	// 修改 writer 为 直接 返回错误, 并不进行真正的写入
	// 此时 任何 scope 的日志记录, 最终都会调用 handler 回调函数, 在回调函数中, 如果写入失败, 会将错误写入到 errSink中
	pt := funcs.Load().(patchTable)
	pt.write = func(zapcore.Entry, []zapcore.Field) error {
		return errors.New("bad")
	}
	funcs.Store(pt)

	defaultScope.Error("TestBadWriter") // 这一句是不会输出的, 因为已经把 write 给修改了
}

// 测试defaultScope的日志输出, 包含json, logCaller, stack
func TestDefaultScopeLogging(t *testing.T) {
	cases := []struct {
		f          func() // 测试用例
		pat        string // 期望输出正则匹配
		json       bool   // 输出为json
		caller     bool   // 是否添加logCaller
		wantExit   bool   // 期望结果是否会panic退出
		stackLevel Level  // 堆栈跟踪级别
	}{
		{
			f:   func() { Debug("Hello") },
			pat: timePattern + "\tdebug\tdefault\tHello",
		},
		{
			f:   func() { Debugf("Hello") },
			pat: timePattern + "\tdebug\tdefault\tHello",
		},
		{
			f:   func() { Debugf("%s", "Hello") },
			pat: timePattern + "\tdebug\tdefault\tHello",
		},

		{
			f:   func() { Info("Hello") },
			pat: timePattern + "\tinfo\tdefault\tHello",
		},
		{
			f:   func() { Infof("Hello") },
			pat: timePattern + "\tinfo\tdefault\tHello",
		},
		{
			f:   func() { Infof("%s", "Hello") },
			pat: timePattern + "\tinfo\tdefault\tHello",
		},

		{
			f:   func() { Warn("Hello") },
			pat: timePattern + "\twarn\tdefault\tHello",
		},
		{
			f:   func() { Warnf("Hello") },
			pat: timePattern + "\twarn\tdefault\tHello",
		},
		{
			f:   func() { Warnf("%s", "Hello") },
			pat: timePattern + "\twarn\tdefault\tHello",
		},

		{
			f:   func() { Error("Hello") },
			pat: timePattern + "\terror\tdefault\tHello",
		},
		{
			f:   func() { Errorf("Hello") },
			pat: timePattern + "\terror\tdefault\tHello",
		},
		{
			f:   func() { Errorf("%s", "Hello") },
			pat: timePattern + "\terror\tdefault\tHello",
		},

		{
			f:        func() { Fatal("Hello") },
			pat:      timePattern + "\tfatal\tdefault\tHello",
			wantExit: true,
		},
		{
			f:        func() { Fatalf("Hello") },
			pat:      timePattern + "\tfatal\tdefault\tHello",
			wantExit: true,
		},
		{
			f:        func() { Fatalf("%s", "Hello") },
			pat:      timePattern + "\tfatal\tdefault\tHello",
			wantExit: true,
		},

		{
			f:      func() { Debug("Hello") },
			pat:    timePattern + "\tdebug\tdefault\tlog/.*_test.go:.*\tHello",
			caller: true,
		},

		{
			f: func() { WithLabels("foo", "bar").WithLabels("baz", 123, "qux", 0.123).Debug("Hello") },
			pat: "{\"@level\":\"debug\",\"@time\":\"" + timePattern + "\",\"@scope\":\"default\",\"@caller\":\"log/.*_test.go:.*\",\"@msg\":\"Hello\"," + "\"foo\":\"bar\",\"baz\":123,\"qux\":0.123," +
				"\"@stack\":\".*\"}",
			json:       true,
			caller:     true,
			stackLevel: DebugLevel,
		},
		{
			f: func() { WithLabels("foo", "bar").WithLabels("baz", 123, "qux", 0.123).Info("Hello") },
			pat: "{\"@level\":\"info\",\"@time\":\"" + timePattern + "\",\"@scope\":\"default\",\"@caller\":\"log/.*_test.go:.*\",\"@msg\":\"Hello\"," + "\"foo\":\"bar\",\"baz\":123,\"qux\":0.123," +
				"\"@stack\":\".*\"}",
			json:       true,
			caller:     true,
			stackLevel: DebugLevel,
		},
		{
			f: func() { WithLabels("foo", "bar").WithLabels("baz", 123, "qux", 0.123).Warn("Hello") },
			pat: "{\"@level\":\"warn\",\"@time\":\"" + timePattern + "\",\"@scope\":\"default\",\"@caller\":\"log/.*_test.go:.*\",\"@msg\":\"Hello\"," + "\"foo\":\"bar\",\"baz\":123,\"qux\":0.123," +
				"\"@stack\":\".*\"}",
			json:       true,
			caller:     true,
			stackLevel: DebugLevel,
		},
		{
			f: func() { WithLabels("foo", "bar").WithLabels("baz", 123, "qux", 0.123).Error("Hello") },
			pat: "{\"@level\":\"error\",\"@time\":\"" + timePattern + "\",\"@scope\":\"default\",\"@caller\":\"log/.*_test.go:.*\",\"@msg\":\"Hello\"," + "\"foo\":\"bar\",\"baz\":123,\"qux\":0.123," +
				"\"@stack\":\".*\"}",
			json:       true,
			caller:     true,
			stackLevel: DebugLevel,
		},
		{
			f: func() { WithLabels("foo", "bar").WithLabels("baz", 123, "qux", 0.123).Fatal("Hello") },
			pat: "{\"@level\":\"fatal\",\"@time\":\"" + timePattern + "\",\"@scope\":\"default\",\"@caller\":\"log/.*_test.go:.*\",\"@msg\":\"Hello\"," + "\"foo\":\"bar\",\"baz\":123,\"qux\":0.123," +
				"\"@stack\":\".*\"}",
			json:       true,
			caller:     true,
			wantExit:   true,
			stackLevel: DebugLevel,
		},

		{
			f: func() { Debug("Hello") },
			pat: "{\"@level\":\"debug\",\"@time\":\"" + timePattern + "\",\"@scope\":\"default\",\"@caller\":\"log/.*_test.go:.*\",\"@msg\":\"Hello\"," +
				"\"@stack\":\".*\"}",
			json:       true,
			caller:     true,
			stackLevel: DebugLevel,
		},
		{
			f: func() { Info("Hello") },
			pat: "{\"@level\":\"info\",\"@time\":\"" + timePattern + "\",\"@scope\":\"default\",\"@caller\":\"log/.*_test.go:.*\",\"@msg\":\"Hello\"," +
				"\"@stack\":\".*\"}",
			json:       true,
			caller:     true,
			stackLevel: DebugLevel,
		},
		{
			f: func() { Warn("Hello") },
			pat: "{\"@level\":\"warn\",\"@time\":\"" + timePattern + "\",\"@scope\":\"default\",\"@caller\":\"log/.*_test.go:.*\",\"@msg\":\"Hello\"," +
				"\"@stack\":\".*\"}",
			json:       true,
			caller:     true,
			stackLevel: DebugLevel,
		},
		{
			f: func() { Error("Hello") },
			pat: "{\"@level\":\"error\",\"@time\":\"" + timePattern + "\",\"@scope\":\"default\",\"@caller\":\"log/.*_test.go:.*\",\"@msg\":\"Hello\"," +
				"\"@stack\":\".*\"}",
			json:       true,
			caller:     true,
			stackLevel: DebugLevel,
		},
		{
			f: func() { Fatal("Hello") },
			pat: "{\"@level\":\"fatal\",\"@time\":\"" + timePattern + "\",\"@scope\":\"default\",\"@caller\":\"log/.*_test.go:.*\",\"@msg\":\"Hello\"," +
				"\"@stack\":\".*\"}",
			json:       true,
			caller:     true,
			wantExit:   true,
			stackLevel: DebugLevel,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var exitCalled bool
			lines, err := captureStdout(func() {
				o := testOptions()
				o.JSONEncoding = c.json

				if err := Configure(o); err != nil {
					t.Errorf("Got err '%v', expecting success", err)
				}

				// hack 掉 Fatal 场景下, 需要退出程序的细节, 保证程序不会退出
				pt := funcs.Load().(patchTable)
				pt.exitProcess = func(_ int) {
					exitCalled = true
				}
				funcs.Store(pt)

				// 设置 scope
				defaultScope.SetOutputLevel(DebugLevel)
				defaultScope.SetStackTraceLevel(c.stackLevel)
				defaultScope.SetLogCallers(c.caller)

				c.f()
				_ = Sync()
			})

			if exitCalled != c.wantExit {
				var verb string
				if c.wantExit {
					verb = " never"
				}
				t.Errorf("os.Exit%s called", verb)
			}

			if err != nil {
				t.Errorf("Got error '%v', expected success", err)
			}

			if match, _ := regexp.MatchString(c.pat, lines[0]); !match {
				t.Errorf("Got logEntry '%v', expected a match with '%v'", lines[0], c.pat)
			} else {
				t.Logf("Test success, got: %v", lines[0])
			}
		})
	}
}

// 测试其他scope的日志输出
func TestBasicScopesLogging(t *testing.T) {
	resetGlobals()
	s := RegisterScope("testScope", "z", 0)

	cases := []struct {
		f          func()
		pat        string
		json       bool
		caller     bool
		wantExit   bool
		stackLevel Level
	}{
		{
			// zap.Field is no longer supported, prints like regular Sprint.
			f:   func() { s.Debug("Hello", zap.String("key", "value"), zap.Int("intkey", 123)) },
			pat: timePattern + "\tdebug\ttestScope\tHello{key 15 0 value <nil>} {intkey 11 123  <nil>}",
		},
		{
			f:   func() { s.Debug("Hello", " some", " fields") },
			pat: timePattern + "\tdebug\ttestScope\tHello some fields",
		},
		{
			f:   func() { s.Debugf("Hello") },
			pat: timePattern + "\tdebug\ttestScope\tHello",
		},
		{
			f:   func() { s.Debugf("%s", "Hello") },
			pat: timePattern + "\tdebug\ttestScope\tHello",
		},

		{
			f:   func() { s.Info("Hello", zap.String("key", "value"), zap.Int("intkey", 123)) },
			pat: timePattern + "\tinfo\ttestScope\tHello{key 15 0 value <nil>} {intkey 11 123  <nil>}",
		},
		{
			f:   func() { s.Info("Hello", " some", " fields") },
			pat: timePattern + "\tinfo\ttestScope\tHello some fields",
		},
		{
			f:   func() { s.Infof("Hello") },
			pat: timePattern + "\tinfo\ttestScope\tHello",
		},
		{
			f:   func() { s.Infof("%s", "Hello") },
			pat: timePattern + "\tinfo\ttestScope\tHello",
		},

		{
			f:   func() { s.Warn("Hello", zap.String("key", "value"), zap.Int("intkey", 123)) },
			pat: timePattern + "\twarn\ttestScope\tHello{key 15 0 value <nil>} {intkey 11 123  <nil>}",
		},
		{
			f:   func() { s.Warn("Hello", " some", " fields") },
			pat: timePattern + "\twarn\ttestScope\tHello some fields",
		},
		{
			f:   func() { s.Warnf("Hello") },
			pat: timePattern + "\twarn\ttestScope\tHello",
		},
		{
			f:   func() { s.Warnf("%s", "Hello") },
			pat: timePattern + "\twarn\ttestScope\tHello",
		},

		{
			f:   func() { s.Error("Hello", zap.String("key", "value"), zap.Int("intkey", 123)) },
			pat: timePattern + "\terror\ttestScope\tHello{key 15 0 value <nil>} {intkey 11 123  <nil>}",
		},
		{
			f:   func() { s.Errorf("Hello") },
			pat: timePattern + "\terror\ttestScope\tHello",
		},
		{
			f:   func() { s.Errorf("%s", "Hello") },
			pat: timePattern + "\terror\ttestScope\tHello",
		},

		{
			f:        func() { s.Fatal("Hello") },
			pat:      timePattern + "\tfatal\ttestScope\tHello",
			wantExit: true,
		},
		{
			f:        func() { s.Fatal("Hello", zap.String("key", "value"), zap.Int("intkey", 123)) },
			pat:      timePattern + "\tfatal\ttestScope\tHello{key 15 0 value <nil>} {intkey 11 123  <nil>}",
			wantExit: true,
		},
		{
			f:        func() { s.Fatal("Hello", " some", " fields") },
			pat:      timePattern + "\tfatal\ttestScope\tHello some fields",
			wantExit: true,
		},
		{
			f:        func() { s.Fatalf("Hello") },
			pat:      timePattern + "\tfatal\ttestScope\tHello",
			wantExit: true,
		},
		{
			f:        func() { s.Fatalf("%s", "Hello") },
			pat:      timePattern + "\tfatal\ttestScope\tHello",
			wantExit: true,
		},

		{
			f:      func() { s.Debug("Hello") },
			pat:    timePattern + "\tdebug\ttestScope\tlog/.*_test.go:.*\tHello",
			caller: true,
		},

		{
			f: func() { s.Debug("Hello") },
			pat: "{\"@level\":\"debug\",\"@time\":\"" + timePattern + "\",\"@scope\":\"testScope\",\"@caller\":\"log/.*_test.go:.*\",\"@msg\":\"Hello\"," +
				"\"@stack\":\".*\"}",
			json:       true,
			caller:     true,
			stackLevel: DebugLevel,
		},
		{
			f: func() { s.Info("Hello") },
			pat: "{\"@level\":\"info\",\"@time\":\"" + timePattern + "\",\"@scope\":\"testScope\",\"@caller\":\"log/.*_test.go:.*\",\"@msg\":\"Hello\"," +
				"\"@stack\":\".*\"}",
			json:       true,
			caller:     true,
			stackLevel: DebugLevel,
		},
		{
			f: func() { s.Warn("Hello") },
			pat: "{\"@level\":\"warn\",\"@time\":\"" + timePattern + "\",\"@scope\":\"testScope\",\"@caller\":\"log/.*_test.go:.*\",\"@msg\":\"Hello\"," +
				"\"@stack\":\".*\"}",
			json:       true,
			caller:     true,
			stackLevel: DebugLevel,
		},
		{
			f: func() { s.Error("Hello") },
			pat: "{\"@level\":\"error\",\"@time\":\"" + timePattern + "\",\"@scope\":\"testScope\",\"@caller\":\"log/.*_test.go:.*\"," +
				"\"@msg\":\"Hello\"," +
				"\"@stack\":\".*\"}",
			json:       true,
			caller:     true,
			stackLevel: DebugLevel,
		},
		{
			f: func() { s.Fatal("Hello") },
			pat: "{\"@level\":\"fatal\",\"@time\":\"" + timePattern + "\",\"@scope\":\"testScope\",\"@caller\":\"log/.*_test.go:.*\"," +
				"\"@msg\":\"Hello\"," +
				"\"@stack\":\".*\"}",
			json:       true,
			caller:     true,
			wantExit:   true,
			stackLevel: DebugLevel,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var exitCalled bool
			lines, err := captureStdout(func() {
				o := testOptions()
				o.JSONEncoding = c.json

				if err := Configure(o); err != nil {
					t.Errorf("Got err '%v', expecting success", err)
				}

				pt := funcs.Load().(patchTable)
				pt.exitProcess = func(_ int) {
					exitCalled = true
				}
				funcs.Store(pt)

				s.SetOutputLevel(DebugLevel)
				s.SetStackTraceLevel(c.stackLevel)
				s.SetLogCallers(c.caller)

				c.f()
				_ = Sync()
			})

			if exitCalled != c.wantExit {
				var verb string
				if c.wantExit {
					verb = " never"
				}
				t.Errorf("os.Exit%s called", verb)
			}

			if err != nil {
				t.Errorf("Got error '%v', expected success", err)
			}

			if match, _ := regexp.MatchString(c.pat, lines[0]); !match {
				t.Errorf("Got '%v', expected a match with '%v'", lines[0], c.pat)
			} else {
				t.Logf("Test success, got: %v", lines[0])
			}
		})
	}
}

// 测试 带 labels 的 scope 日志输出
func TestDefaultWithLabel(t *testing.T) {
	lines, err := captureStdout(func() {
		o := testOptions()
		Configure(o)
		funcs.Store(funcs.Load().(patchTable))
		defaultScope.SetLogCallers(true)
		WithLabels("foo", "bar").WithLabels("baz", 123, "qux", 0.123).Error("Hello")

		_ = Sync()
	})
	if err != nil {
		t.Errorf("Got error '%v', expected success", err)
	}

	mustRegexMatchString(t, lines[0], `Hello	foo=bar baz=123 qux=0.123`)
}

// 测试 带labels 的 其他 scope 日志输出
func TestScopeWithLabel(t *testing.T) {
	resetGlobals()
	const name = "TestScope"
	const desc = "Desc"
	s := RegisterScope(name, desc, 0)
	s.SetOutputLevel(DebugLevel)

	lines, err := captureStdout(func() {
		Configure(DefaultOptions())
		funcs.Store(funcs.Load().(patchTable))
		s2 := s.WithLabels("foo", "bar").WithLabels("baz", 123, "qux", 0.123)
		s2.Debug("Hello")
		// s 应该是未被修改的
		s.Debug("Hello")

		_ = Sync()
	})
	if err != nil {
		t.Errorf("Got error '%v', expected success", err)
	}

	mustRegexMatchString(t, lines[0], `Hello	foo=bar baz=123 qux=0.123`)
	mustRegexMatchString(t, lines[1], "Hello$")
}

// 测试 命令行设置的 日志级别  和 Configure 函数中 给 scope 设置的 日志级别 是否匹配
func TestDefaultEnabled(t *testing.T) {
	cases := []struct {
		level        Level
		debugEnabled bool
		infoEnabled  bool
		warnEnabled  bool
		errorEnabled bool
		fatalEnabled bool
	}{
		{DebugLevel, true, true, true, true, true},
		{InfoLevel, false, true, true, true, true},
		{WarnLevel, false, false, true, true, true},
		{ErrorLevel, false, false, false, true, true},
		{FatalLevel, false, false, false, false, true},
		{NoneLevel, false, false, false, false, false},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			o := testOptions()
			o.SetOutputLevel(DefaultScopeName, c.level) // 命令行设置

			_ = Configure(o) // 这里会将 命令行的设置 给 update 到 对应的 scope 对象中

			pt := funcs.Load().(patchTable)
			pt.exitProcess = func(_ int) {
			}
			funcs.Store(pt)

			if c.debugEnabled != DebugEnabled() {
				t.Errorf("Got %v, expecting %v", DebugEnabled(), c.debugEnabled)
			}

			if c.infoEnabled != InfoEnabled() {
				t.Errorf("Got %v, expecting %v", InfoEnabled(), c.infoEnabled)
			}

			if c.warnEnabled != WarnEnabled() {
				t.Errorf("Got %v, expecting %v", WarnEnabled(), c.warnEnabled)
			}

			if c.errorEnabled != ErrorEnabled() {
				t.Errorf("Got %v, expecting %v", ErrorEnabled(), c.errorEnabled)
			}

			if c.fatalEnabled != FatalEnabled() {
				t.Errorf("Got %v, expecting %v", FatalEnabled(), c.fatalEnabled)
			}
		})
	}
}

// 测试 命令行设置的 日志级别  和 Configure 函数中 给 scope 设置的 日志级别 是否匹配
func TestScopeEnabled(t *testing.T) {
	resetGlobals()
	const name = "TestEnabled"
	const desc = "Desc"
	s := RegisterScope(name, desc, 0)

	if n := s.Name(); n != name {
		t.Errorf("Got %s, expected %s", n, name)
	}

	if d := s.Description(); d != desc {
		t.Errorf("Got %s, expected %s", d, desc)
	}

	cases := []struct {
		level        Level
		debugEnabled bool
		infoEnabled  bool
		warnEnabled  bool
		errorEnabled bool
		fatalEnabled bool
	}{
		{NoneLevel, false, false, false, false, false},
		{FatalLevel, false, false, false, false, true},
		{ErrorLevel, false, false, false, true, true},
		{WarnLevel, false, false, true, true, true},
		{InfoLevel, false, true, true, true, true},
		{DebugLevel, true, true, true, true, true},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			s.SetOutputLevel(c.level)

			if c.debugEnabled != s.DebugEnabled() {
				t.Errorf("Got %v, expected %v", s.DebugEnabled(), c.debugEnabled)
			}

			if c.infoEnabled != s.InfoEnabled() {
				t.Errorf("Got %v, expected %v", s.InfoEnabled(), c.infoEnabled)
			}

			if c.warnEnabled != s.WarnEnabled() {
				t.Errorf("Got %v, expected %v", s.WarnEnabled(), c.warnEnabled)
			}

			if c.errorEnabled != s.ErrorEnabled() {
				t.Errorf("Got %v, expected %v", s.ErrorEnabled(), c.errorEnabled)
			}

			if c.fatalEnabled != s.FatalEnabled() {
				t.Errorf("Got %v, expected %v", s.FatalEnabled(), c.fatalEnabled)
			}

			if c.level != s.GetOutputLevel() {
				t.Errorf("Got %v, expected %v", s.GetOutputLevel(), c.level)
			}
		})
	}
}

// 测试 带 labels 的json 输出
func TestScopeJSON(t *testing.T) {
	resetGlobals()
	const name = "TestScope"
	const desc = "Desc"
	s := RegisterScope(name, desc, 0)
	s.SetOutputLevel(DebugLevel)

	lines, err := captureStdout(func() {
		o := DefaultOptions()
		o.JSONEncoding = true
		Configure(o)
		funcs.Store(funcs.Load().(patchTable))
		s.WithLabels("foo", "bar", "baz", 123).Debug("Hello")

		_ = Sync()
	})
	if err != nil {
		t.Errorf("Got error '%v', expected success", err)
	}

	mustRegexMatchString(t, lines[0], `{.*"@msg":"Hello","foo":"bar","baz":123}`)
}

// 测试日志截获, 截获 golang 标准库, grpc, 标准 zap 日志
func TestCapture(t *testing.T) {
	resetGlobals()
	lines, _ := captureStdout(func() {
		o := DefaultOptions()
		o.SetLogCallers(DefaultScopeName, true)
		o.SetOutputLevel(DefaultScopeName, DebugLevel)
		o.LogGrpc = true
		_ = Configure(o)

		// 使用 golang 标准库 输出日志
		log.Println("golang")

		// 使用 grpc 日志库 输出日志
		grpclog.Error("grpc-error")
		grpclog.Warning("grpc-warn")
		grpclog.Info("grpc-info")

		// 使用 zap 输出 日志
		zap.L().Error("zap-error")
		zap.L().Warn("zap-warn")
		zap.L().Info("zap-info")
		zap.L().Debug("zap-debug")

		l := zap.L().With(zap.String("a", "b"))
		l.Error("zap-with")

		entry := zapcore.Entry{
			Message: "zap-write",
			Level:   zapcore.ErrorLevel,
		}
		_ = zap.L().Core().Write(entry, nil)

		defaultScope.SetOutputLevel(NoneLevel) // 设置为 none

		// 由于日志设置为了none, 所以下面这些日志都会被丢弃
		log.Println("golang-2")
		grpclog.Error("grpc-error-2")
		grpclog.Warning("grpc-warn-2")
		grpclog.Info("grpc-info-2")
		zap.L().Error("zap-error-2")
		zap.L().Warn("zap-warn-2")
		zap.L().Info("zap-info-2")
		zap.L().Debug("zap-debug-2")
	})

	patterns := []string{
		timePattern + "\tinfo\tlog/.*_test.go:.*\tgolang",
		timePattern + "\tinfo\tlog/.*_test.go:.*\tgrpc-error",
		timePattern + "\tinfo\tlog/.*_test.go:.*\tgrpc-warn",
		timePattern + "\tinfo\tlog/.*_test.go:.*\tgrpc-info",
		timePattern + "\terror\tlog/.*_test.go:.*\tzap-error",
		timePattern + "\twarn\tlog/.*_test.go:.*\tzap-warn",
		timePattern + "\tinfo\tlog/.*_test.go:.*\tzap-info",
		timePattern + "\tdebug\tlog/.*_test.go:.*\tzap-debug",
		timePattern + "\terror\tlog/.*_test.go:.*\tzap-with",
		timePattern + "\terror\tzap-write",
		"",
	}

	if len(lines) > len(patterns) {
		t.Errorf("Expecting %d lines of output, but got %d", len(patterns), len(lines))

		for i := len(patterns); i < len(lines); i++ {
			t.Errorf("  Extra line of output: %s", lines[i])
		}
	}

	for i, pat := range patterns {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			match, _ := regexp.MatchString(pat, lines[i])
			if !match {
				t.Errorf("Got '%s', expecting to match '%s'", lines[i], pat)
			} else {
				t.Logf("Test success, got: %v", lines[i])
			}
		})
	}

	lines, _ = captureStdout(func() {
		o := DefaultOptions()
		o.SetStackTraceLevel(DefaultScopeName, DebugLevel)
		o.SetOutputLevel(DefaultScopeName, DebugLevel)
		_ = Configure(o)
		log.Println("golang")
	})

	for _, line := range lines {
		if strings.Contains(line, "_test.go") {
			return
		}
	}

	t.Error("Could not find stack trace info in output")
}

// 测试 klog 的截获 是否正常
func TestKlog(t *testing.T) {
	resetGlobals()
	gf := klogVerboseFlag()
	gf.Value.Set("5")
	cases := []struct {
		log      func()
		expected string
	}{
		{
			func() { klog.Error("foo") },
			"error\tklog\tfoo$",
		},
		{
			func() { klog.Error("foo") },
			"foo$",
		},
		{
			func() { klog.Errorf("fmt %v", "item") },
			"fmt item$",
		},
		{
			func() { klog.Errorf("fmt %v", "item") },
			"fmt item$",
		},
		{
			func() { klog.ErrorS(errors.New("my error"), "info") },
			"error\tklog\tmy error: info",
		},
		{
			func() { klog.Info("a", "b") },
			"info\tklog\tab",
		},
		{
			func() { klog.InfoS("msg", "key", 1) },
			"info\tklog\tmsg\tkey=1",
		},
		{
			func() { klog.ErrorS(errors.New("my error"), "info", "key", 1) },
			"error\tklog\tmy error: info\tkey=1",
		},
		{
			func() { klog.V(4).Info("debug") },
			"debug\tklog\tdebug$",
		},
	}
	for _, tt := range cases {
		t.Run(tt.expected, func(t *testing.T) {
			lines := runTest(t, tt.log)
			mustMatchLength(t, 1, lines)
			mustRegexMatchString(t, lines[0], tt.expected)
		})
	}
}

func runTest(t *testing.T, f func()) []string {
	lines, err := captureStdout(func() {
		Configure(DefaultOptions())
		f()
		_ = Sync()
	})
	if err != nil {
		t.Fatalf("Got error '%v', expected success", err)
	}
	if lines[len(lines)-1] == "" {
		return lines[:len(lines)-1]
	}
	return lines
}

func mustRegexMatchString(t *testing.T, got, want string) {
	t.Helper()
	match, _ := regexp.MatchString(want, got)

	if !match {
		t.Fatalf("Got '%v', expected a match with '%v'", got, want)
	} else {
		t.Logf("Test success, got: %v", got)
	}
}

func mustMatchLength(t *testing.T, l int, items []string) {
	t.Helper()
	if len(items) != l {
		t.Fatalf("expected %v items, got %v: %v", l, len(items), items)
	}
}
