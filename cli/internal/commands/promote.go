package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	environmentv1 "github.com/yallaops/yallaops/cli/internal/gen/environment/v1"
)

func newPromoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "promote <release-id> <target-env>",
		Short: "Promote a release to the next environment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			releaseID, targetEnv := args[0], args[1]

			c, err := dial(cmd.Context())
			if err != nil {
				return err
			}
			defer c.Close()

			resp, err := c.Environment.PromoteRelease(cmd.Context(), &environmentv1.PromoteReleaseRequest{
				ReleaseId: releaseID,
				TargetEnv: targetEnv,
			})
			if err != nil {
				return fmt.Errorf("promote: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "promoted release %s to %s: status=%s\n",
				releaseID, targetEnv, resp.GetEnvironment().GetStatus())
			return nil
		},
	}

	return cmd
}
