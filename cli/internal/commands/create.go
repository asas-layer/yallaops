package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	releasev1 "github.com/yallaops/yallaops/cli/internal/gen/release/v1"
)

func newCreateCmd() *cobra.Command {
	var service, version string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new release",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dial(cmd.Context())
			if err != nil {
				return err
			}
			defer c.Close()

			resp, err := c.Release.CreateRelease(cmd.Context(), &releasev1.CreateReleaseRequest{
				Service: service,
				Version: version,
			})
			if err != nil {
				return fmt.Errorf("create: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "created release %s (%s@%s) status=%s\n",
				resp.GetRelease().GetId(), service, version, resp.GetRelease().GetStatus())
			return nil
		},
	}

	cmd.Flags().StringVar(&service, "service", "", "service name (required)")
	cmd.Flags().StringVar(&version, "version", "", "release version (required)")
	cmd.MarkFlagRequired("service")
	cmd.MarkFlagRequired("version")

	return cmd
}
