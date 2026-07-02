package release

import (
	"fmt"

	"github.com/yallaops/yallaops/core/internal/store"
)

var validTransitions = map[store.ReleaseStatus][]store.ReleaseStatus{
	store.ReleaseStatusDraft:   {store.ReleaseStatusRunning},
	store.ReleaseStatusRunning: {store.ReleaseStatusDeployed, store.ReleaseStatusFailed},
}

type StateMachine struct{}

func (sm *StateMachine) CanTransition(from, to store.ReleaseStatus) bool {
	targets, ok := validTransitions[from]
	if !ok {
		return false
	}
	for _, t := range targets {
		if t == to {
			return true
		}
	}
	return false
}

func (sm *StateMachine) ValidateTransition(from, to store.ReleaseStatus) error {
	if !sm.CanTransition(from, to) {
		return fmt.Errorf("invalid transition from %s to %s", from, to)
	}
	return nil
}
