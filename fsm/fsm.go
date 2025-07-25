package fsm

import (
	"errors"
	"fmt"
)

// State represents a state in the FSM.
type State string

// Event represents an event that triggers a state transition.
type Event string

// Action is a function that can be executed when entering or exiting a state.
type Action func(args ...interface{}) error

// Guard is a function that determines if a transition is allowed.
type Guard func(args ...interface{}) bool

// Transition defines a state transition: current state, event, and next state.
type Transition struct {
	From  State
	Event Event
	To    State
}

// FSM represents a Finite State Machine.
type FSM struct {
	machineID           string
	currentState        State
	transitions         map[State]map[Event]State
	entryActions        map[State]Action
	exitActions         map[State]Action
	guards              map[State]map[Event]Guard
	transitionCallbacks map[State]map[Event]Action
}

// NewFSM creates a new FSM with an initial state and a list of transitions.
func NewFSM(machineID string, initialState State, transitions []Transition) (*FSM, error) {
	fsm := &FSM{
		machineID:           machineID,
		currentState:        initialState,
		transitions:         make(map[State]map[Event]State),
		entryActions:        make(map[State]Action),
		exitActions:         make(map[State]Action),
		guards:              make(map[State]map[Event]Guard),
		transitionCallbacks: make(map[State]map[Event]Action),
	}

	for _, t := range transitions {
		if _, ok := fsm.transitions[t.From]; !ok {
			fsm.transitions[t.From] = make(map[Event]State)
		}
		fsm.transitions[t.From][t.Event] = t.To
	}

	return fsm, nil
}

// CurrentState returns the current state of the FSM.
func (f *FSM) CurrentState() State {
	return f.currentState
}

// Transition attempts to transition the FSM to a new state based on an event.
func (f *FSM) Transition(event Event, args ...interface{}) error {
	nextStates, ok := f.transitions[f.currentState]
	if !ok {
		return fmt.Errorf("no transitions defined from state %s", f.currentState)
	}

	nextState, ok := nextStates[event]
	if !ok {
		return fmt.Errorf("no transition defined for event %s from state %s", event, f.currentState)
	}

	// Check guard if registered
	if guardsForState, ok := f.guards[f.currentState]; ok {
		if guard, ok := guardsForState[event]; ok {
			if !guard(args...) {
				return errors.New("transition denied by guard")
			}
		}
	}

	// Execute exit action of the current state
	if exitAction, ok := f.exitActions[f.currentState]; ok {
		if err := exitAction(args...); err != nil {
			return fmt.Errorf("exit action failed for state %s: %w", f.currentState, err)
		}
	}

	// Execute transition callback if registered
	if callbacksForState, ok := f.transitionCallbacks[f.currentState]; ok {
		if callback, ok := callbacksForState[event]; ok {
			if err := callback(args...); err != nil {
				return fmt.Errorf("transition callback failed for event %s from state %s: %w", event, f.currentState, err)
			}
		}
	}

	f.currentState = nextState

	// Execute entry action of the new state
	if entryAction, ok := f.entryActions[f.currentState]; ok {
		if err := entryAction(args...); err != nil {
			return fmt.Errorf("entry action failed for state %s: %w", f.currentState, err)
		}
	}

	return nil
}

// OnTransition registers a callback function to be executed when a specific transition occurs.
func (f *FSM) OnTransition(from State, event Event, callback Action) error {
	if _, ok := f.transitions[from]; !ok {
		return errors.New("no transitions defined from this state")
	}
	if _, ok := f.transitions[from][event]; !ok {
		return errors.New("no transition defined for this event from this state")
	}

	if _, ok := f.transitionCallbacks[from]; !ok {
		f.transitionCallbacks[from] = make(map[Event]Action)
	}
	f.transitionCallbacks[from][event] = callback
	return nil
}

// OnEntry registers an action to be executed when entering a state.
func (f *FSM) OnEntry(state State, action Action) {
	f.entryActions[state] = action
}

// OnExit registers an action to be executed when exiting a state.
func (f *FSM) OnExit(state State, action Action) {
	f.exitActions[state] = action
}

// AddGuard registers a guard function for a specific transition.
func (f *FSM) AddGuard(from State, event Event, guard Guard) error {
	if _, ok := f.transitions[from]; !ok {
		return errors.New("no transitions defined from this state")
	}
	if _, ok := f.transitions[from][event]; !ok {
		return errors.New("no transition defined for this event from this state")
	}

	if _, ok := f.guards[from]; !ok {
		f.guards[from] = make(map[Event]Guard)
	}
	f.guards[from][event] = guard
	return nil
}
