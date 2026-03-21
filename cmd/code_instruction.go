package cmd

import "github.com/spf13/cobra"

func init() {
	codeCmd.AddCommand(codeInstructionCmd)
}

var codeInstructionCmd = &cobra.Command{
	Use:   "instruction",
	Short: "Read agent instruction text",
}
