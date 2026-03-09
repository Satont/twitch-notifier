CREATE TABLE `channels` (
	`id` text PRIMARY KEY NOT NULL,
	`channel_id` text NOT NULL,
	`service` text DEFAULT 'twitch' NOT NULL,
	`is_live` integer DEFAULT false NOT NULL,
	`title` text,
	`category` text,
	`updated_at` text
);
--> statement-breakpoint
CREATE TABLE `chat_settings` (
	`id` text PRIMARY KEY NOT NULL,
	`chat_id` text NOT NULL,
	`game_change_notification` integer DEFAULT true NOT NULL,
	`title_change_notification` integer DEFAULT false NOT NULL,
	`game_and_title_change_notification` integer DEFAULT false NOT NULL,
	`offline_notification` integer DEFAULT true NOT NULL,
	`image_in_notification` integer DEFAULT true NOT NULL,
	`language` text DEFAULT 'en' NOT NULL,
	FOREIGN KEY (`chat_id`) REFERENCES `chats`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE UNIQUE INDEX `chat_settings_chat_id_unique` ON `chat_settings` (`chat_id`);--> statement-breakpoint
CREATE TABLE `chats` (
	`id` text PRIMARY KEY NOT NULL,
	`chat_id` text NOT NULL,
	`service` text DEFAULT 'telegram' NOT NULL
);
--> statement-breakpoint
CREATE TABLE `follows` (
	`id` text PRIMARY KEY NOT NULL,
	`channel_id` text NOT NULL,
	`chat_id` text NOT NULL,
	FOREIGN KEY (`channel_id`) REFERENCES `channels`(`id`) ON UPDATE no action ON DELETE cascade,
	FOREIGN KEY (`chat_id`) REFERENCES `chats`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE TABLE `streams` (
	`id` text PRIMARY KEY NOT NULL,
	`channel_id` text NOT NULL,
	`is_live` integer DEFAULT true NOT NULL,
	`title` text,
	`category` text,
	`titles` text DEFAULT '[]' NOT NULL,
	`categories` text DEFAULT '[]' NOT NULL,
	`started_at` text,
	`updated_at` text,
	`ended_at` text,
	FOREIGN KEY (`channel_id`) REFERENCES `channels`(`id`) ON UPDATE no action ON DELETE cascade
);
