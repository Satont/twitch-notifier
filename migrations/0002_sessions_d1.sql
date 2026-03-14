-- Migration number: 0002 	 2026-03-14T11:45:27.036Z
CREATE TABLE `sessions` (
  `key` text PRIMARY KEY NOT NULL,
  `value` text NOT NULL,
  `expires_at` integer
);
