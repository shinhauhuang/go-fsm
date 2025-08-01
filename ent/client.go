// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/shinhauhuang/go-fsm/ent/migrate"

	"github.com/shinhauhuang/go-fsm/ent/statemachine"
	"github.com/shinhauhuang/go-fsm/ent/statetransition"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// StateMachine is the client for interacting with the StateMachine builders.
	StateMachine *StateMachineClient
	// StateTransition is the client for interacting with the StateTransition builders.
	StateTransition *StateTransitionClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	client := &Client{config: newConfig(opts...)}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.StateMachine = NewStateMachineClient(c.config)
	c.StateTransition = NewStateTransitionClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// newConfig creates a new config for the client.
func newConfig(opts ...Option) config {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	return cfg
}

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:             ctx,
		config:          cfg,
		StateMachine:    NewStateMachineClient(cfg),
		StateTransition: NewStateTransitionClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:             ctx,
		config:          cfg,
		StateMachine:    NewStateMachineClient(cfg),
		StateTransition: NewStateTransitionClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		StateMachine.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.StateMachine.Use(hooks...)
	c.StateTransition.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.StateMachine.Intercept(interceptors...)
	c.StateTransition.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *StateMachineMutation:
		return c.StateMachine.mutate(ctx, m)
	case *StateTransitionMutation:
		return c.StateTransition.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// StateMachineClient is a client for the StateMachine schema.
type StateMachineClient struct {
	config
}

// NewStateMachineClient returns a client for the StateMachine from the given config.
func NewStateMachineClient(c config) *StateMachineClient {
	return &StateMachineClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `statemachine.Hooks(f(g(h())))`.
func (c *StateMachineClient) Use(hooks ...Hook) {
	c.hooks.StateMachine = append(c.hooks.StateMachine, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `statemachine.Intercept(f(g(h())))`.
func (c *StateMachineClient) Intercept(interceptors ...Interceptor) {
	c.inters.StateMachine = append(c.inters.StateMachine, interceptors...)
}

// Create returns a builder for creating a StateMachine entity.
func (c *StateMachineClient) Create() *StateMachineCreate {
	mutation := newStateMachineMutation(c.config, OpCreate)
	return &StateMachineCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of StateMachine entities.
func (c *StateMachineClient) CreateBulk(builders ...*StateMachineCreate) *StateMachineCreateBulk {
	return &StateMachineCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *StateMachineClient) MapCreateBulk(slice any, setFunc func(*StateMachineCreate, int)) *StateMachineCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &StateMachineCreateBulk{err: fmt.Errorf("calling to StateMachineClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*StateMachineCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &StateMachineCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for StateMachine.
func (c *StateMachineClient) Update() *StateMachineUpdate {
	mutation := newStateMachineMutation(c.config, OpUpdate)
	return &StateMachineUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *StateMachineClient) UpdateOne(sm *StateMachine) *StateMachineUpdateOne {
	mutation := newStateMachineMutation(c.config, OpUpdateOne, withStateMachine(sm))
	return &StateMachineUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *StateMachineClient) UpdateOneID(id int) *StateMachineUpdateOne {
	mutation := newStateMachineMutation(c.config, OpUpdateOne, withStateMachineID(id))
	return &StateMachineUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for StateMachine.
func (c *StateMachineClient) Delete() *StateMachineDelete {
	mutation := newStateMachineMutation(c.config, OpDelete)
	return &StateMachineDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *StateMachineClient) DeleteOne(sm *StateMachine) *StateMachineDeleteOne {
	return c.DeleteOneID(sm.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *StateMachineClient) DeleteOneID(id int) *StateMachineDeleteOne {
	builder := c.Delete().Where(statemachine.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &StateMachineDeleteOne{builder}
}

// Query returns a query builder for StateMachine.
func (c *StateMachineClient) Query() *StateMachineQuery {
	return &StateMachineQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeStateMachine},
		inters: c.Interceptors(),
	}
}

// Get returns a StateMachine entity by its id.
func (c *StateMachineClient) Get(ctx context.Context, id int) (*StateMachine, error) {
	return c.Query().Where(statemachine.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *StateMachineClient) GetX(ctx context.Context, id int) *StateMachine {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryHistory queries the history edge of a StateMachine.
func (c *StateMachineClient) QueryHistory(sm *StateMachine) *StateTransitionQuery {
	query := (&StateTransitionClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := sm.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(statemachine.Table, statemachine.FieldID, id),
			sqlgraph.To(statetransition.Table, statetransition.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, statemachine.HistoryTable, statemachine.HistoryColumn),
		)
		fromV = sqlgraph.Neighbors(sm.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *StateMachineClient) Hooks() []Hook {
	return c.hooks.StateMachine
}

// Interceptors returns the client interceptors.
func (c *StateMachineClient) Interceptors() []Interceptor {
	return c.inters.StateMachine
}

func (c *StateMachineClient) mutate(ctx context.Context, m *StateMachineMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&StateMachineCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&StateMachineUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&StateMachineUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&StateMachineDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown StateMachine mutation op: %q", m.Op())
	}
}

// StateTransitionClient is a client for the StateTransition schema.
type StateTransitionClient struct {
	config
}

// NewStateTransitionClient returns a client for the StateTransition from the given config.
func NewStateTransitionClient(c config) *StateTransitionClient {
	return &StateTransitionClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `statetransition.Hooks(f(g(h())))`.
func (c *StateTransitionClient) Use(hooks ...Hook) {
	c.hooks.StateTransition = append(c.hooks.StateTransition, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `statetransition.Intercept(f(g(h())))`.
func (c *StateTransitionClient) Intercept(interceptors ...Interceptor) {
	c.inters.StateTransition = append(c.inters.StateTransition, interceptors...)
}

// Create returns a builder for creating a StateTransition entity.
func (c *StateTransitionClient) Create() *StateTransitionCreate {
	mutation := newStateTransitionMutation(c.config, OpCreate)
	return &StateTransitionCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of StateTransition entities.
func (c *StateTransitionClient) CreateBulk(builders ...*StateTransitionCreate) *StateTransitionCreateBulk {
	return &StateTransitionCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *StateTransitionClient) MapCreateBulk(slice any, setFunc func(*StateTransitionCreate, int)) *StateTransitionCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &StateTransitionCreateBulk{err: fmt.Errorf("calling to StateTransitionClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*StateTransitionCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &StateTransitionCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for StateTransition.
func (c *StateTransitionClient) Update() *StateTransitionUpdate {
	mutation := newStateTransitionMutation(c.config, OpUpdate)
	return &StateTransitionUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *StateTransitionClient) UpdateOne(st *StateTransition) *StateTransitionUpdateOne {
	mutation := newStateTransitionMutation(c.config, OpUpdateOne, withStateTransition(st))
	return &StateTransitionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *StateTransitionClient) UpdateOneID(id int) *StateTransitionUpdateOne {
	mutation := newStateTransitionMutation(c.config, OpUpdateOne, withStateTransitionID(id))
	return &StateTransitionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for StateTransition.
func (c *StateTransitionClient) Delete() *StateTransitionDelete {
	mutation := newStateTransitionMutation(c.config, OpDelete)
	return &StateTransitionDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *StateTransitionClient) DeleteOne(st *StateTransition) *StateTransitionDeleteOne {
	return c.DeleteOneID(st.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *StateTransitionClient) DeleteOneID(id int) *StateTransitionDeleteOne {
	builder := c.Delete().Where(statetransition.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &StateTransitionDeleteOne{builder}
}

// Query returns a query builder for StateTransition.
func (c *StateTransitionClient) Query() *StateTransitionQuery {
	return &StateTransitionQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeStateTransition},
		inters: c.Interceptors(),
	}
}

// Get returns a StateTransition entity by its id.
func (c *StateTransitionClient) Get(ctx context.Context, id int) (*StateTransition, error) {
	return c.Query().Where(statetransition.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *StateTransitionClient) GetX(ctx context.Context, id int) *StateTransition {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryMachine queries the machine edge of a StateTransition.
func (c *StateTransitionClient) QueryMachine(st *StateTransition) *StateMachineQuery {
	query := (&StateMachineClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := st.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(statetransition.Table, statetransition.FieldID, id),
			sqlgraph.To(statemachine.Table, statemachine.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, statetransition.MachineTable, statetransition.MachineColumn),
		)
		fromV = sqlgraph.Neighbors(st.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *StateTransitionClient) Hooks() []Hook {
	return c.hooks.StateTransition
}

// Interceptors returns the client interceptors.
func (c *StateTransitionClient) Interceptors() []Interceptor {
	return c.inters.StateTransition
}

func (c *StateTransitionClient) mutate(ctx context.Context, m *StateTransitionMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&StateTransitionCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&StateTransitionUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&StateTransitionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&StateTransitionDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown StateTransition mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		StateMachine, StateTransition []ent.Hook
	}
	inters struct {
		StateMachine, StateTransition []ent.Interceptor
	}
)
