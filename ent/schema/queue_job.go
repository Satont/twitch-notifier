package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/internal/queue"
)

type QueueJob struct {
	ent.Schema
}

func (QueueJob) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("queue_name"),
		field.Bytes("data"),
		field.Int("retries"),
		field.Int("max_retries"),
		field.Time("added_at"),
		field.Int("ttl"),
		field.Enum("status").GoType(queue.JobStatus("")).Default(queue.JobStatusPending.String()),
		field.String("fail_reason").Optional(),
	}
}

func (QueueJob) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("queue_name"),
	}
}
