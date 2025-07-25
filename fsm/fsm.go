package fsm

// State represents a state in the state machine.
type State string

// Event represents an event that triggers a state transition.
type Event string

// Action is a function that is executed during a state transition.
type Action func()

// Guard is a function that returns true if a transition is allowed, false otherwise.
type Guard func() bool

// Transition defines a state transition, including the next state, an optional guard, and an optional action.
type Transition struct {
	NextState State
	Guard     Guard
	Action    Action
}

// StateMachine represents a finite state machine.
type StateMachine struct {
	CurrentState State
	Transitions  map[State]map[Event]Transition
}

// NewStateMachine creates a new state machine with initial state and transitions.
func NewStateMachine(initialState State, transitions map[State]map[Event]Transition) *StateMachine {
	return &StateMachine{
		CurrentState: initialState,
		Transitions:  transitions,
	}
}

// SendEvent processes an event and transitions the state machine.
func (sm *StateMachine) SendEvent(event Event) (State, bool) {
	if stateTransitions, ok := sm.Transitions[sm.CurrentState]; ok {
		if transition, ok := stateTransitions[event]; ok {
			// Check guard condition if present
			if transition.Guard != nil && !transition.Guard() {
				return sm.CurrentState, false // Guard condition not met
			}

			// Execute action if present
			if transition.Action != nil {
				transition.Action()
			}

			sm.CurrentState = transition.NextState
			return transition.NextState, true
		}
	}
	return sm.CurrentState, false // No transition for this event in the current state or guard failed
}
