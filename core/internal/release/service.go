package release

import (
	"context"
	"fmt"

	"github.com/yallaops/yallaops/core/internal/store"
)

type Service struct {
	q  store.Querier
	sm *StateMachine
}

func NewService(q store.Querier) *Service {
	return &Service{q: q, sm: &StateMachine{}}
}

func (s *Service) Create(ctx context.Context, id, service, version string) (store.Release, error) {
	r, err := s.q.CreateRelease(ctx, store.CreateReleaseParams{
		ID:      id,
		Service: service,
		Version: version,
	})
	if err != nil {
		return store.Release{}, fmt.Errorf("release service: create: %w", err)
	}
	return r, nil
}

func (s *Service) Get(ctx context.Context, id string) (store.Release, error) {
	r, err := s.q.GetRelease(ctx, id)
	if err != nil {
		return store.Release{}, fmt.Errorf("release service: get: %w", err)
	}
	return r, nil
}

func (s *Service) List(ctx context.Context, svc string, limit, offset int32) ([]store.Release, error) {
	releases, err := s.q.ListReleases(ctx, store.ListReleasesParams{
		Service: svc,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("release service: list: %w", err)
	}
	return releases, nil
}

func (s *Service) Transition(ctx context.Context, id string, to store.ReleaseStatus) (store.Release, error) {
	current, err := s.q.GetRelease(ctx, id)
	if err != nil {
		return store.Release{}, fmt.Errorf("release service: transition: get: %w", err)
	}

	if err := s.sm.ValidateTransition(current.Status, to); err != nil {
		return store.Release{}, fmt.Errorf("release service: transition: %w", err)
	}

	updated, err := s.q.UpdateReleaseStatus(ctx, store.UpdateReleaseStatusParams{
		ID:     id,
		Status: to,
	})
	if err != nil {
		return store.Release{}, fmt.Errorf("release service: transition: update: %w", err)
	}
	return updated, nil
}
