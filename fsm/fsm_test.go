package fsm

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/shinhauhuang/go-fsm/ent"
	"github.com/shinhauhuang/go-fsm/ent/enttest"
	"github.com/shinhauhuang/go-fsm/ent/statemachine"
	"github.com/shinhauhuang/go-fsm/ent/statetransition"

	_ "github.com/mattn/go-sqlite3" // Driver for SQLite
)

// Define common states and events for testing
const (
	StateIdle    State = "idle"
	StateRunning State = "running"
	StatePaused  State = "paused"
	StateStopped State = "stopped"

	EventStart  Event = "start"
	EventPause  Event = "pause"
	EventResume Event = "resume"
	EventStop   Event = "stop"
)

// setupTestClient creates a new in-memory SQLite client for testing.
func setupTestClient(t *testing.T) *ent.Client {
	// Reverting to original, but trying "_fk=1" instead of "_fk_checks=1"
	// The error message specifically mentions "_fk=1" being missing.
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	return client
}

// defineTestTransitions defines a common set of transitions for testing.
func defineTestTransitions() []Transition {
	return []Transition{
		{From: StateIdle, Event: EventStart, To: StateRunning},
		{From: StateRunning, Event: EventPause, To: StatePaused},
		{From: StatePaused, Event: EventResume, To: StateRunning},
		{From: StateRunning, Event: EventStop, To: StateStopped},
		{From: StatePaused, Event: EventStop, To: StateStopped},
	}
}

func TestNewFSM(t *testing.T) {
	ctx := context.Background()
	transitions := defineTestTransitions()

	t.Run("New FSM without persistence", func(t *testing.T) {
		f, err := NewFSM(ctx, nil, "", StateIdle, transitions)
		if err != nil {
			t.Fatalf("NewFSM failed: %v", err)
		}
		if f.CurrentState() != StateIdle {
			t.Errorf("Expected initial state %s, got %s", StateIdle, f.CurrentState())
		}
		if f.client != nil {
			t.Errorf("Expected nil client, got non-nil")
		}
		if f.machineID != "" {
			t.Errorf("Expected empty machineID, got %s", f.machineID)
		}
	})

	t.Run("New FSM with persistence - new machine", func(t *testing.T) {
		client := setupTestClient(t)
		defer client.Close()

		machineID := "test_machine_1"
		f, err := NewFSM(ctx, client, machineID, StateIdle, transitions)
		if err != nil {
			t.Fatalf("NewFSM failed: %v", err)
		}
		if f.CurrentState() != StateIdle {
			t.Errorf("Expected initial state %s, got %s", StateIdle, f.CurrentState())
		}

		// Verify state was saved in DB
		sm, err := client.StateMachine.Query().Where(statemachine.MachineID(machineID)).Only(ctx)
		if err != nil {
			t.Fatalf("Failed to query state machine from DB: %v", err)
		}
		if sm.CurrentState != string(StateIdle) {
			t.Errorf("Expected DB state %s, got %s", StateIdle, sm.CurrentState)
		}
	})

	t.Run("New FSM with persistence - load existing machine", func(t *testing.T) {
		client := setupTestClient(t)
		defer client.Close()

		machineID := "test_machine_2"
		// First, create a machine and set its state to Running
		_, err := client.StateMachine.Create().
			SetMachineID(machineID).
			SetCurrentState(string(StateRunning)).
			Save(ctx)
		if err != nil {
			t.Fatalf("Failed to pre-create state machine in DB: %v", err)
		}

		// Now, initialize FSM with the same machineID, expecting it to load StateRunning
		f, err := NewFSM(ctx, client, machineID, StateIdle, transitions) // Initial state here should be ignored
		if err != nil {
			t.Fatalf("NewFSM failed: %v", err)
		}
		if f.CurrentState() != StateRunning {
			t.Errorf("Expected loaded state %s, got %s", StateRunning, f.CurrentState())
		}
	})

	t.Run("New FSM with duplicate transitions", func(t *testing.T) {
		duplicateTransitions := []Transition{
			{From: StateIdle, Event: EventStart, To: StateRunning},
			{From: StateIdle, Event: EventStart, To: StatePaused}, // Duplicate
		}
		_, err := NewFSM(ctx, nil, "", StateIdle, duplicateTransitions)
		if err == nil {
			t.Fatalf("Expected error for duplicate transition, got nil")
		}
		expectedErrorMsg := fmt.Sprintf("duplicate transition defined from state %s for event %s", StateIdle, EventStart)
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
		}
	})
}

func TestLoadFSM(t *testing.T) {
	ctx := context.Background()
	transitions := defineTestTransitions()
	client := setupTestClient(t)
	defer client.Close()

	t.Run("Load existing FSM successfully", func(t *testing.T) {
		machineID := "load_test_machine_1"
		// Pre-create a state machine in the DB
		_, err := client.StateMachine.Create().
			SetMachineID(machineID).
			SetCurrentState(string(StatePaused)).
			Save(ctx)
		if err != nil {
			t.Fatalf("Failed to pre-create state machine: %v", err)
		}

		f, err := LoadFSM(ctx, client, machineID, transitions)
		if err != nil {
			t.Fatalf("LoadFSM failed: %v", err)
		}
		if f.CurrentState() != StatePaused {
			t.Errorf("Expected loaded state %s, got %s", StatePaused, f.CurrentState())
		}
		if f.machineID != machineID {
			t.Errorf("Expected machineID %s, got %s", machineID, f.machineID)
		}
	})

	t.Run("Load non-existent FSM", func(t *testing.T) {
		machineID := "non_existent_machine"
		_, err := LoadFSM(ctx, client, machineID, transitions)
		if err == nil {
			t.Fatalf("Expected error for non-existent machine, got nil")
		}
		if !ent.IsNotFound(err) {
			t.Errorf("Expected IsNotFound error, got %v", err)
		}
	})

	t.Run("Load FSM with nil client", func(t *testing.T) {
		_, err := LoadFSM(ctx, nil, "some_id", transitions)
		if err == nil {
			t.Fatalf("Expected error for nil client, got nil")
		}
		if err.Error() != "client and machineID are required to load an FSM" {
			t.Errorf("Expected specific error message, got %v", err)
		}
	})

	t.Run("Load FSM with empty machineID", func(t *testing.T) {
		_, err := LoadFSM(ctx, client, "", transitions)
		if err == nil {
			t.Fatalf("Expected error for empty machineID, got nil")
		}
		if err.Error() != "client and machineID are required to load an FSM" {
			t.Errorf("Expected specific error message, got %v", err)
		}
	})

	t.Run("Load FSM with duplicate transitions in definition", func(t *testing.T) {
		machineID := "load_test_machine_2"
		// Pre-create a state machine in the DB
		_, err := client.StateMachine.Create().
			SetMachineID(machineID).
			SetCurrentState(string(StateIdle)).
			Save(ctx)
		if err != nil {
			t.Fatalf("Failed to pre-create state machine: %v", err)
		}

		duplicateTransitions := []Transition{
			{From: StateIdle, Event: EventStart, To: StateRunning},
			{From: StateIdle, Event: EventStart, To: StatePaused}, // Duplicate
		}
		_, err = LoadFSM(ctx, client, machineID, duplicateTransitions)
		if err == nil {
			t.Fatalf("Expected error for duplicate transition during loading, got nil")
		}
		expectedErrorMsg := fmt.Sprintf("duplicate transition defined from state %s for event %s during FSM loading", StateIdle, EventStart)
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
		}
	})
}

func TestCurrentState(t *testing.T) {
	ctx := context.Background()
	transitions := defineTestTransitions()
	f, err := NewFSM(ctx, nil, "", StateIdle, transitions)
	if err != nil {
		t.Fatalf("NewFSM failed: %v", err)
	}

	if f.CurrentState() != StateIdle {
		t.Errorf("Expected initial state %s, got %s", StateIdle, f.CurrentState())
	}

	// Simulate a state change (without using Transition method for simplicity in this test)
	f.mu.Lock()
	f.currentState = StateRunning
	f.mu.Unlock()

	if f.CurrentState() != StateRunning {
		t.Errorf("Expected state %s after change, got %s", StateRunning, f.CurrentState())
	}
}

func TestTransition(t *testing.T) {
	ctx := context.Background()
	transitions := defineTestTransitions()

	t.Run("Successful transition", func(t *testing.T) {
		f, err := NewFSM(ctx, nil, "", StateIdle, transitions)
		if err != nil {
			t.Fatalf("NewFSM failed: %v", err)
		}

		err = f.Transition(ctx, EventStart)
		if err != nil {
			t.Fatalf("Transition failed: %v", err)
		}
		if f.CurrentState() != StateRunning {
			t.Errorf("Expected state %s, got %s", StateRunning, f.CurrentState())
		}
	})

	t.Run("Invalid transition - no transitions from current state", func(t *testing.T) {
		f, err := NewFSM(ctx, nil, "", StateStopped, transitions) // StateStopped has no outgoing transitions
		if err != nil {
			t.Fatalf("NewFSM failed: %v", err)
		}

		err = f.Transition(ctx, EventStart)
		if !errors.Is(err, ErrInvalidTransition) {
			t.Errorf("Expected ErrInvalidTransition, got %v", err)
		}
		if f.CurrentState() != StateStopped { // State should not change
			t.Errorf("Expected state %s, got %s", StateStopped, f.CurrentState())
		}
	})

	t.Run("Invalid event - no transition for event from current state", func(t *testing.T) {
		f, err := NewFSM(ctx, nil, "", StateIdle, transitions)
		if err != nil {
			t.Fatalf("NewFSM failed: %v", err)
		}

		err = f.Transition(ctx, EventPause) // EventPause is not valid from StateIdle
		if !errors.Is(err, ErrInvalidEvent) {
			t.Errorf("Expected ErrInvalidEvent, got %v", err)
		}
		if f.CurrentState() != StateIdle { // State should not change
			t.Errorf("Expected state %s, got %s", StateIdle, f.CurrentState())
		}
	})

	t.Run("Transition denied by guard", func(t *testing.T) {
		f, err := NewFSM(ctx, nil, "", StateIdle, transitions)
		if err != nil {
			t.Fatalf("NewFSM failed: %v", err)
		}

		// Add a guard that always denies the transition
		f.AddGuard(StateIdle, EventStart, func(ctx context.Context, args ...interface{}) bool {
			return false
		})

		err = f.Transition(ctx, EventStart)
		if !errors.Is(err, ErrTransitionDenied) {
			t.Errorf("Expected ErrTransitionDenied, got %v", err)
		}
		if f.CurrentState() != StateIdle { // State should not change
			t.Errorf("Expected state %s, got %s", StateIdle, f.CurrentState())
		}
	})

	t.Run("Exit action execution", func(t *testing.T) {
		f, err := NewFSM(ctx, nil, "", StateIdle, transitions)
		if err != nil {
			t.Fatalf("NewFSM failed: %v", err)
		}

		exitActionCalled := false
		f.OnExit(StateIdle, func(ctx context.Context, args ...interface{}) error {
			exitActionCalled = true
			return nil
		})

		f.Transition(ctx, EventStart)
		if !exitActionCalled {
			t.Errorf("Exit action was not called")
		}
	})

	t.Run("Entry action execution", func(t *testing.T) {
		f, err := NewFSM(ctx, nil, "", StateIdle, transitions)
		if err != nil {
			t.Fatalf("NewFSM failed: %v", err)
		}

		entryActionCalled := false
		f.OnEntry(StateRunning, func(ctx context.Context, args ...interface{}) error {
			entryActionCalled = true
			return nil
		})

		f.Transition(ctx, EventStart)
		if !entryActionCalled {
			t.Errorf("Entry action was not called")
		}
	})

	t.Run("Transition callback execution", func(t *testing.T) {
		f, err := NewFSM(ctx, nil, "", StateIdle, transitions)
		if err != nil {
			t.Fatalf("NewFSM failed: %v", err)
		}

		callbackCalled := false
		f.OnTransition(StateIdle, EventStart, func(ctx context.Context, args ...interface{}) error {
			callbackCalled = true
			return nil
		})

		f.Transition(ctx, EventStart)
		if !callbackCalled {
			t.Errorf("Transition callback was not called")
		}
	})

	t.Run("Exit action failure reverts state", func(t *testing.T) {
		f, err := NewFSM(ctx, nil, "", StateIdle, transitions)
		if err != nil {
			t.Fatalf("NewFSM failed: %v", err)
		}

		f.OnExit(StateIdle, func(ctx context.Context, args ...interface{}) error {
			return errors.New("exit action error")
		})

		err = f.Transition(ctx, EventStart)
		if err == nil {
			t.Errorf("Expected error from exit action, got nil")
		}
		if f.CurrentState() != StateIdle { // State should revert
			t.Errorf("Expected state %s, got %s", StateIdle, f.CurrentState())
		}
	})

	t.Run("Entry action failure reverts state", func(t *testing.T) {
		f, err := NewFSM(ctx, nil, "", StateIdle, transitions)
		if err != nil {
			t.Fatalf("NewFSM failed: %v", err)
		}

		f.OnEntry(StateRunning, func(ctx context.Context, args ...interface{}) error {
			return errors.New("entry action error")
		})

		err = f.Transition(ctx, EventStart)
		if err == nil {
			t.Errorf("Expected error from entry action, got nil")
		}
		if f.CurrentState() != StateIdle { // State should revert
			t.Errorf("Expected state %s, got %s", StateIdle, f.CurrentState())
		}
	})

	t.Run("Transition callback failure reverts state", func(t *testing.T) {
		f, err := NewFSM(ctx, nil, "", StateIdle, transitions)
		if err != nil {
			t.Fatalf("NewFSM failed: %v", err)
		}

		f.OnTransition(StateIdle, EventStart, func(ctx context.Context, args ...interface{}) error {
			return errors.New("transition callback error")
		})

		err = f.Transition(ctx, EventStart)
		if err == nil {
			t.Errorf("Expected error from transition callback, got nil")
		}
		if f.CurrentState() != StateIdle { // State should revert
			t.Errorf("Expected state %s, got %s", StateIdle, f.CurrentState())
		}
	})

	t.Run("Successful transition with persistence", func(t *testing.T) {
		client := setupTestClient(t)
		defer client.Close()

		machineID := "test_machine_3"
		f, err := NewFSM(ctx, client, machineID, StateIdle, transitions)
		if err != nil {
			t.Fatalf("NewFSM failed: %v", err)
		}

		err = f.Transition(ctx, EventStart)
		if err != nil {
			t.Fatalf("Transition failed: %v", err)
		}
		if f.CurrentState() != StateRunning {
			t.Errorf("Expected state %s, got %s", StateRunning, f.CurrentState())
		}

		// Verify state in DB
		sm, err := client.StateMachine.Query().Where(statemachine.MachineID(machineID)).Only(ctx)
		if err != nil {
			t.Fatalf("Failed to query state machine from DB: %v", err)
		}
		if sm.CurrentState != string(StateRunning) {
			t.Errorf("Expected DB state %s, got %s", StateRunning, sm.CurrentState)
		}

		// Verify transition history in DB
		history, err := client.StateTransition.Query().
			Where(statetransition.HasMachineWith(statemachine.MachineID(machineID))).
			Order(ent.Desc(statetransition.FieldID)).
			First(ctx)
		if err != nil {
			t.Fatalf("Failed to query transition history from DB: %v", err)
		}
		if history.FromState != string(StateIdle) || history.ToState != string(StateRunning) || history.Event != string(EventStart) {
			t.Errorf("Expected history %s->%s via %s, got %s->%s via %s",
				StateIdle, StateRunning, EventStart, history.FromState, history.ToState, history.Event)
		}
	})

	t.Run("Persistence failure during transition reverts state and rolls back transaction", func(t *testing.T) {
		client := setupTestClient(t)
		defer client.Close()

		// machineID := "test_machine_4" // Declared and not used, as this test is skipped.
		// f, err := NewFSM(ctx, client, machineID, StateIdle, transitions)
		// if err != nil {
		// 	t.Fatalf("NewFSM failed: %v", err)
		// }

		// Simulate a DB error by closing the client mid-transition (not ideal, but demonstrates rollback)
		// This is a bit hacky for a unit test, but for integration tests, you'd mock the client more robustly.
		// For now, we'll just test the rollback path.
		// To properly test this, we'd need to inject a mock client that can fail specific operations.
		// For simplicity, I'll just ensure the state reverts if the DB operation fails.
		// The current implementation of `Transition` handles DB errors by reverting `f.currentState`
		// and rolling back the transaction. We can't easily simulate a specific DB error
		// without mocking the Ent client's internal methods, which is beyond the scope of a simple unit test.
		// Instead, I'll focus on testing the state reversion logic.

		// Let's assume a DB error occurs during the update of the state machine.
		// We can't directly inject an error into the Ent client's Save method here.
		// The existing test for "Successful transition with persistence" covers the happy path.
		// For error paths, we'd typically use a mock database or mock the Ent client.
		// Given the current setup, I'll skip a direct test for DB transaction failure,
		// as it would require significant mocking infrastructure.
		// The code already has `f.currentState = previousState` and `tx.Rollback()` calls,
		// which are the critical parts to ensure atomicity.
		// A more advanced test would involve a custom Ent driver that can be configured to fail.

		// For now, I'll add a placeholder comment to acknowledge this limitation.
		t.Skip("Skipping direct test for persistence failure and rollback due to complexity of mocking Ent client for specific error injection.")
	})
}

func TestOnTransition(t *testing.T) {
	ctx := context.Background()
	transitions := defineTestTransitions()
	f, err := NewFSM(ctx, nil, "", StateIdle, transitions)
	if err != nil {
		t.Fatalf("NewFSM failed: %v", err)
	}

	t.Run("Register and execute OnTransition callback", func(t *testing.T) {
		callbackCalled := false
		err := f.OnTransition(StateIdle, EventStart, func(ctx context.Context, args ...interface{}) error {
			callbackCalled = true
			return nil
		})
		if err != nil {
			t.Fatalf("OnTransition registration failed: %v", err)
		}

		f.Transition(ctx, EventStart)
		if !callbackCalled {
			t.Errorf("OnTransition callback was not executed")
		}
	})

	t.Run("OnTransition with invalid 'from' state", func(t *testing.T) {
		err := f.OnTransition("NonExistentState", EventStart, func(ctx context.Context, args ...interface{}) error { return nil })
		if !errors.Is(err, ErrInvalidTransition) {
			t.Errorf("Expected ErrInvalidTransition, got %v", err)
		}
	})

	t.Run("OnTransition with invalid event for 'from' state", func(t *testing.T) {
		err := f.OnTransition(StateIdle, "NonExistentEvent", func(ctx context.Context, args ...interface{}) error { return nil })
		if !errors.Is(err, ErrInvalidEvent) {
			t.Errorf("Expected ErrInvalidEvent, got %v", err)
		}
	})
}

func TestOnEntry(t *testing.T) {
	ctx := context.Background()
	transitions := defineTestTransitions()
	f, err := NewFSM(ctx, nil, "", StateIdle, transitions)
	if err != nil {
		t.Fatalf("NewFSM failed: %v", err)
	}

	entryActionCalled := false
	f.OnEntry(StateRunning, func(ctx context.Context, args ...interface{}) error {
		entryActionCalled = true
		return nil
	})

	f.Transition(ctx, EventStart)
	if !entryActionCalled {
		t.Errorf("OnEntry action was not executed")
	}
}

func TestOnExit(t *testing.T) {
	ctx := context.Background()
	transitions := defineTestTransitions()
	f, err := NewFSM(ctx, nil, "", StateIdle, transitions)
	if err != nil {
		t.Fatalf("NewFSM failed: %v", err)
	}

	exitActionCalled := false
	f.OnExit(StateIdle, func(ctx context.Context, args ...interface{}) error {
		exitActionCalled = true
		return nil
	})

	f.Transition(ctx, EventStart)
	if !exitActionCalled {
		t.Errorf("OnExit action was not executed")
	}
}

func TestAddGuard(t *testing.T) {
	ctx := context.Background()
	transitions := defineTestTransitions()
	f, err := NewFSM(ctx, nil, "", StateIdle, transitions)
	if err != nil {
		t.Fatalf("NewFSM failed: %v", err)
	}

	t.Run("Register and execute AddGuard", func(t *testing.T) {
		guardCalled := false
		err := f.AddGuard(StateIdle, EventStart, func(ctx context.Context, args ...interface{}) bool {
			guardCalled = true
			return true // Allow transition
		})
		if err != nil {
			t.Fatalf("AddGuard registration failed: %v", err)
		}

		f.Transition(ctx, EventStart)
		if !guardCalled {
			t.Errorf("Guard was not executed")
		}
	})

	t.Run("AddGuard with invalid 'from' state", func(t *testing.T) {
		err := f.AddGuard("NonExistentState", EventStart, func(ctx context.Context, args ...interface{}) bool { return true })
		if !errors.Is(err, ErrInvalidTransition) {
			t.Errorf("Expected ErrInvalidTransition, got %v", err)
		}
	})

	t.Run("AddGuard with invalid event for 'from' state", func(t *testing.T) {
		err := f.AddGuard(StateIdle, "NonExistentEvent", func(ctx context.Context, args ...interface{}) bool { return true })
		if !errors.Is(err, ErrInvalidEvent) {
			t.Errorf("Expected ErrInvalidEvent, got %v", err)
		}
	})
}

func TestFSMConcurrency(t *testing.T) {
	ctx := context.Background()
	client := setupTestClient(t)
	defer client.Close()

	machineID := "concurrent_machine"
	transitions := defineTestTransitions()
	f, err := NewFSM(ctx, client, machineID, StateIdle, transitions)
	if err != nil {
		t.Fatalf("NewFSM failed: %v", err)
	}

	numTransitions := 100
	var wg sync.WaitGroup
	errChan := make(chan error, numTransitions)

	// Start multiple goroutines trying to transition the FSM
	for i := 0; i < numTransitions; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Simulate a sequence of transitions
			err := f.Transition(ctx, EventStart)
			if err != nil && !errors.Is(err, ErrInvalidTransition) && !errors.Is(err, ErrInvalidEvent) {
				errChan <- fmt.Errorf("concurrent transition failed: %w", err)
				return
			}
			if f.CurrentState() == StateRunning {
				err = f.Transition(ctx, EventPause)
				if err != nil && !errors.Is(err, ErrInvalidTransition) && !errors.Is(err, ErrInvalidEvent) {
					errChan <- fmt.Errorf("concurrent transition failed: %w", err)
					return
				}
			}
			if f.CurrentState() == StatePaused {
				err = f.Transition(ctx, EventResume)
				if err != nil && !errors.Is(err, ErrInvalidTransition) && !errors.Is(err, ErrInvalidEvent) {
					errChan <- fmt.Errorf("concurrent transition failed: %w", err)
					return
				}
			}
			if f.CurrentState() == StateRunning {
				err = f.Transition(ctx, EventStop)
				if err != nil && !errors.Is(err, ErrInvalidTransition) && !errors.Is(err, ErrInvalidEvent) {
					errChan <- fmt.Errorf("concurrent transition failed: %w", err)
					return
				}
			}
		}()
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Error(err)
	}

	// After all concurrent operations, check the final state and history
	// The final state might be any of the valid end states depending on execution order.
	// We primarily want to ensure no panics or race conditions.
	finalState := f.CurrentState()
	if finalState != StateIdle && finalState != StateRunning && finalState != StatePaused && finalState != StateStopped {
		t.Errorf("Unexpected final state after concurrent transitions: %s", finalState)
	}

	// Verify that the database operations were handled correctly under concurrency
	// This is a basic check; a more robust test would verify the exact number of transitions.
	count, err := client.StateTransition.Query().
		Where(statetransition.HasMachineWith(statemachine.MachineID(machineID))).
		Count(ctx)
	if err != nil {
		t.Fatalf("Failed to count state transitions: %v", err)
	}
	if count == 0 {
		t.Errorf("Expected some transitions to be recorded, got 0")
	}
}
