package cmd

import (
	"context"
	"flag"
	"os"

	"github.com/bloodsteel/easynetes/pkg/config"
	"github.com/bloodsteel/easynetes/pkg/log"
	"github.com/bloodsteel/easynetes/pkg/signal"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	author         = []byte{10, 110, 98, 119, 104, 103, 64, 98, 108, 111, 111, 100, 115, 116, 101, 101, 108, 10, 10}
	loggingOptions = log.DefaultOptions()
	cfg            = &config.Config{}
	configFile     string
)

var (
	// Branch 当前运行的branch
	Branch string
	// Commit 当前运行的commit
	Commit string
	// BuildTime 构建时间
	BuildTime string
)

// EasynetesAPI 返回 easynetes 命令
func EasynetesAPI() *cobra.Command {
	easynetesAPICmd := &cobra.Command{
		Use:               "easynetes-api",
		Short:             "easynetes http apiserver.",
		Long:              "easynetes http apiserver.",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Version:           "dev",
		Args:              cobra.ExactArgs(0),
		Example:           "easynetes-api --config config.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			// start program...  printf some flags
			log.WithLabels(
				"git_branch", Branch,
				"git_commit", Commit,
				"built_time", BuildTime,
				"version", "dev",
			).Info("easynetes-api is starting...")

			// set root ctx & listen os signal
			ctx := signal.WithContextFunc(context.Background(), func() {
				log.Info("sync log...")
				_ = log.Sync()
			})

			// initialize global configure
			cfg.Load(configFile)

			// initialize application
			app, err := InitializeApplication(cfg)
			if err != nil {
				log.WithLabels("error", err).Fatal("cannot initialize application")
			}

			// start application
			g, ctx := errgroup.WithContext(ctx)
			g.Go(func() error {
				log.WithLabels(
					"bind_address", cfg.Server.BindAddress,
					"secure_port", cfg.Server.SecurePort,
					"insecure_port", cfg.Server.InsecurePort,
					"read_timeout", cfg.Server.ReadTimeout,
					"write_timeout", cfg.Server.WriteTimeout,
				).Info("starting the http server")
				return app.ListenAndServe(ctx)
			})
			if err := g.Wait(); err != nil {
				log.WithLabels("error", err).Error("program terminated")
			}
			return nil
		},
	}

	// 禁用 completion 子命令
	easynetesAPICmd.CompletionOptions.DisableDefaultCmd = true

	// attach 日志的flag 到 指定的 cmd
	loggingOptions.AttachCobraFlags(easynetesAPICmd)
	log.EnableKlogWithCobra()

	easynetesAPICmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Set configuration file path")
	easynetesAPICmd.MarkPersistentFlagRequired("config")

	// 添加go原生flag
	easynetesAPICmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	// 设置命令输出模板
	versionTemplate := easynetesAPICmd.VersionTemplate()
	helpTemplate := easynetesAPICmd.HelpTemplate()
	easynetesAPICmd.SetHelpTemplate(helpTemplate + string(author))
	easynetesAPICmd.SetVersionTemplate(versionTemplate + string(author))

	easynetesAPICmd.SetArgs(os.Args[1:])
	return easynetesAPICmd
}
