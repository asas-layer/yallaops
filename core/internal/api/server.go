package api

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	releasev1 "github.com/yallaops/yallaops/core/internal/gen/release/v1"
)

type Server struct {
	grpc *grpc.Server
}

func NewServer(releaseHandler *ReleaseHandler) *Server {
	g := grpc.NewServer()
	releasev1.RegisterReleaseServiceServer(g, releaseHandler)
	reflection.Register(g)
	return &Server{grpc: g}
}

func (s *Server) ListenAndServe(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("api server: listen: %w", err)
	}
	return s.grpc.Serve(lis)
}

func (s *Server) GracefulStop() {
	s.grpc.GracefulStop()
}
