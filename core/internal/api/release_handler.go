package api

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	releasev1 "github.com/yallaops/yallaops/core/internal/gen/release/v1"
	"github.com/yallaops/yallaops/core/internal/release"
	"github.com/yallaops/yallaops/core/internal/store"

	"github.com/google/uuid"
)

type ReleaseHandler struct {
	releasev1.UnimplementedReleaseServiceServer
	svc *release.Service
}

func NewReleaseHandler(svc *release.Service) *ReleaseHandler {
	return &ReleaseHandler{svc: svc}
}

func (h *ReleaseHandler) CreateRelease(ctx context.Context, req *releasev1.CreateReleaseRequest) (*releasev1.CreateReleaseResponse, error) {
	if req.Service == "" || req.Version == "" {
		return nil, status.Error(codes.InvalidArgument, "service and version are required")
	}

	id := fmt.Sprintf("%s-%s-%s", req.Service, req.Version, uuid.New().String()[:8])

	r, err := h.svc.Create(ctx, id, req.Service, req.Version)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create release: %v", err)
	}

	return &releasev1.CreateReleaseResponse{Release: toProtoRelease(r)}, nil
}

func (h *ReleaseHandler) GetRelease(ctx context.Context, req *releasev1.GetReleaseRequest) (*releasev1.GetReleaseResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	r, err := h.svc.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "release not found: %v", err)
	}

	return &releasev1.GetReleaseResponse{Release: toProtoRelease(r)}, nil
}

func (h *ReleaseHandler) ListReleases(ctx context.Context, req *releasev1.ListReleasesRequest) (*releasev1.ListReleasesResponse, error) {
	if req.Service == "" {
		return nil, status.Error(codes.InvalidArgument, "service is required")
	}

	limit := req.PageSize
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	releases, err := h.svc.List(ctx, req.Service, limit, 0)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list releases: %v", err)
	}

	out := make([]*releasev1.Release, len(releases))
	for i, r := range releases {
		out[i] = toProtoRelease(r)
	}

	return &releasev1.ListReleasesResponse{Releases: out}, nil
}

func (h *ReleaseHandler) UpdateReleaseStatus(ctx context.Context, req *releasev1.UpdateReleaseStatusRequest) (*releasev1.UpdateReleaseStatusResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	to, err := protoStatusToStore(req.Status)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r, err := h.svc.Transition(ctx, req.Id, to)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "update status: %v", err)
	}

	return &releasev1.UpdateReleaseStatusResponse{Release: toProtoRelease(r)}, nil
}

func toProtoRelease(r store.Release) *releasev1.Release {
	return &releasev1.Release{
		Id:        r.ID,
		Service:   r.Service,
		Version:   r.Version,
		Status:    storeStatusToProto(r.Status),
		CreatedAt: timestamppb.New(r.CreatedAt),
		UpdatedAt: timestamppb.New(r.UpdatedAt),
	}
}

func storeStatusToProto(s store.ReleaseStatus) releasev1.ReleaseStatus {
	switch s {
	case store.ReleaseStatusDraft:
		return releasev1.ReleaseStatus_RELEASE_STATUS_DRAFT
	case store.ReleaseStatusRunning:
		return releasev1.ReleaseStatus_RELEASE_STATUS_RUNNING
	case store.ReleaseStatusDeployed:
		return releasev1.ReleaseStatus_RELEASE_STATUS_DEPLOYED
	case store.ReleaseStatusFailed:
		return releasev1.ReleaseStatus_RELEASE_STATUS_FAILED
	default:
		return releasev1.ReleaseStatus_RELEASE_STATUS_UNSPECIFIED
	}
}

func protoStatusToStore(s releasev1.ReleaseStatus) (store.ReleaseStatus, error) {
	switch s {
	case releasev1.ReleaseStatus_RELEASE_STATUS_DRAFT:
		return store.ReleaseStatusDraft, nil
	case releasev1.ReleaseStatus_RELEASE_STATUS_RUNNING:
		return store.ReleaseStatusRunning, nil
	case releasev1.ReleaseStatus_RELEASE_STATUS_DEPLOYED:
		return store.ReleaseStatusDeployed, nil
	case releasev1.ReleaseStatus_RELEASE_STATUS_FAILED:
		return store.ReleaseStatusFailed, nil
	default:
		return "", fmt.Errorf("unknown status: %v", s)
	}
}
