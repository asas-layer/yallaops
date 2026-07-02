package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	environmentv1 "github.com/yallaops/yallaops/cli/internal/gen/environment/v1"
	releasev1 "github.com/yallaops/yallaops/cli/internal/gen/release/v1"
)

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [release-id]",
		Short: "Show release status, or list all releases if no ID is given",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := dial(cmd.Context())
			if err != nil {
				return err
			}
			defer c.Close()

			if len(args) == 1 {
				return showRelease(cmd, c.Release, c.Environment, args[0])
			}
			return listReleases(cmd, c.Release)
		},
	}

	return cmd
}

func showRelease(cmd *cobra.Command, releases releasev1.ReleaseServiceClient, environments environmentv1.EnvironmentServiceClient, id string) error {
	relResp, err := releases.GetRelease(cmd.Context(), &releasev1.GetReleaseRequest{Id: id})
	if err != nil {
		return fmt.Errorf("status: get release: %w", err)
	}

	envResp, err := environments.ListEnvironments(cmd.Context(), &environmentv1.ListEnvironmentsRequest{ReleaseId: id})
	if err != nil {
		return fmt.Errorf("status: list environments: %w", err)
	}

	rel := relResp.GetRelease()
	fmt.Fprintf(cmd.OutOrStdout(), "%s  %s@%s  %s\n", rel.GetId(), rel.GetService(), rel.GetVersion(), rel.GetStatus())
	for _, env := range envResp.GetEnvironments() {
		fmt.Fprintf(cmd.OutOrStdout(), "  %-10s %s\n", env.GetEnvName(), env.GetStatus())
	}
	return nil
}

func listReleases(cmd *cobra.Command, releases releasev1.ReleaseServiceClient) error {
	resp, err := releases.ListReleases(cmd.Context(), &releasev1.ListReleasesRequest{})
	if err != nil {
		return fmt.Errorf("status: list releases: %w", err)
	}

	for _, rel := range resp.GetReleases() {
		fmt.Fprintf(cmd.OutOrStdout(), "%s  %s@%s  %s\n", rel.GetId(), rel.GetService(), rel.GetVersion(), rel.GetStatus())
	}
	return nil
}
