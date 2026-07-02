// Package client wires up gRPC connections to the YallaOps control plane.
package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	environmentv1 "github.com/yallaops/yallaops/cli/internal/gen/environment/v1"
	releasev1 "github.com/yallaops/yallaops/cli/internal/gen/release/v1"
)

// Client bundles the gRPC service clients the CLI talks to.
type Client struct {
	conn        *grpc.ClientConn
	Release     releasev1.ReleaseServiceClient
	Environment environmentv1.EnvironmentServiceClient
}

// Dial opens a gRPC connection to the control plane at endpoint.
//
// TODO(phase-3): switch to TLS transport credentials once the control plane
// exposes a secure listener; insecure is fine for local dev against
// localhost:50051.
func Dial(ctx context.Context, endpoint string) (*Client, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("client: dial %s: %w", endpoint, err)
	}

	return &Client{
		conn:        conn,
		Release:     releasev1.NewReleaseServiceClient(conn),
		Environment: environmentv1.NewEnvironmentServiceClient(conn),
	}, nil
}

// Close releases the underlying gRPC connection.
func (c *Client) Close() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("client: close: %w", err)
	}
	return nil
}
