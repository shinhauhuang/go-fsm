package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// StateMachine holds the schema definition for the StateMachine entity.
type StateMachine struct {
	ent.Schema
}

// Fields of the StateMachine.
func (StateMachine) Fields() []ent.Field {
	return []ent.Field{
		field.String("machine_id").
			Unique().
			NotEmpty(),
		field.String("current_state").
			NotEmpty(),
	}
}

// Edges of the StateMachine.
func (StateMachine) Edges() []ent.Edge {
	return []ent.Edge{
		// Create a one-to-many relationship with StateTransition.
		// This means a StateMachine can have many history records.
		edge.To("history", StateTransition.Type),
	}
}
