package schema

import (
	"entgo.io/ent"
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
	return nil
}
