-- Modify "chat_settings" table
ALTER TABLE "chat_settings" ADD COLUMN "image_in_notification" boolean NOT NULL DEFAULT true;
