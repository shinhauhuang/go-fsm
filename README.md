# Go FSM

A simple and flexible Finite State Machine (FSM) library for Go.

## Core Concepts

-   `fsm.State`: Represents a state in the machine (e.g., `Locked`).
-   `fsm.Event`: Represents an event that can trigger a transition (e.g., `Coin`).
-   `fsm.Transition`: A struct defining a transition rule: `{From, Event, To}`.
-   `fsm.Action`: A function `func(args ...interface{}) error` executed on entry, exit, or during a transition.
-   `fsm.Guard`: A function `func(args ...interface{}) bool` that must return `true` for a transition to be allowed.

## How to Use

### 1. Define Transitions

Create a slice of `fsm.Transition` to define the state machine's rules.

```go
transitions := []fsm.Transition{
    {From: Locked, Event: Coin, To: Unlocked},
    {From: Unlocked, Event: Push, To: Locked},
}
```

### 2. Create an FSM Instance

Use `fsm.NewFSM` to create a new state machine instance.

```go
machine, err := fsm.NewFSM("machine-id", initialState, transitions)
```

### 3. Configure Callbacks and Guards

Use the provided methods to attach your custom logic to the FSM.

-   `OnEntry(state, action)`: Executes `action` when entering `state`.
-   `OnExit(state, action)`: Executes `action` when exiting `state`.
-   `AddGuard(from, event, guard)`: Executes `guard` before the transition from `from` state on `event`. The transition is denied if the guard returns `false`.
-   `OnTransition(from, event, action)`: Executes `action` during the transition from `from` state on `event`.

```go
// Example: Add an entry action for the "Unlocked" state
machine.OnEntry(Unlocked, func(args ...interface{}) error {
    fmt.Println("State is now Unlocked")
    return nil
})

// Example: Add a guard for the "Coin" event in the "Locked" state
machine.AddGuard(Locked, Coin, func(args ...interface{}) bool {
    // ... logic to validate coin ...
    return true
})
```

### 4. Trigger a Transition

Use the `Transition` method to send an event to the state machine.

```go
err := machine.Transition(Coin)
if err != nil {
    // Handle transition error (e.g., guard denied, no transition found)
}
```

### 5. Get the Current State

Use the `CurrentState` method to check the machine's current state at any time.

```go
currentState := machine.CurrentState()
fmt.Printf("The current state is: %s\n", currentState)
