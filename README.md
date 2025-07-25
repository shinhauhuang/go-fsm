# Go FSM

A simple and flexible Finite State Machine (FSM) library for Go.

## Core Concepts

-   `fsm.State`: Represents a state in the machine (e.g., `Locked`).
-   `fsm.Event`: Represents an event that can trigger a transition (e.g., `Coin`).
-   `fsm.Transition`: A struct defining a transition rule: `{From, Event, To}`.
-   `fsm.Action`: A function `func(args ...interface{}) error` executed on entry, exit, or during a transition.
-   `fsm.Guard`: A function `func(args ...interface{}) bool` that must return `true` for a transition to be allowed.

## Database Configuration

This project uses `ent` for persistence and supports both **SQLite** and **MariaDB**. The database driver can be configured using a `.env` file or by setting environment variables directly.

A `.env.example` file is provided. To use it, copy it to `.env` and fill in your details:

```sh
cp .env.example .env
```

-   `DB_DRIVER`: Set to `mariadb` to use MariaDB. If unset or set to any other value, it will default to an in-memory SQLite database.
-   `DB_DSN`: When using `mariadb`, this **must** be set to your MariaDB Data Source Name (DSN). The application will exit if it's not set.
    -   Format: `user:password@tcp(host:port)/dbname?parseTime=true`

### Using Docker for Local Development

This project includes a `docker-compose.yml` file to easily spin up a local MariaDB and Adminer (a database management tool) instance.

**1. Set up your environment:**

Copy the example environment file and edit it if necessary. The defaults are configured to work with the `docker-compose.yml` file out of the box.

```sh
cp .env.example .env
```

**2. Start the services:**

```sh
docker-compose up -d
```

**3. Initialize the database schema:**

After the database container is running, run the initialization script. This only needs to be done once.

```sh
go run db/init.go
```

**4. Run the application:**

The Go application will now connect to the initialized MariaDB instance running in Docker.

```sh
go run main.go
```

**4. Access Adminer:**

You can manage the database by visiting `http://localhost:8080` in your browser.
-   **System**: `MariaDB`
-   **Server**: `db` (the service name from `docker-compose.yml`)
-   **Username**: `fsm_user` (from your `.env` file)
-   **Password**: `fsm_pass` (from your `.env` file)
-   **Database**: `fsm_db` (from your `.env` file)

**5. Stop the services:**

```sh
docker-compose down
```

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
