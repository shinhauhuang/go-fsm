// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shinhauhuang/go-fsm/ent/predicate"
	"github.com/shinhauhuang/go-fsm/ent/statemachine"
	"github.com/shinhauhuang/go-fsm/ent/statetransition"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// StateTransitionUpdate is the builder for updating StateTransition entities.
type StateTransitionUpdate struct {
	config
	hooks    []Hook
	mutation *StateTransitionMutation
}

// Where appends a list predicates to the StateTransitionUpdate builder.
func (stu *StateTransitionUpdate) Where(ps ...predicate.StateTransition) *StateTransitionUpdate {
	stu.mutation.Where(ps...)
	return stu
}

// SetFromState sets the "from_state" field.
func (stu *StateTransitionUpdate) SetFromState(s string) *StateTransitionUpdate {
	stu.mutation.SetFromState(s)
	return stu
}

// SetNillableFromState sets the "from_state" field if the given value is not nil.
func (stu *StateTransitionUpdate) SetNillableFromState(s *string) *StateTransitionUpdate {
	if s != nil {
		stu.SetFromState(*s)
	}
	return stu
}

// SetToState sets the "to_state" field.
func (stu *StateTransitionUpdate) SetToState(s string) *StateTransitionUpdate {
	stu.mutation.SetToState(s)
	return stu
}

// SetNillableToState sets the "to_state" field if the given value is not nil.
func (stu *StateTransitionUpdate) SetNillableToState(s *string) *StateTransitionUpdate {
	if s != nil {
		stu.SetToState(*s)
	}
	return stu
}

// SetEvent sets the "event" field.
func (stu *StateTransitionUpdate) SetEvent(s string) *StateTransitionUpdate {
	stu.mutation.SetEvent(s)
	return stu
}

// SetNillableEvent sets the "event" field if the given value is not nil.
func (stu *StateTransitionUpdate) SetNillableEvent(s *string) *StateTransitionUpdate {
	if s != nil {
		stu.SetEvent(*s)
	}
	return stu
}

// SetTimestamp sets the "timestamp" field.
func (stu *StateTransitionUpdate) SetTimestamp(t time.Time) *StateTransitionUpdate {
	stu.mutation.SetTimestamp(t)
	return stu
}

// SetNillableTimestamp sets the "timestamp" field if the given value is not nil.
func (stu *StateTransitionUpdate) SetNillableTimestamp(t *time.Time) *StateTransitionUpdate {
	if t != nil {
		stu.SetTimestamp(*t)
	}
	return stu
}

// SetMachineID sets the "machine" edge to the StateMachine entity by ID.
func (stu *StateTransitionUpdate) SetMachineID(id int) *StateTransitionUpdate {
	stu.mutation.SetMachineID(id)
	return stu
}

// SetNillableMachineID sets the "machine" edge to the StateMachine entity by ID if the given value is not nil.
func (stu *StateTransitionUpdate) SetNillableMachineID(id *int) *StateTransitionUpdate {
	if id != nil {
		stu = stu.SetMachineID(*id)
	}
	return stu
}

// SetMachine sets the "machine" edge to the StateMachine entity.
func (stu *StateTransitionUpdate) SetMachine(s *StateMachine) *StateTransitionUpdate {
	return stu.SetMachineID(s.ID)
}

// Mutation returns the StateTransitionMutation object of the builder.
func (stu *StateTransitionUpdate) Mutation() *StateTransitionMutation {
	return stu.mutation
}

// ClearMachine clears the "machine" edge to the StateMachine entity.
func (stu *StateTransitionUpdate) ClearMachine() *StateTransitionUpdate {
	stu.mutation.ClearMachine()
	return stu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (stu *StateTransitionUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, stu.sqlSave, stu.mutation, stu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (stu *StateTransitionUpdate) SaveX(ctx context.Context) int {
	affected, err := stu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (stu *StateTransitionUpdate) Exec(ctx context.Context) error {
	_, err := stu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (stu *StateTransitionUpdate) ExecX(ctx context.Context) {
	if err := stu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (stu *StateTransitionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(statetransition.Table, statetransition.Columns, sqlgraph.NewFieldSpec(statetransition.FieldID, field.TypeInt))
	if ps := stu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := stu.mutation.FromState(); ok {
		_spec.SetField(statetransition.FieldFromState, field.TypeString, value)
	}
	if value, ok := stu.mutation.ToState(); ok {
		_spec.SetField(statetransition.FieldToState, field.TypeString, value)
	}
	if value, ok := stu.mutation.Event(); ok {
		_spec.SetField(statetransition.FieldEvent, field.TypeString, value)
	}
	if value, ok := stu.mutation.Timestamp(); ok {
		_spec.SetField(statetransition.FieldTimestamp, field.TypeTime, value)
	}
	if stu.mutation.MachineCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   statetransition.MachineTable,
			Columns: []string{statetransition.MachineColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(statemachine.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := stu.mutation.MachineIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   statetransition.MachineTable,
			Columns: []string{statetransition.MachineColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(statemachine.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, stu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{statetransition.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	stu.mutation.done = true
	return n, nil
}

// StateTransitionUpdateOne is the builder for updating a single StateTransition entity.
type StateTransitionUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *StateTransitionMutation
}

// SetFromState sets the "from_state" field.
func (stuo *StateTransitionUpdateOne) SetFromState(s string) *StateTransitionUpdateOne {
	stuo.mutation.SetFromState(s)
	return stuo
}

// SetNillableFromState sets the "from_state" field if the given value is not nil.
func (stuo *StateTransitionUpdateOne) SetNillableFromState(s *string) *StateTransitionUpdateOne {
	if s != nil {
		stuo.SetFromState(*s)
	}
	return stuo
}

// SetToState sets the "to_state" field.
func (stuo *StateTransitionUpdateOne) SetToState(s string) *StateTransitionUpdateOne {
	stuo.mutation.SetToState(s)
	return stuo
}

// SetNillableToState sets the "to_state" field if the given value is not nil.
func (stuo *StateTransitionUpdateOne) SetNillableToState(s *string) *StateTransitionUpdateOne {
	if s != nil {
		stuo.SetToState(*s)
	}
	return stuo
}

// SetEvent sets the "event" field.
func (stuo *StateTransitionUpdateOne) SetEvent(s string) *StateTransitionUpdateOne {
	stuo.mutation.SetEvent(s)
	return stuo
}

// SetNillableEvent sets the "event" field if the given value is not nil.
func (stuo *StateTransitionUpdateOne) SetNillableEvent(s *string) *StateTransitionUpdateOne {
	if s != nil {
		stuo.SetEvent(*s)
	}
	return stuo
}

// SetTimestamp sets the "timestamp" field.
func (stuo *StateTransitionUpdateOne) SetTimestamp(t time.Time) *StateTransitionUpdateOne {
	stuo.mutation.SetTimestamp(t)
	return stuo
}

// SetNillableTimestamp sets the "timestamp" field if the given value is not nil.
func (stuo *StateTransitionUpdateOne) SetNillableTimestamp(t *time.Time) *StateTransitionUpdateOne {
	if t != nil {
		stuo.SetTimestamp(*t)
	}
	return stuo
}

// SetMachineID sets the "machine" edge to the StateMachine entity by ID.
func (stuo *StateTransitionUpdateOne) SetMachineID(id int) *StateTransitionUpdateOne {
	stuo.mutation.SetMachineID(id)
	return stuo
}

// SetNillableMachineID sets the "machine" edge to the StateMachine entity by ID if the given value is not nil.
func (stuo *StateTransitionUpdateOne) SetNillableMachineID(id *int) *StateTransitionUpdateOne {
	if id != nil {
		stuo = stuo.SetMachineID(*id)
	}
	return stuo
}

// SetMachine sets the "machine" edge to the StateMachine entity.
func (stuo *StateTransitionUpdateOne) SetMachine(s *StateMachine) *StateTransitionUpdateOne {
	return stuo.SetMachineID(s.ID)
}

// Mutation returns the StateTransitionMutation object of the builder.
func (stuo *StateTransitionUpdateOne) Mutation() *StateTransitionMutation {
	return stuo.mutation
}

// ClearMachine clears the "machine" edge to the StateMachine entity.
func (stuo *StateTransitionUpdateOne) ClearMachine() *StateTransitionUpdateOne {
	stuo.mutation.ClearMachine()
	return stuo
}

// Where appends a list predicates to the StateTransitionUpdate builder.
func (stuo *StateTransitionUpdateOne) Where(ps ...predicate.StateTransition) *StateTransitionUpdateOne {
	stuo.mutation.Where(ps...)
	return stuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (stuo *StateTransitionUpdateOne) Select(field string, fields ...string) *StateTransitionUpdateOne {
	stuo.fields = append([]string{field}, fields...)
	return stuo
}

// Save executes the query and returns the updated StateTransition entity.
func (stuo *StateTransitionUpdateOne) Save(ctx context.Context) (*StateTransition, error) {
	return withHooks(ctx, stuo.sqlSave, stuo.mutation, stuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (stuo *StateTransitionUpdateOne) SaveX(ctx context.Context) *StateTransition {
	node, err := stuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (stuo *StateTransitionUpdateOne) Exec(ctx context.Context) error {
	_, err := stuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (stuo *StateTransitionUpdateOne) ExecX(ctx context.Context) {
	if err := stuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (stuo *StateTransitionUpdateOne) sqlSave(ctx context.Context) (_node *StateTransition, err error) {
	_spec := sqlgraph.NewUpdateSpec(statetransition.Table, statetransition.Columns, sqlgraph.NewFieldSpec(statetransition.FieldID, field.TypeInt))
	id, ok := stuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "StateTransition.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := stuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, statetransition.FieldID)
		for _, f := range fields {
			if !statetransition.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != statetransition.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := stuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := stuo.mutation.FromState(); ok {
		_spec.SetField(statetransition.FieldFromState, field.TypeString, value)
	}
	if value, ok := stuo.mutation.ToState(); ok {
		_spec.SetField(statetransition.FieldToState, field.TypeString, value)
	}
	if value, ok := stuo.mutation.Event(); ok {
		_spec.SetField(statetransition.FieldEvent, field.TypeString, value)
	}
	if value, ok := stuo.mutation.Timestamp(); ok {
		_spec.SetField(statetransition.FieldTimestamp, field.TypeTime, value)
	}
	if stuo.mutation.MachineCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   statetransition.MachineTable,
			Columns: []string{statetransition.MachineColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(statemachine.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := stuo.mutation.MachineIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   statetransition.MachineTable,
			Columns: []string{statetransition.MachineColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(statemachine.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &StateTransition{config: stuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, stuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{statetransition.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	stuo.mutation.done = true
	return _node, nil
}
