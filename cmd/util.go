package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(utilCmd)
}

var utilCmd = &cobra.Command{
	Use:   "util",
	Short: "Utility operations",
}
