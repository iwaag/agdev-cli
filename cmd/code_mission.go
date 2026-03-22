package cmd

import (
	"agdev/internal/agcode"
	"agdev/internal/output"

	"github.com/spf13/cobra"
)

func init() {
	codeCmd.AddCommand(codeMissionCmd)
}

var codeMissionCmd = &cobra.Command{
	Use:   "mission <mission_id>",
	Short: "Get mission information",
	Args:  exactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := agcode.NewClient()

		mission, err := client.GetMission(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		return output.WriteJSON(cmd.OutOrStdout(), mission)
	},
}
