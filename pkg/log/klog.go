package log

import (
	goflag "flag"
	"sync"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

var (
	// KlogScope 注册一个 klog scope
	KlogScope     = RegisterScope("klog", "kubernetes logger", 0)
	configureKlog = sync.Once{}
)

// EnableKlogWithCobra 启用 klog 和 cobra / pflags 一起工作
// 像client-go这样的k8s库使用的是 klog 库
func EnableKlogWithCobra() {
	gf := klogVerboseFlag()
	pflag.CommandLine.AddFlag(pflag.PFlagFromGoFlag(
		&goflag.Flag{
			Name:     "vklog",
			Value:    gf.Value,
			DefValue: gf.DefValue,
			Usage:    gf.Usage + ". Like -v flag. ex: --vklog=9",
		}))
}

// isKlogVerbose 如果 klog 的 verbosity 是 非0值 那么返回 true
func klogVerbose() bool {
	gf := klogVerboseFlag()
	return gf.Value.String() != "0"
}

// KlogVerboseFlag 初始化klog 的 flags, 并从命令行获取 -v 设置的值
func klogVerboseFlag() *goflag.Flag {
	// 设置klog
	fs := &goflag.FlagSet{}
	klog.InitFlags(fs)
	// --v= flag of klog.
	return fs.Lookup("v")
}
