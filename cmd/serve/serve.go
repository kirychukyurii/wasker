package serve

import (
	"context"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func init() {
	_ = Command.PersistentFlags()
}

var Command = &cobra.Command{
	Use:          "serve",
	Short:        "",
	Example:      "",
	SilenceUsage: true,
	PreRun:       func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		fx.New(Module, fx.NopLogger).Run()
	},
}

var Module = fx.Options(
	fx.Invoke(runApplication),
)

func runApplication(lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {

			return nil
		},
		OnStop: func(context.Context) error {

			return nil
		},
	})
}
