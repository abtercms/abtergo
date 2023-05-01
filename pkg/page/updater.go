package page

import (
	"fmt"

	"github.com/qmuntal/stateless"
)

// Trigger represent a trigger for a status change of a Page.
type Trigger string

const (
	// Activate represent a trigger for a status change of a Page which aims to move the status to Active.
	Activate Trigger = "activate"
	// Inactivate represent a trigger for a status change of a Page which aims to move the status to Inactive.
	Inactivate Trigger = "inactivate"
)

// Updater represent business logic for various processes.
type Updater interface {
	Transition(status Status, trigger Trigger) (Status, error)
}

type updater struct{}

// Transition attempts to apply a trigger starting the state machine in a given status.
func (u *updater) Transition(status Status, trigger Trigger) (Status, error) {
	sm := u.getStateMachine(status)

	err := sm.Fire(trigger)
	if err != nil {
		return status, fmt.Errorf("invalid state transition. err: %w", err)
	}

	return sm.MustState().(Status), nil
}

func (u *updater) getStateMachine(state Status) *stateless.StateMachine {
	sm := stateless.NewStateMachine(state)

	sm.Configure(Draft).
		Permit(Inactivate, Inactive).
		Permit(Activate, Active)

	sm.Configure(Active).
		Permit(Inactivate, Inactive)

	sm.Configure(Inactive).
		Permit(Activate, Active)

	return sm
}

// NewUpdater creates a new Updater instance.
func NewUpdater() Updater {
	return &updater{}
}
