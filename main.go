package main

import (
	"context"
	"fmt"
	"go-fsm/ent"
	"go-fsm/ent/statemachine"
	"go-fsm/fsm"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
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
	// --- Database Configuration ---
	dbDriver := os.Getenv("DB_DRIVER")
	var client *ent.Client
	var err error

	switch dbDriver {
	case "mariadb":
		dsn := os.Getenv("DB_DSN")
		if dsn == "" {
			log.Fatal("DB_DRIVER is 'mariadb' but DB_DSN is not set. Please set the MariaDB DSN.")
		}
		client, err = ent.Open("mysql", dsn)
	default:
		fmt.Println("DB_DRIVER not set or invalid, using default: sqlite3")
		client, err = ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	}

	if err != nil {
		log.Fatalf("failed opening connection to database: %v", err)
	}
	defer client.Close()

	// 2. Run the auto migration tool to create the schema
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	fmt.Println("--- FSM with Ent Persistence and History ---")

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

	fmt.Println("\n3. Sending PUSH event (should succeed and lock)")
	err = turnstile.Transition(ctx, Push)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
	fmt.Printf("  Current state: %s\n", turnstile.CurrentState())

	// --- Query and display transition history ---
	fmt.Println("\n--- Querying Transition History ---")
	machine, err := client.StateMachine.Query().
		Where(statemachine.MachineIDEQ(machineID)).
		WithHistory().
		Only(ctx)
	if err != nil {
		log.Fatalf("failed to query machine history: %v", err)
	}

	for _, record := range machine.Edges.History {
		fmt.Printf("  - From: %s, To: %s, Event: %s, Time: %s\n",
			record.FromState, record.ToState, record.Event, record.Timestamp.Format("15:04:05"))
	}
}
