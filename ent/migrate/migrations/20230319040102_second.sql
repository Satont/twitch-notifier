-- create index "follow_channel_id_chat_id" to table: "follows"
CREATE UNIQUE INDEX "follow_channel_id_chat_id" ON "follows" ("channel_id", "chat_id");
