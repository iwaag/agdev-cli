package cmd

import (
	"context"

	"agdev/internal/config"

	"github.com/spf13/cobra"
)

var (
	cfg      config.Config
	rootOpts rootOptions
)

type rootOptions struct {
	json bool
}

var rootCmd = &cobra.Command{
	Use:           "agdev",
	Short:         "Agent-oriented CLI bridge for AGDEV backends",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		loaded, err := config.Load(rootOpts.json)
		if err != nil {
			return err
		}

		cfg = loaded
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&rootOpts.json, "json", false, "Emit machine-readable JSON output")
}

func ExecuteContext(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}

func currentConfig() config.Config {
	return cfg
}
