package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(videoCmd)
}

var videoCmd = &cobra.Command{
	Use:   "video",
	Short: "Video-related operations",
}
