package cmd

import (
	"agdev/internal/api"
	"agdev/internal/output"

	"github.com/spf13/cobra"
)

func init() {
	imageCmd.AddCommand(imageGenerateCmd)
}

var imageGenerateCmd = &cobra.Command{
	Use:   "generate <input-image> <prompt>",
	Short: "Queue an image generation request",
	Args:  exactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		response := api.AcceptedResponse{
			Status:  "accepted",
			Kind:    "image.generate",
			Message: "skeleton command: backend call not implemented yet",
			Inputs: map[string]string{
				"input_image": args[0],
				"prompt":      args[1],
			},
		}

		text := "accepted image.generate request"
		return output.WriteSuccess(cmd.OutOrStdout(), currentConfig().OutputJSON, text, response)
	},
}
