// Code generated by ent, DO NOT EDIT.

package channel

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the channel type in the database.
	Label = "channel"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldChannelID holds the string denoting the channel_id field in the database.
	FieldChannelID = "channel_id"
	// FieldService holds the string denoting the service field in the database.
	FieldService = "service"
	// FieldIsLive holds the string denoting the is_live field in the database.
	FieldIsLive = "is_live"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldCategory holds the string denoting the category field in the database.
	FieldCategory = "category"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// EdgeFollows holds the string denoting the follows edge name in mutations.
	EdgeFollows = "follows"
	// EdgeStreams holds the string denoting the streams edge name in mutations.
	EdgeStreams = "streams"
	// Table holds the table name of the channel in the database.
	Table = "channels"
	// FollowsTable is the table that holds the follows relation/edge.
	FollowsTable = "follows"
	// FollowsInverseTable is the table name for the Follow entity.
	// It exists in this package in order to avoid circular dependency with the "follow" package.
	FollowsInverseTable = "follows"
	// FollowsColumn is the table column denoting the follows relation/edge.
	FollowsColumn = "channel_id"
	// StreamsTable is the table that holds the streams relation/edge.
	StreamsTable = "streams"
	// StreamsInverseTable is the table name for the Stream entity.
	// It exists in this package in order to avoid circular dependency with the "stream" package.
	StreamsInverseTable = "streams"
	// StreamsColumn is the table column denoting the streams relation/edge.
	StreamsColumn = "channel_id"
)

// Columns holds all SQL columns for channel fields.
var Columns = []string{
	FieldID,
	FieldChannelID,
	FieldService,
	FieldIsLive,
	FieldTitle,
	FieldCategory,
	FieldUpdatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultIsLive holds the default value on creation for the "is_live" field.
	DefaultIsLive bool
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// Service defines the type for the "service" enum field.
type Service string

// Service values.
const (
	ServiceTwitch Service = "twitch"
)

func (s Service) String() string {
	return string(s)
}

// ServiceValidator is a validator for the "service" field enum values. It is called by the builders before save.
func ServiceValidator(s Service) error {
	switch s {
	case ServiceTwitch:
		return nil
	default:
		return fmt.Errorf("channel: invalid enum value for service field: %q", s)
	}
}