package release

import (
	"testing"

	"github.com/yallaops/yallaops/core/internal/store"
)

func TestStateMachine_CanTransition(t *testing.T) {
	sm := &StateMachine{}

	tests := []struct {
		from store.ReleaseStatus
		to   store.ReleaseStatus
		want bool
	}{
		{store.ReleaseStatusDraft, store.ReleaseStatusRunning, true},
		{store.ReleaseStatusRunning, store.ReleaseStatusDeployed, true},
		{store.ReleaseStatusRunning, store.ReleaseStatusFailed, true},
		{store.ReleaseStatusDraft, store.ReleaseStatusDeployed, false},
		{store.ReleaseStatusDeployed, store.ReleaseStatusRunning, false},
		{store.ReleaseStatusFailed, store.ReleaseStatusRunning, false},
	}

	for _, tt := range tests {
		got := sm.CanTransition(tt.from, tt.to)
		if got != tt.want {
			t.Errorf("CanTransition(%s, %s) = %v, want %v", tt.from, tt.to, got, tt.want)
		}
	}
}
