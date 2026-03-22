package cmd

import (
	"agdev/internal/auth"

	"github.com/spf13/cobra"
)

type authOptions struct {
	token string
}

var globalAuthOpts authOptions

func init() {
	rootCmd.PersistentFlags().StringVar(&globalAuthOpts.token, "token", "", "API token used for authenticated commands")
}

func withAuth(runE func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		resolver, err := auth.DefaultResolver()
		if err != nil {
			return err
		}

		token, err := resolver.Resolve(cmd.Context(), globalAuthOpts.token)
		if err != nil {
			return err
		}

		cmd.SetContext(auth.WithToken(cmd.Context(), token))

		return runE(cmd, args)
	}
}
