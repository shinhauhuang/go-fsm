package fsm

import (
	"context" // Import context for database operations
	"errors"
	"fmt"
	"sync" // Import the sync package for mutex

	"github.com/shinhauhuang/go-fsm/ent"              // Import the generated Ent client
	"github.com/shinhauhuang/go-fsm/ent/statemachine" // Import statemachine query
)

var (
	// ErrTransitionDenied is returned when a guard function denies a transition.
	ErrTransitionDenied = errors.New("transition denied by guard")
	// ErrInvalidTransition is returned for an invalid transition.
	ErrInvalidTransition = errors.New("invalid transition")
	// ErrInvalidEvent is returned for an invalid event.
	ErrInvalidEvent = errors.New("invalid event")
)

// State represents a state in the FSM.
type State string

// Event represents an event that triggers a state transition.
type Event string

// Action is a function that can be executed when entering or exiting a state.
// It receives a context, allowing for operations like database queries.
type Action func(ctx context.Context, args ...interface{}) error

// Guard is a function that determines if a transition is allowed.
// It receives a context, allowing for operations like database queries.
type Guard func(ctx context.Context, args ...interface{}) bool

// Transition defines a state transition: current state, event, and next state.
type Transition struct {
	From  State
	Event Event
	To    State
}

// FSM represents a Finite State Machine.
type FSM struct {
	mu                  sync.RWMutex // Mutex to ensure thread safety
	client              *ent.Client  // Ent client for persistence
	machineID           string       // Unique ID for this FSM instance
	currentState        State
	transitions         map[State]map[Event]State
	entryActions        map[State]Action
	exitActions         map[State]Action
	guards              map[State]map[Event]Guard
	transitionCallbacks map[State]map[Event]Action
}

// NewFSM creates a new FSM with an initial state, a list of transitions, and an Ent client for persistence.
// It will try to load the state from the database if a machineID is provided.
func NewFSM(ctx context.Context, client *ent.Client, machineID string, initialState State, transitions []Transition) (*FSM, error) {
	fsm := &FSM{
		client:              client,
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

	// Try to load the state from the database if machineID is provided
	if machineID != "" {
		sm, err := client.StateMachine.Query().Where(statemachine.MachineID(machineID)).Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				// If not found, create a new entry
				_, err := client.StateMachine.Create().
					SetMachineID(machineID).
					SetCurrentState(string(initialState)).
					Save(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to create new state machine entry: %w", err)
				}
			} else {
				return nil, fmt.Errorf("failed to query state machine: %w", err)
			}
		} else {
			fsm.currentState = State(sm.CurrentState)
		}
	}

	return fsm, nil
}

// CurrentState returns the current state of the FSM.
func (f *FSM) CurrentState() State {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.currentState
}

// Transition attempts to transition the FSM to a new state based on an event.
func (f *FSM) Transition(ctx context.Context, event Event, args ...interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	nextStates, ok := f.transitions[f.currentState]
	if !ok {
		return fmt.Errorf("%w: no transitions defined from state %s", ErrInvalidTransition, f.currentState)
	}

	nextState, ok := nextStates[event]
	if !ok {
		return fmt.Errorf("%w: no transition defined for event %s from state %s", ErrInvalidEvent, event, f.currentState)
	}

	// Check guard if registered
	if guardsForState, ok := f.guards[f.currentState]; ok {
		if guard, ok := guardsForState[event]; ok {
			if !guard(ctx, args...) {
				return ErrTransitionDenied
			}
		}
	}

	// Execute exit action of the current state
	if exitAction, ok := f.exitActions[f.currentState]; ok {
		if err := exitAction(ctx, args...); err != nil {
			return fmt.Errorf("exit action failed for state %s: %w", f.currentState, err)
		}
	}

	// Execute transition callback if registered
	if callbacksForState, ok := f.transitionCallbacks[f.currentState]; ok {
		if callback, ok := callbacksForState[event]; ok {
			if err := callback(ctx, args...); err != nil {
				return fmt.Errorf("transition callback failed for event %s from state %s: %w", event, f.currentState, err)
			}
		}
	}

	previousState := f.currentState
	f.currentState = nextState

	// Execute entry action of the new state
	if entryAction, ok := f.entryActions[f.currentState]; ok {
		if err := entryAction(ctx, args...); err != nil {
			// Revert state if entry action fails
			f.currentState = previousState
			return fmt.Errorf("entry action failed for state %s: %w", f.currentState, err)
		}
	}

	// Persist the new state and the transition history to the database
	if f.client != nil && f.machineID != "" {
		// Use a transaction to ensure atomicity
		tx, err := f.client.Tx(ctx)
		if err != nil {
			f.currentState = previousState // Revert state
			return fmt.Errorf("failed to start transaction: %w", err)
		}

		// Get the StateMachine node
		sm, err := tx.StateMachine.Query().Where(statemachine.MachineID(f.machineID)).Only(ctx)
		if err != nil {
			f.currentState = previousState // Revert state
			return fmt.Errorf("failed to query state machine for update: %w", err)
		}

		// Create the history record
		_, err = tx.StateTransition.Create().
			SetFromState(string(previousState)).
			SetToState(string(nextState)).
			SetEvent(string(event)).
			SetMachine(sm).
			Save(ctx)
		if err != nil {
			f.currentState = previousState // Revert state
			if rerr := tx.Rollback(); rerr != nil {
				return fmt.Errorf("failed to rollback transaction: %v, original error: %w", rerr, err)
			}
			return fmt.Errorf("failed to create transition history: %w", err)
		}

		// Update the current state of the machine
		_, err = sm.Update().SetCurrentState(string(f.currentState)).Save(ctx)
		if err != nil {
			f.currentState = previousState // Revert state
			if rerr := tx.Rollback(); rerr != nil {
				return fmt.Errorf("failed to rollback transaction: %v, original error: %w", rerr, err)
			}
			return fmt.Errorf("failed to persist state: %w", err)
		}

		if err := tx.Commit(); err != nil {
			f.currentState = previousState // Revert state
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	return nil
}

// OnTransition registers a callback function to be executed when a specific transition occurs.
func (f *FSM) OnTransition(from State, event Event, callback Action) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.transitions[from]; !ok {
		return fmt.Errorf("%w: no transitions defined from state %s", ErrInvalidTransition, from)
	}
	if _, ok := f.transitions[from][event]; !ok {
		return fmt.Errorf("%w: no transition defined for event %s from state %s", ErrInvalidEvent, event, from)
	}

	if _, ok := f.transitionCallbacks[from]; !ok {
		f.transitionCallbacks[from] = make(map[Event]Action)
	}
	f.transitionCallbacks[from][event] = callback
	return nil
}

// OnEntry registers an action to be executed when entering a state.
func (f *FSM) OnEntry(state State, action Action) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.entryActions[state] = action
}

// OnExit registers an action to be executed when exiting a state.
func (f *FSM) OnExit(state State, action Action) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.exitActions[state] = action
}

// AddGuard registers a guard function for a specific transition.
func (f *FSM) AddGuard(from State, event Event, guard Guard) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.transitions[from]; !ok {
		return fmt.Errorf("%w: no transitions defined from state %s", ErrInvalidTransition, from)
	}
	if _, ok := f.transitions[from][event]; !ok {
		return fmt.Errorf("%w: no transition defined for event %s from state %s", ErrInvalidEvent, event, from)
	}

	if _, ok := f.guards[from]; !ok {
		f.guards[from] = make(map[Event]Guard)
	}
	f.guards[from][event] = guard
	return nil
}
