package cmd

import (
	"fmt"

	"agdev/internal/app"
	"agdev/internal/openapi"

	"github.com/spf13/cobra"
)

type utilOpenAPIOptions struct {
	out  string
	tags []string
}

var utilOpenAPIOpts utilOpenAPIOptions

func init() {
	utilCmd.AddCommand(utilOpenAPICmd)
	utilOpenAPICmd.Flags().StringVarP(&utilOpenAPIOpts.out, "out", "o", "", "Output file path")
	utilOpenAPICmd.Flags().StringSliceVar(&utilOpenAPIOpts.tags, "tags", nil, "Keep only operations with the specified tags")
}

var utilOpenAPICmd = &cobra.Command{
	Use:   "openapi <base_url>",
	Short: "Fetch and save an OpenAPI document",
	Args:  exactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		documentURL, err := openapi.ResolveDocumentURL(args[0])
		if err != nil {
			return err
		}

		payload, err := openapi.FetchDocument(cmd.Context(), documentURL)
		if err != nil {
			return err
		}

		if len(utilOpenAPIOpts.tags) > 0 {
			openapi.FilterOperationsByTags(payload, utilOpenAPIOpts.tags)
		}

		outputPath, err := openapi.ResolveOutputPath(utilOpenAPIOpts.out, payload)
		if err != nil {
			return err
		}

		if err := openapi.WriteDocument(outputPath, payload); err != nil {
			return err
		}

		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Saved %s\n", outputPath); err != nil {
			return app.WithExitCode(app.ExitInternal, err)
		}

		return nil
	},
}
