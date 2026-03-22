package cmd

import (
	"agdev/internal/agcode"
	"agdev/internal/auth"
	"agdev/internal/output"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	codeCmd.AddCommand(codeMissionCmd)
}

var codeMissionCmd = &cobra.Command{
	Use:   "mission <mission_id>",
	Short: "Get mission information",
	Args:  exactArgs(1),
	RunE: withAuth(func(cmd *cobra.Command, args []string) error {
		token, err := auth.TokenFromContext(cmd.Context())
		if err != nil {
			return err
		}

		client := agcode.NewClient(agcode.Config{
			BaseURL:   os.Getenv("AGCODE_API_URL"),
			AuthToken: token,
		})

		mission, err := client.GetMission(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		return output.WriteJSON(cmd.OutOrStdout(), mission)
	}),
}
