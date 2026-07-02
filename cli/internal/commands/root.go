// Package commands implements the yallaops CLI's Cobra command tree.
package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/yallaops/yallaops/cli/internal/client"
	"github.com/yallaops/yallaops/cli/internal/config"
)

var endpointOverride string

// Execute runs the root command, returning any error for main to report.
func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "yallaops",
		Short: "Manage YallaOps releases from the command line",
		Long:  "yallaops talks to the YallaOps control plane over gRPC to create releases, promote them across environments, and inspect their status.",
	}

	root.PersistentFlags().StringVar(&endpointOverride, "endpoint", "", "control plane gRPC endpoint (overrides the current context)")

	root.AddCommand(newCreateCmd())
	root.AddCommand(newPromoteCmd())
	root.AddCommand(newStatusCmd())
	root.AddCommand(newDashboardCmd())

	return root
}

// dial loads the active config context and opens a gRPC client, honoring
// --endpoint as an override.
func dial(ctx context.Context) (*client.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	endpoint := endpointOverride
	if endpoint == "" {
		current, err := cfg.Current()
		if err != nil {
			return nil, err
		}
		endpoint = current.Endpoint
	}

	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	c, err := client.Dial(dialCtx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("commands: connect to %s: %w", endpoint, err)
	}
	return c, nil
}
