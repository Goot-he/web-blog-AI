-- ============================================================
-- 由 ent schema 自动生成，请勿手改
-- 重新生成：
--   go run entgo.io/ent/cmd/ent schema ./internal/infras/repo/ent/schema --dialect mysql --version 8.0
-- ============================================================


-- Create "categories" table
CREATE TABLE `categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL,
  `updated_at` timestamp NULL,
  `deleted_at` timestamp NULL,
  `category` varchar(255) NOT NULL,
  `type_id` bigint unsigned NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
-- Create "interaction_count" table
CREATE TABLE `interaction_count` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL,
  `updated_at` timestamp NULL,
  `object_id` bigint unsigned NOT NULL,
  `object_type` enum('article','comment','life') NOT NULL,
  `interaction_type` enum('like','comment','view','favor') NOT NULL,
  `count` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_object` (`object_id`, `object_type`, `interaction_type`)
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
-- Create "interaction_like" table
CREATE TABLE `interaction_like` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL,
  `updated_at` timestamp NULL,
  `user_id` bigint unsigned NOT NULL,
  `object_type` enum('article','comment','life') NOT NULL,
  `object_id` bigint unsigned NOT NULL,
  `status` enum('unknown','thumb_up','nothing') NOT NULL DEFAULT 'nothing',
  `version` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_like_user_object` (`user_id`, `object_type`, `object_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
-- Create "users" table
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL,
  `updated_at` timestamp NULL,
  `deleted_at` timestamp NULL,
  `birth_date` timestamp NULL,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(255) NULL DEFAULT '',
  `avatar` varchar(255) NULL DEFAULT '',
  `phone` varchar(255) NULL DEFAULT '',
  `sex` enum('male','female','unknown') NOT NULL DEFAULT 'unknown',
  `role` enum('admin','user') NOT NULL DEFAULT 'user',
  `signature` varchar(255) NULL DEFAULT '',
  `status` enum('approved','blocked') NOT NULL DEFAULT 'approved',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `username` (`username`)
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
-- Create "articles" table
CREATE TABLE `articles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL,
  `updated_at` timestamp NULL,
  `deleted_at` timestamp NULL,
  `published_at` timestamp NULL,
  `title` varchar(255) NOT NULL,
  `summary` varchar(255) NOT NULL,
  `content` longtext NOT NULL,
  `coverage` varchar(255) NULL DEFAULT '',
  `is_top` bool NULL DEFAULT false,
  `is_published` bool NULL DEFAULT true,
  `comment_count` bigint unsigned NULL DEFAULT 0,
  `like_count` bigint unsigned NULL DEFAULT 0,
  `favor_count` bigint unsigned NULL DEFAULT 0,
  `view_count` bigint unsigned NULL DEFAULT 0,
  `category_id` bigint unsigned NOT NULL,
  `author_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `articles_categories_articles` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE NO ACTION,
  CONSTRAINT `articles_users_articles` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`) ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
-- Create "comments" table
CREATE TABLE `comments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL,
  `updated_at` timestamp NULL,
  `deleted_at` timestamp NULL,
  `root_id` bigint unsigned NULL DEFAULT 0,
  `reply_id` bigint unsigned NULL DEFAULT 0,
  `reply_user_id` bigint unsigned NULL DEFAULT 0,
  `content` text NOT NULL,
  `like_count` bigint unsigned NOT NULL DEFAULT 0,
  `child_count` bigint unsigned NOT NULL DEFAULT 0,
  `article_id` bigint unsigned NOT NULL,
  `publish_user_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `comments_articles_comments` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON DELETE NO ACTION,
  CONSTRAINT `comments_users_comments` FOREIGN KEY (`publish_user_id`) REFERENCES `users` (`id`) ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
-- Create "user_status" table
CREATE TABLE `user_status` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL,
  `updated_at` timestamp NULL,
  `deleted_at` timestamp NULL,
  `article_count` bigint unsigned NOT NULL DEFAULT 0,
  `like_count` bigint unsigned NOT NULL DEFAULT 0,
  `receive_like_count` bigint unsigned NOT NULL DEFAULT 0,
  `favor_count` bigint unsigned NOT NULL DEFAULT 0,
  `receive_favor_count` bigint unsigned NOT NULL DEFAULT 0,
  `follow_count` bigint unsigned NOT NULL DEFAULT 0,
  `receive_follow_count` bigint unsigned NOT NULL DEFAULT 0,
  `view_count` bigint unsigned NOT NULL DEFAULT 0,
  `comment_count` bigint unsigned NOT NULL DEFAULT 0,
  `user_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uk_user_status_user` (`user_id`),
  CONSTRAINT `user_status_users_statu` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci;

