// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/satont/twitch-notifier/ent/chat"
	"github.com/satont/twitch-notifier/ent/chatsettings"
	"github.com/satont/twitch-notifier/ent/predicate"
)

// ChatSettingsUpdate is the builder for updating ChatSettings entities.
type ChatSettingsUpdate struct {
	config
	hooks    []Hook
	mutation *ChatSettingsMutation
}

// Where appends a list predicates to the ChatSettingsUpdate builder.
func (csu *ChatSettingsUpdate) Where(ps ...predicate.ChatSettings) *ChatSettingsUpdate {
	csu.mutation.Where(ps...)
	return csu
}

// SetGameChangeNotification sets the "game_change_notification" field.
func (csu *ChatSettingsUpdate) SetGameChangeNotification(b bool) *ChatSettingsUpdate {
	csu.mutation.SetGameChangeNotification(b)
	return csu
}

// SetNillableGameChangeNotification sets the "game_change_notification" field if the given value is not nil.
func (csu *ChatSettingsUpdate) SetNillableGameChangeNotification(b *bool) *ChatSettingsUpdate {
	if b != nil {
		csu.SetGameChangeNotification(*b)
	}
	return csu
}

// SetOfflineNotification sets the "offline_notification" field.
func (csu *ChatSettingsUpdate) SetOfflineNotification(b bool) *ChatSettingsUpdate {
	csu.mutation.SetOfflineNotification(b)
	return csu
}

// SetNillableOfflineNotification sets the "offline_notification" field if the given value is not nil.
func (csu *ChatSettingsUpdate) SetNillableOfflineNotification(b *bool) *ChatSettingsUpdate {
	if b != nil {
		csu.SetOfflineNotification(*b)
	}
	return csu
}

// SetChatLanguage sets the "chat_language" field.
func (csu *ChatSettingsUpdate) SetChatLanguage(cl chatsettings.ChatLanguage) *ChatSettingsUpdate {
	csu.mutation.SetChatLanguage(cl)
	return csu
}

// SetNillableChatLanguage sets the "chat_language" field if the given value is not nil.
func (csu *ChatSettingsUpdate) SetNillableChatLanguage(cl *chatsettings.ChatLanguage) *ChatSettingsUpdate {
	if cl != nil {
		csu.SetChatLanguage(*cl)
	}
	return csu
}

// SetChatID sets the "chat_id" field.
func (csu *ChatSettingsUpdate) SetChatID(u uuid.UUID) *ChatSettingsUpdate {
	csu.mutation.SetChatID(u)
	return csu
}

// SetChat sets the "chat" edge to the Chat entity.
func (csu *ChatSettingsUpdate) SetChat(c *Chat) *ChatSettingsUpdate {
	return csu.SetChatID(c.ID)
}

// Mutation returns the ChatSettingsMutation object of the builder.
func (csu *ChatSettingsUpdate) Mutation() *ChatSettingsMutation {
	return csu.mutation
}

// ClearChat clears the "chat" edge to the Chat entity.
func (csu *ChatSettingsUpdate) ClearChat() *ChatSettingsUpdate {
	csu.mutation.ClearChat()
	return csu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (csu *ChatSettingsUpdate) Save(ctx context.Context) (int, error) {
	return withHooks[int, ChatSettingsMutation](ctx, csu.sqlSave, csu.mutation, csu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (csu *ChatSettingsUpdate) SaveX(ctx context.Context) int {
	affected, err := csu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (csu *ChatSettingsUpdate) Exec(ctx context.Context) error {
	_, err := csu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (csu *ChatSettingsUpdate) ExecX(ctx context.Context) {
	if err := csu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (csu *ChatSettingsUpdate) check() error {
	if v, ok := csu.mutation.ChatLanguage(); ok {
		if err := chatsettings.ChatLanguageValidator(v); err != nil {
			return &ValidationError{Name: "chat_language", err: fmt.Errorf(`ent: validator failed for field "ChatSettings.chat_language": %w`, err)}
		}
	}
	if _, ok := csu.mutation.ChatID(); csu.mutation.ChatCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "ChatSettings.chat"`)
	}
	return nil
}

func (csu *ChatSettingsUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := csu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(chatsettings.Table, chatsettings.Columns, sqlgraph.NewFieldSpec(chatsettings.FieldID, field.TypeUUID))
	if ps := csu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := csu.mutation.GameChangeNotification(); ok {
		_spec.SetField(chatsettings.FieldGameChangeNotification, field.TypeBool, value)
	}
	if value, ok := csu.mutation.OfflineNotification(); ok {
		_spec.SetField(chatsettings.FieldOfflineNotification, field.TypeBool, value)
	}
	if value, ok := csu.mutation.ChatLanguage(); ok {
		_spec.SetField(chatsettings.FieldChatLanguage, field.TypeEnum, value)
	}
	if csu.mutation.ChatCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   chatsettings.ChatTable,
			Columns: []string{chatsettings.ChatColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chat.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := csu.mutation.ChatIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   chatsettings.ChatTable,
			Columns: []string{chatsettings.ChatColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chat.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, csu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{chatsettings.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	csu.mutation.done = true
	return n, nil
}

// ChatSettingsUpdateOne is the builder for updating a single ChatSettings entity.
type ChatSettingsUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ChatSettingsMutation
}

// SetGameChangeNotification sets the "game_change_notification" field.
func (csuo *ChatSettingsUpdateOne) SetGameChangeNotification(b bool) *ChatSettingsUpdateOne {
	csuo.mutation.SetGameChangeNotification(b)
	return csuo
}

// SetNillableGameChangeNotification sets the "game_change_notification" field if the given value is not nil.
func (csuo *ChatSettingsUpdateOne) SetNillableGameChangeNotification(b *bool) *ChatSettingsUpdateOne {
	if b != nil {
		csuo.SetGameChangeNotification(*b)
	}
	return csuo
}

// SetOfflineNotification sets the "offline_notification" field.
func (csuo *ChatSettingsUpdateOne) SetOfflineNotification(b bool) *ChatSettingsUpdateOne {
	csuo.mutation.SetOfflineNotification(b)
	return csuo
}

// SetNillableOfflineNotification sets the "offline_notification" field if the given value is not nil.
func (csuo *ChatSettingsUpdateOne) SetNillableOfflineNotification(b *bool) *ChatSettingsUpdateOne {
	if b != nil {
		csuo.SetOfflineNotification(*b)
	}
	return csuo
}

// SetChatLanguage sets the "chat_language" field.
func (csuo *ChatSettingsUpdateOne) SetChatLanguage(cl chatsettings.ChatLanguage) *ChatSettingsUpdateOne {
	csuo.mutation.SetChatLanguage(cl)
	return csuo
}

// SetNillableChatLanguage sets the "chat_language" field if the given value is not nil.
func (csuo *ChatSettingsUpdateOne) SetNillableChatLanguage(cl *chatsettings.ChatLanguage) *ChatSettingsUpdateOne {
	if cl != nil {
		csuo.SetChatLanguage(*cl)
	}
	return csuo
}

// SetChatID sets the "chat_id" field.
func (csuo *ChatSettingsUpdateOne) SetChatID(u uuid.UUID) *ChatSettingsUpdateOne {
	csuo.mutation.SetChatID(u)
	return csuo
}

// SetChat sets the "chat" edge to the Chat entity.
func (csuo *ChatSettingsUpdateOne) SetChat(c *Chat) *ChatSettingsUpdateOne {
	return csuo.SetChatID(c.ID)
}

// Mutation returns the ChatSettingsMutation object of the builder.
func (csuo *ChatSettingsUpdateOne) Mutation() *ChatSettingsMutation {
	return csuo.mutation
}

// ClearChat clears the "chat" edge to the Chat entity.
func (csuo *ChatSettingsUpdateOne) ClearChat() *ChatSettingsUpdateOne {
	csuo.mutation.ClearChat()
	return csuo
}

// Where appends a list predicates to the ChatSettingsUpdate builder.
func (csuo *ChatSettingsUpdateOne) Where(ps ...predicate.ChatSettings) *ChatSettingsUpdateOne {
	csuo.mutation.Where(ps...)
	return csuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (csuo *ChatSettingsUpdateOne) Select(field string, fields ...string) *ChatSettingsUpdateOne {
	csuo.fields = append([]string{field}, fields...)
	return csuo
}

// Save executes the query and returns the updated ChatSettings entity.
func (csuo *ChatSettingsUpdateOne) Save(ctx context.Context) (*ChatSettings, error) {
	return withHooks[*ChatSettings, ChatSettingsMutation](ctx, csuo.sqlSave, csuo.mutation, csuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (csuo *ChatSettingsUpdateOne) SaveX(ctx context.Context) *ChatSettings {
	node, err := csuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (csuo *ChatSettingsUpdateOne) Exec(ctx context.Context) error {
	_, err := csuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (csuo *ChatSettingsUpdateOne) ExecX(ctx context.Context) {
	if err := csuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (csuo *ChatSettingsUpdateOne) check() error {
	if v, ok := csuo.mutation.ChatLanguage(); ok {
		if err := chatsettings.ChatLanguageValidator(v); err != nil {
			return &ValidationError{Name: "chat_language", err: fmt.Errorf(`ent: validator failed for field "ChatSettings.chat_language": %w`, err)}
		}
	}
	if _, ok := csuo.mutation.ChatID(); csuo.mutation.ChatCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "ChatSettings.chat"`)
	}
	return nil
}

func (csuo *ChatSettingsUpdateOne) sqlSave(ctx context.Context) (_node *ChatSettings, err error) {
	if err := csuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(chatsettings.Table, chatsettings.Columns, sqlgraph.NewFieldSpec(chatsettings.FieldID, field.TypeUUID))
	id, ok := csuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "ChatSettings.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := csuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, chatsettings.FieldID)
		for _, f := range fields {
			if !chatsettings.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != chatsettings.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := csuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := csuo.mutation.GameChangeNotification(); ok {
		_spec.SetField(chatsettings.FieldGameChangeNotification, field.TypeBool, value)
	}
	if value, ok := csuo.mutation.OfflineNotification(); ok {
		_spec.SetField(chatsettings.FieldOfflineNotification, field.TypeBool, value)
	}
	if value, ok := csuo.mutation.ChatLanguage(); ok {
		_spec.SetField(chatsettings.FieldChatLanguage, field.TypeEnum, value)
	}
	if csuo.mutation.ChatCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   chatsettings.ChatTable,
			Columns: []string{chatsettings.ChatColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chat.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := csuo.mutation.ChatIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   chatsettings.ChatTable,
			Columns: []string{chatsettings.ChatColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chat.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &ChatSettings{config: csuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, csuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{chatsettings.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	csuo.mutation.done = true
	return _node, nil
}
