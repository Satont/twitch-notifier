-- create "chats" table
CREATE TABLE "chats" ("id" uuid NOT NULL, "chat_id" character varying NOT NULL, "service" character varying NOT NULL, PRIMARY KEY ("id"));
-- create index "chat_chat_id_service" to table: "chats"
CREATE UNIQUE INDEX "chat_chat_id_service" ON "chats" ("chat_id", "service");
-- create "chat_settings" table
CREATE TABLE "chat_settings" ("id" uuid NOT NULL, "game_change_notification" boolean NOT NULL DEFAULT true, "offline_notification" boolean NOT NULL DEFAULT true, "chat_language" character varying NOT NULL DEFAULT 'en', "chat_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "chat_settings_chats_settings" FOREIGN KEY ("chat_id") REFERENCES "chats" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- create index "chat_settings_chat_id_key" to table: "chat_settings"
CREATE UNIQUE INDEX "chat_settings_chat_id_key" ON "chat_settings" ("chat_id");
-- create "channels" table
CREATE TABLE "channels" ("id" uuid NOT NULL, "channel_id" character varying NOT NULL, "service" character varying NOT NULL, "is_live" boolean NOT NULL DEFAULT false, "title" character varying NULL, "category" character varying NULL, "updated_at" timestamptz NULL, PRIMARY KEY ("id"));
-- create index "channel_channel_id_service" to table: "channels"
CREATE UNIQUE INDEX "channel_channel_id_service" ON "channels" ("channel_id", "service");
-- create "follows" table
CREATE TABLE "follows" ("id" uuid NOT NULL, "channel_id" uuid NOT NULL, "chat_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "follows_channels_follows" FOREIGN KEY ("channel_id") REFERENCES "channels" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "follows_chats_follows" FOREIGN KEY ("chat_id") REFERENCES "chats" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- create "streams" table
CREATE TABLE "streams" ("id" character varying NOT NULL, "titles" text[] NULL, "categories" text[] NULL, "started_at" timestamptz NULL, "updated_at" timestamptz NULL, "ended_at" timestamptz NULL, "channel_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "streams_channels_streams" FOREIGN KEY ("channel_id") REFERENCES "channels" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
