// Code generated by ent, DO NOT EDIT.

package statetransition

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the statetransition type in the database.
	Label = "state_transition"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldFromState holds the string denoting the from_state field in the database.
	FieldFromState = "from_state"
	// FieldToState holds the string denoting the to_state field in the database.
	FieldToState = "to_state"
	// FieldEvent holds the string denoting the event field in the database.
	FieldEvent = "event"
	// FieldTimestamp holds the string denoting the timestamp field in the database.
	FieldTimestamp = "timestamp"
	// EdgeMachine holds the string denoting the machine edge name in mutations.
	EdgeMachine = "machine"
	// Table holds the table name of the statetransition in the database.
	Table = "state_transitions"
	// MachineTable is the table that holds the machine relation/edge.
	MachineTable = "state_transitions"
	// MachineInverseTable is the table name for the StateMachine entity.
	// It exists in this package in order to avoid circular dependency with the "statemachine" package.
	MachineInverseTable = "state_machines"
	// MachineColumn is the table column denoting the machine relation/edge.
	MachineColumn = "state_machine_history"
)

// Columns holds all SQL columns for statetransition fields.
var Columns = []string{
	FieldID,
	FieldFromState,
	FieldToState,
	FieldEvent,
	FieldTimestamp,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "state_transitions"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"state_machine_history",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultTimestamp holds the default value on creation for the "timestamp" field.
	DefaultTimestamp func() time.Time
)

// OrderOption defines the ordering options for the StateTransition queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByFromState orders the results by the from_state field.
func ByFromState(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFromState, opts...).ToFunc()
}

// ByToState orders the results by the to_state field.
func ByToState(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldToState, opts...).ToFunc()
}

// ByEvent orders the results by the event field.
func ByEvent(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEvent, opts...).ToFunc()
}

// ByTimestamp orders the results by the timestamp field.
func ByTimestamp(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTimestamp, opts...).ToFunc()
}

// ByMachineField orders the results by machine field.
func ByMachineField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMachineStep(), sql.OrderByField(field, opts...))
	}
}
func newMachineStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MachineInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, MachineTable, MachineColumn),
	)
}
