package cmd

import (
	"agdev/internal/instruction"
	"agdev/internal/output"

	"github.com/spf13/cobra"
)

type codeInstructionCommonOptions struct {
	version string
}

var codeInstructionCommonOpts = codeInstructionCommonOptions{}

func init() {
	codeInstructionCmd.AddCommand(codeInstructionCommonCmd)
	codeInstructionCommonCmd.Flags().StringVar(&codeInstructionCommonOpts.version, "version", "latest", "Instruction version to read")
}

var codeInstructionCommonCmd = &cobra.Command{
	Use:   "common",
	Short: "Read the common agent instruction text",
	Args:  exactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		doc, err := instruction.Get("code", "common", codeInstructionCommonOpts.version)
		if err != nil {
			return err
		}

		payload := map[string]string{
			"scope":            doc.Scope,
			"name":             doc.Name,
			"version":          doc.Version,
			"resolved_version": doc.ResolvedVersion,
			"instruction":      doc.Body,
		}

		return output.WriteSuccess(cmd.OutOrStdout(), currentConfig().OutputJSON, doc.Body, payload)
	},
}
