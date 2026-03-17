package cmd

import (
	"agdev/internal/app"

	"github.com/spf13/cobra"
)

func exactArgs(count int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(count)(cmd, args); err != nil {
			return app.WithExitCode(app.ExitUsage, err)
		}

		return nil
	}
}
