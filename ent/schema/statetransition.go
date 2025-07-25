package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// StateTransition holds the schema definition for the StateTransition entity.
type StateTransition struct {
	ent.Schema
}

// Fields of the StateTransition.
func (StateTransition) Fields() []ent.Field {
	return []ent.Field{
		field.String("from_state"),
		field.String("to_state"),
		field.String("event"),
		field.Time("timestamp").
			Default(time.Now),
	}
}

// Edges of the StateTransition.
func (StateTransition) Edges() []ent.Edge {
	return []ent.Edge{
		// Create an inverse-edge to the StateMachine entity.
		// This creates a "history" edge on the StateMachine entity.
		edge.From("machine", StateMachine.Type).
			Ref("history").
			Unique(),
	}
}
