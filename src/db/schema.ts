import { sqliteTable, text, integer } from 'drizzle-orm/sqlite-core';
import { relations } from 'drizzle-orm';
import { randomUUID } from 'node:crypto';

// Chat table
export const chats = sqliteTable('chats', {
  id: text('id').primaryKey().$defaultFn(() => randomUUID()),
  chatId: text('chat_id').notNull(),
  service: text('service', { enum: ['telegram'] }).notNull().default('telegram'),
});

export const chatsRelations = relations(chats, ({ one, many }) => ({
  settings: one(chatSettings, {
    fields: [chats.id],
    references: [chatSettings.chatId],
  }),
  follows: many(follows),
}));

// Chat Settings table
export const chatSettings = sqliteTable('chat_settings', {
  id: text('id').primaryKey().$defaultFn(() => randomUUID()),
  chatId: text('chat_id').notNull().unique().references(() => chats.id, { onDelete: 'cascade' }),
  gameChangeNotification: integer('game_change_notification', { mode: 'boolean' }).notNull().default(true),
  titleChangeNotification: integer('title_change_notification', { mode: 'boolean' }).notNull().default(false),
  gameAndTitleChangeNotification: integer('game_and_title_change_notification', { mode: 'boolean' }).notNull().default(false),
  offlineNotification: integer('offline_notification', { mode: 'boolean' }).notNull().default(true),
  imageInNotification: integer('image_in_notification', { mode: 'boolean' }).notNull().default(true),
  language: text('language', { enum: ['ru', 'en', 'uk'] }).notNull().default('en'),
});

export const chatSettingsRelations = relations(chatSettings, ({ one }) => ({
  chat: one(chats, {
    fields: [chatSettings.chatId],
    references: [chats.id],
  }),
}));

// Channel table
export const channels = sqliteTable('channels', {
  id: text('id').primaryKey().$defaultFn(() => randomUUID()),
  channelId: text('channel_id').notNull(),
  service: text('service', { enum: ['twitch'] }).notNull().default('twitch'),
  isLive: integer('is_live', { mode: 'boolean' }).notNull().default(false),
  title: text('title'),
  category: text('category'),
  updatedAt: text('updated_at').$defaultFn(() => new Date().toISOString()),
});

export const channelsRelations = relations(channels, ({ many }) => ({
  follows: many(follows),
  streams: many(streams),
}));

// Follow table
export const follows = sqliteTable('follows', {
  id: text('id').primaryKey().$defaultFn(() => randomUUID()),
  channelId: text('channel_id').notNull().references(() => channels.id, { onDelete: 'cascade' }),
  chatId: text('chat_id').notNull().references(() => chats.id, { onDelete: 'cascade' }),
});

export const followsRelations = relations(follows, ({ one }) => ({
  channel: one(channels, {
    fields: [follows.channelId],
    references: [channels.id],
  }),
  chat: one(chats, {
    fields: [follows.chatId],
    references: [chats.id],
  }),
}));

// Stream table
export const streams = sqliteTable('streams', {
  id: text('id').primaryKey(), // Twitch stream ID
  channelId: text('channel_id').notNull().references(() => channels.id, { onDelete: 'cascade' }),
  isLive: integer('is_live', { mode: 'boolean' }).notNull().default(true),
  title: text('title'),
  category: text('category'),
  titles: text('titles', { mode: 'json' }).$type<string[]>().notNull().default([]),
  categories: text('categories', { mode: 'json' }).$type<string[]>().notNull().default([]),
  startedAt: text('started_at').$defaultFn(() => new Date().toISOString()),
  updatedAt: text('updated_at').$defaultFn(() => new Date().toISOString()),
  endedAt: text('ended_at'),
});

export const streamsRelations = relations(streams, ({ one }) => ({
  channel: one(channels, {
    fields: [streams.channelId],
    references: [channels.id],
  }),
}));

// Types for insert and select
export type Chat = typeof chats.$inferSelect;
export type NewChat = typeof chats.$inferInsert;

export type ChatSettings = typeof chatSettings.$inferSelect;
export type NewChatSettings = typeof chatSettings.$inferInsert;

export type Channel = typeof channels.$inferSelect;
export type NewChannel = typeof channels.$inferInsert;

export type Follow = typeof follows.$inferSelect;
export type NewFollow = typeof follows.$inferInsert;

export type Stream = typeof streams.$inferSelect;
export type NewStream = typeof streams.$inferInsert;
