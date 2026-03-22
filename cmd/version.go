package cmd

import (
	"agdev/internal/output"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print build version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := map[string]string{
			"version": version,
			"commit":  commit,
			"date":    date,
		}

		return output.WriteJSON(cmd.OutOrStdout(), payload)
	},
}
