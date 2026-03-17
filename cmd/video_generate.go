package cmd

import (
	"agdev/internal/api"
	"agdev/internal/output"

	"github.com/spf13/cobra"
)

func init() {
	videoCmd.AddCommand(videoGenerateCmd)
}

var videoGenerateCmd = &cobra.Command{
	Use:   "generate <first-frame> <last-frame>",
	Short: "Queue a video generation request",
	Args:  exactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		response := api.AcceptedResponse{
			Status:  "accepted",
			Kind:    "video.generate",
			Message: "skeleton command: backend call not implemented yet",
			Inputs: map[string]string{
				"first_frame": args[0],
				"last_frame":  args[1],
			},
		}

		text := "accepted video.generate request"
		return output.WriteSuccess(cmd.OutOrStdout(), currentConfig().OutputJSON, text, response)
	},
}
