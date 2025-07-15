package cmd

import (
	"github.com/risy007/wx-gw/internal"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var configFile string

func init() {
	pf := ServeCmd.PersistentFlags()
	pf.StringVarP(&configFile, "config", "c",
		"./config.yaml", "这是服务配置文件默认位置，也可以通过命令行参数指定")

	rootCmd.AddCommand(ServeCmd)
}

var ServeCmd = &cobra.Command{
	Use:     "serve",
	Short:   "启动服务",
	Example: "{execfile} serve -c ./config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		runApplication()
	},
}

func runApplication() {
	fx.New(
		internal.Module,
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.NopLogger,
	).Run()
}
