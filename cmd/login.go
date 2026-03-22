package cmd

import (
	"fmt"
	"strings"

	"agdev/internal/app"
	"agdev/internal/auth"

	"github.com/spf13/cobra"
)

type loginOptions struct {
	user string
}

var loginOpts loginOptions

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVar(&loginOpts.user, "user", "", "Keycloak username")
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in with Keycloak and save a local session",
	Args:  exactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := auth.NewKeycloakClientFromEnv()
		if err != nil {
			return err
		}

		username, err := resolveLoginUser(loginOpts.user, client.DefaultUser())
		if err != nil {
			return err
		}

		password, err := auth.PromptPassword("Password: ")
		if err != nil {
			return app.WithExitCode(app.ExitUsage, fmt.Errorf("read password: %w", err))
		}
		if strings.TrimSpace(password) == "" {
			return app.WithExitCode(app.ExitUsage, fmt.Errorf("password is required"))
		}

		session, err := client.LoginPassword(cmd.Context(), username, password)
		if err != nil {
			return err
		}

		store, err := auth.NewFileStore()
		if err != nil {
			return app.WithExitCode(app.ExitInternal, err)
		}
		if err := store.WriteSession(cmd.Context(), session); err != nil {
			return app.WithExitCode(app.ExitInternal, fmt.Errorf("save login session: %w", err))
		}

		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Logged in as %s\n", username); err != nil {
			return app.WithExitCode(app.ExitInternal, err)
		}

		return nil
	},
}

func resolveLoginUser(flagUser, defaultUser string) (string, error) {
	if user := strings.TrimSpace(flagUser); user != "" {
		return user, nil
	}
	if user := strings.TrimSpace(defaultUser); user != "" {
		return user, nil
	}

	user, err := auth.PromptLine("Username: ")
	if err != nil {
		return "", app.WithExitCode(app.ExitUsage, fmt.Errorf("read username: %w", err))
	}
	if user == "" {
		return "", app.WithExitCode(app.ExitUsage, fmt.Errorf("username is required"))
	}

	return user, nil
}
