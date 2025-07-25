package main

import (
	"fmt"
	"go-fsm/fsm" // Import the fsm package
)

// Define states for a turnstile
const (
	Locked   fsm.State = "LOCKED"
	Unlocked fsm.State = "UNLOCKED"
)

// Define events for a turnstile
const (
	Coin fsm.Event = "COIN"
	Push fsm.Event = "PUSH"
)

func main() {
	fmt.Println("--- FSM with OnEntry/OnExit, Guards, and Callbacks ---")

	// 1. Define all transitions
	transitions := []fsm.Transition{
		{From: Locked, Event: Coin, To: Unlocked},
		{From: Unlocked, Event: Push, To: Locked},
	}

	// 2. Create a new FSM instance
	turnstile, err := fsm.NewFSM("turnstile-01", Locked, transitions)
	if err != nil {
		fmt.Printf("Failed to create FSM: %v\n", err)
		return
	}

	// 3. Configure actions, guards, and callbacks
	turnstile.OnEntry(Unlocked, func(args ...interface{}) error {
		fmt.Println("  [OnEntry] Unlocked: Please pass through.")
		return nil
	})
	turnstile.OnExit(Locked, func(args ...interface{}) error {
		fmt.Println("  [OnExit] Locked: Processing payment...")
		return nil
	})
	turnstile.AddGuard(Locked, Coin, func(args ...interface{}) bool {
		fmt.Println("  [Guard] Checking if coin is valid... (approved)")
		return true
	})
	turnstile.OnTransition(Locked, Coin, func(args ...interface{}) error {
		fmt.Println("  [OnTransition] Coin transition is happening.")
		return nil
	})

	fmt.Printf("Initial state: %s\n", turnstile.CurrentState())

	// --- Simulate Events ---

	fmt.Println("\n1. Sending PUSH event (should fail, no transition)")
	err = turnstile.Transition(Push)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
	fmt.Printf("  Current state: %s\n", turnstile.CurrentState())

	fmt.Println("\n2. Sending COIN event (should succeed and unlock)")
	err = turnstile.Transition(Coin)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
	fmt.Printf("  Current state: %s\n", turnstile.CurrentState())

	fmt.Println("\n3. Sending COIN event again (should fail, no transition)")
	err = turnstile.Transition(Coin)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
	fmt.Printf("  Current state: %s\n", turnstile.CurrentState())

	fmt.Println("\n4. Sending PUSH event (should succeed and lock)")
	err = turnstile.Transition(Push)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
	fmt.Printf("  Current state: %s\n", turnstile.CurrentState())
}
