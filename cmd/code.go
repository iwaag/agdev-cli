package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(codeCmd)
}

var codeCmd = &cobra.Command{
	Use:   "code",
	Short: "Code-related operations",
}
