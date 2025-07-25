package main

import (
	"context"
	"fmt"
	"go-fsm/ent"
	"go-fsm/fsm"
	"log"

	_ "github.com/mattn/go-sqlite3"
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
	// 1. Create an in-memory SQLite ent client
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	// 2. Run the auto migration tool to create the schema
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	fmt.Println("--- FSM with Ent Persistence ---")

	// 3. Define all transitions
	transitions := []fsm.Transition{
		{From: Locked, Event: Coin, To: Unlocked},
		{From: Unlocked, Event: Push, To: Locked},
	}

	// 4. Create a new FSM instance with the ent client
	machineID := "turnstile-01"
	turnstile, err := fsm.NewFSM(ctx, client, machineID, Locked, transitions)
	if err != nil {
		log.Fatalf("Failed to create FSM: %v", err)
	}

	// 5. Configure actions, guards, and callbacks
	turnstile.OnEntry(Unlocked, func(args ...interface{}) error {
		fmt.Println("  [OnEntry] Unlocked: Please pass through.")
		return nil
	})
	turnstile.AddGuard(Locked, Coin, func(args ...interface{}) bool {
		fmt.Println("  [Guard] Checking if coin is valid... (approved)")
		return true
	})

	fmt.Printf("Initial state (from DB or initial): %s\n", turnstile.CurrentState())

	// --- Simulate Events ---

	fmt.Println("\n1. Sending PUSH event (should fail, no transition)")
	err = turnstile.Transition(ctx, Push)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
	fmt.Printf("  Current state: %s\n", turnstile.CurrentState())

	fmt.Println("\n2. Sending COIN event (should succeed and unlock)")
	err = turnstile.Transition(ctx, Coin)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
	fmt.Printf("  Current state: %s\n", turnstile.CurrentState())

	// --- Create a new FSM instance with the same ID to test loading from DB ---
	fmt.Println("\n--- Re-creating FSM with same ID to test persistence ---")
	turnstile2, err := fsm.NewFSM(ctx, client, machineID, Locked, transitions)
	if err != nil {
		log.Fatalf("Failed to create second FSM instance: %v", err)
	}
	fmt.Printf("State loaded from DB: %s\n", turnstile2.CurrentState())
}
