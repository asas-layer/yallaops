package commands

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	environmentv1 "github.com/yallaops/yallaops/cli/internal/gen/environment/v1"
	releasev1 "github.com/yallaops/yallaops/cli/internal/gen/release/v1"
	"github.com/yallaops/yallaops/cli/internal/tui"
)

func newDashboardCmd() *cobra.Command {
	var theme string

	cmd := &cobra.Command{
		Use:     "dashboard",
		Aliases: []string{"tui"},
		Short:   "Interactive dashboard of releases and environment status",
		RunE: func(cmd *cobra.Command, args []string) error {
			model := tui.New(theme, dashboardFetch(cmd.Context()))
			_, err := tea.NewProgram(model).Run()
			if err != nil {
				return fmt.Errorf("dashboard: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&theme, "theme", "dark", fmt.Sprintf("dashboard color theme (%v)", tui.ThemeNames))

	return cmd
}

// dashboardFetch loads rows from the control plane. If the control plane is
// unreachable (e.g. not implemented yet), it returns placeholder rows so the
// dashboard still renders — useful for previewing the layout during Phase 2.
func dashboardFetch(ctx context.Context) tui.FetchFunc {
	return func() ([]tui.Row, error) {
		c, err := dial(ctx)
		if err != nil {
			return placeholderRows(), nil
		}
		defer c.Close()

		resp, err := c.Release.ListReleases(ctx, &releasev1.ListReleasesRequest{})
		if err != nil {
			return placeholderRows(), nil
		}

		rows := make([]tui.Row, 0, len(resp.GetReleases()))
		for _, rel := range resp.GetReleases() {
			row := tui.Row{
				ID:      rel.GetId(),
				Service: rel.GetService(),
				Version: rel.GetVersion(),
				Status:  rel.GetStatus().String(),
			}

			envResp, err := c.Environment.ListEnvironments(ctx, &environmentv1.ListEnvironmentsRequest{ReleaseId: rel.GetId()})
			if err == nil {
				for _, env := range envResp.GetEnvironments() {
					row.EnvNames = append(row.EnvNames, env.GetEnvName())
					row.EnvGood = append(row.EnvGood, env.GetStatus().String())
				}
			}
			rows = append(rows, row)
		}
		return rows, nil
	}
}

func placeholderRows() []tui.Row {
	return []tui.Row{
		{
			ID: "release-payment-api-1.4.2", Service: "payment-api", Version: "1.4.2",
			Status:   "RELEASE_STATUS_RUNNING",
			EnvNames: []string{"dev", "staging", "prod"},
			EnvGood:  []string{"ENVIRONMENT_STATUS_DEPLOYED", "ENVIRONMENT_STATUS_PENDING", "ENVIRONMENT_STATUS_BLOCKED"},
		},
		{
			ID: "release-checkout-web-2.0.0", Service: "checkout-web", Version: "2.0.0",
			Status:   "RELEASE_STATUS_DEPLOYED",
			EnvNames: []string{"dev", "staging", "prod"},
			EnvGood:  []string{"ENVIRONMENT_STATUS_DEPLOYED", "ENVIRONMENT_STATUS_DEPLOYED", "ENVIRONMENT_STATUS_DEPLOYED"},
		},
	}
}
