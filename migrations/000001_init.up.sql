CREATE TABLE IF NOT EXISTS `xxx` (
    `id` integer PRIMARY KEY AUTOINCREMENT,
    `project_id` text NOT NULL,
    `project_name` text NOT NULL,
    `user_id` text NOT NULL,
    `user_name` text NOT NULL,
    `action` text NOT NULL,
    `created_at` datetime DEFAULT (datetime('now', 'localtime'))
);

CREATE TABLE IF NOT EXISTS `clusters`(
    `id` integer PRIMARY KEY AUTOINCREMENT,
    `name` text,
    `auth_url` text
);

CREATE TABLE IF NOT EXISTS `image_upload_tasks` (
    `id` integer PRIMARY KEY AUTOINCREMENT,
    `project_id` text,
    `image_id` text,
    `image_name` text,
    `size` integer,
    `cached` integer,
    `uploaded` integer
);