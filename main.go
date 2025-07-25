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

// Example Action functions
func onUnlock() {
	fmt.Println("Action: Turnstile unlocked! Please enter.")
}

func onLock() {
	fmt.Println("Action: Turnstile locked. Insert coin.")
}

// Example Guard functions
func canPush() bool {
	fmt.Println("Guard: Checking if push is allowed...")
	// In a real scenario, this might check if a person is actually pushing
	return true
}

func canCoin() bool {
	fmt.Println("Guard: Checking if coin is valid...")
	// In a real scenario, this might validate the coin
	return true
}

func main() {
	// Define turnstile transitions with actions and guards
	transitions := map[fsm.State]map[fsm.Event]fsm.Transition{
		Locked: {
			Coin: {
				NextState: Unlocked,
				Guard:     canCoin,
				Action:    onUnlock,
			},
		},
		Unlocked: {
			Push: {
				NextState: Locked,
				Guard:     canPush,
				Action:    onLock,
			},
		},
	}

	// Create a new turnstile state machine
	turnstile := fsm.NewStateMachine(Locked, transitions)

	fmt.Printf("Initial state: %s\n", turnstile.CurrentState)

	// Simulate events
	fmt.Println("\n--- Simulating PUSH event (should not change state from LOCKED) ---")
	newState, transitioned := turnstile.SendEvent(Push)
	fmt.Printf("Current state: %s, Transitioned: %t\n", newState, transitioned)

	fmt.Println("\n--- Simulating COIN event (should unlock) ---")
	newState, transitioned = turnstile.SendEvent(Coin)
	fmt.Printf("Current state: %s, Transitioned: %t\n", newState, transitioned)

	fmt.Println("\n--- Simulating COIN event (should not change state from UNLOCKED) ---")
	newState, transitioned = turnstile.SendEvent(Coin)
	fmt.Printf("Current state: %s, Transitioned: %t\n", newState, transitioned)

	fmt.Println("\n--- Simulating PUSH event (should lock) ---")
	newState, transitioned = turnstile.SendEvent(Push)
	fmt.Printf("Current state: %s, Transitioned: %t\n", newState, transitioned)

	fmt.Println("\n--- Simulating PUSH event (should not change state from LOCKED) ---")
	newState, transitioned = turnstile.SendEvent(Push)
	fmt.Printf("Current state: %s, Transitioned: %t\n", newState, transitioned)
}
