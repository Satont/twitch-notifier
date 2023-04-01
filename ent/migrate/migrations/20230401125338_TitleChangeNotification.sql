-- Modify "chat_settings" table
ALTER TABLE "chat_settings" ADD COLUMN "title_change_notification" boolean NOT NULL DEFAULT false;
